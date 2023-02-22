package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/urfave/cli/v2"
)

type Node struct {
	path          string
	bytes         int64
	isFile        bool
	ancestorCount int
	sumBytes      int
	children      []*Node
}

func New(path string, isFile bool) Node {
	return Node{path: path, isFile: isFile}
}

type ByBytes []*Node

func (a ByBytes) Len() int           { return len(a) }
func (a ByBytes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByBytes) Less(i, j int) bool { return a[j].sumBytes < a[i].sumBytes } // Reverse

type PathStat struct {
	path string
	stat *fs.FileInfo
}

type ProgramOptions struct {
	head       bool
	exclusions []*regexp.Regexp
}

var globalOpts ProgramOptions

func main() {
	globalOpts = ProgramOptions{}

	app := &cli.App{
		Name:  "dusc",
		Usage: "disk utilization simple comparison",
		Flags: []cli.Flag{
			// TODO it'd be nice to allow --head, --head 30, etc; seems this flag
			// parsing module doesn't support that
			&cli.BoolFlag{
				Name:  "head",
				Usage: "display the top 20 records",
			},
			&cli.StringSliceFlag{
				Name:  "exclude",
				Usage: "exclude paths by regex",
			},
		},
		Action: func(cCtx *cli.Context) error {
			globalOpts.head = cCtx.Bool("head")

			exclusions := cCtx.StringSlice("exclude")
			for _, exc := range exclusions {
				globalOpts.exclusions = append(globalOpts.exclusions, regexp.MustCompile(exc))
			}

			mainAction(cCtx.Args().Slice())
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func mainAction(cliArgs []string) {
	args := resolveArgList(cliArgs)

	// Create and sort the list

	rows := []*Node{}

	for _, pathStat := range args {
		if !checkPathExclusion(pathStat.path) {
			node := getSize(&pathStat)
			if node != nil {
				rows = append(rows, node)
			}
		}
	}

	sort.Sort(ByBytes(rows))

	displaySortedResults(rows)
}

func checkPathExclusion(path string) bool {
	for _, re := range globalOpts.exclusions {
		if re.MatchString(path) {
			return true
		}
	}

	return false
}

func resolveArgList(args []string) []PathStat {
	resolved := []PathStat{{path: "."}}

	if len(args) > 0 {
		resolved = []PathStat{}
		for _, path := range args {
			resolved = append(resolved, PathStat{path: path})
		}
	}

	if len(resolved) == 1 {
		stat, _ := os.Stat(resolved[0].path)
		resolved[0].stat = &stat

		if stat.IsDir() {
			os.Chdir(resolved[0].path)
			entries, _ := os.ReadDir(".")

			// expanding := args[0]
			resolved = []PathStat{}
			for _, entry := range entries {
				resolved = append(resolved, PathStat{path: entry.Name()})
			}
		}
	}

	return resolved
}

func getSize(pathStat *PathStat) *Node {
	if pathStat.stat == nil {
		fileInfo, err := os.Stat(pathStat.path)
		if err != nil {
			fmt.Println(err)
		}
		pathStat.stat = &fileInfo
	}

	if (*pathStat.stat).IsDir() {
		list := collectSizes(pathStat.path)
		node, _ := list.Get(0)
		if node != nil {
			return node.(*Node)
		} else {
			return nil
		}
	} else {
		node := New(pathStat.path, true)
		node.bytes = (*pathStat.stat).Size()
		node.sumBytes = int(node.bytes)
		return &node
	}
}

func collectSizes(path string) *arraylist.List {
	var dirNode *Node
	dirNode = nil
	var topNode *Node
	topNode = nil

	dirMap := map[string]*Node{}

	list := arraylist.New()

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		for _, re := range globalOpts.exclusions {
			if re.MatchString(path) {
				if info.IsDir() {
					return fs.SkipDir
				} else {
					return nil
				}
			}
		}

		node := New(path, !info.IsDir())
		if !info.IsDir() {
			node.bytes = info.Size()
		}

		// Initial case
		if dirNode == nil {
			dirNode = &node
			topNode = &node

			if info.IsDir() {
				dirMap[path] = &node
			}

			// infos = append(infos, node)
			list.Add(&node)
			return nil
		}

		dirPath := filepath.Dir(path)

		// TODO there is a bug here if passing two directories as cli args, and one
		// starts with "./**"
		if dirPath != dirNode.path {
			dirNode = dirMap[dirPath]
		}

		dirNode.children = append(dirNode.children, &node)

		if info.IsDir() {
			dirMap[path] = &node
		}

		list.Add(&node)
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
		return nil
	}

	var countAncestors func(*Node) (int, int64)
	countAncestors = func(node *Node) (int, int64) {
		count := 0
		if node.isFile {
			count = 1
		}

		byteTotal := node.bytes

		for _, child := range node.children {
			childrenCount, childrenBytes := countAncestors(child)
			count = count + childrenCount
			byteTotal = byteTotal + int64(childrenBytes)
		}
		node.ancestorCount = count
		node.sumBytes = int(byteTotal)
		return count, byteTotal
	}

	if topNode != nil {
		countAncestors(topNode)
	}

	return list
}

func displaySortedResults(nodes []*Node) {
	t := table.NewWriter()

	t.AppendHeader(table.Row{"Path", "Files", "Bytes", "Pct"})

	// calc totals
	var reportTotalFiles int
	var reportTotalBytes int
	longestNonWrap := 0
	for _, value := range nodes {
		reportTotalFiles += value.ancestorCount
		reportTotalBytes += value.sumBytes
		pathLength := len(value.path)
		if pathLength < 80 && pathLength > longestNonWrap {
			longestNonWrap = pathLength
		}
	}

	if globalOpts.head {
		nodes = nodes[0:20]
	}

	for _, value := range nodes {
		pct := float64(value.sumBytes) / float64(reportTotalBytes) * 100
		t.AppendRows([]table.Row{
			{value.path, value.ancestorCount, ByteSize(value.sumBytes), pct},
		})
	}

	t.AppendFooter(table.Row{"TOTALS", reportTotalFiles, ByteSize(reportTotalBytes), "100.0 %"})

	nameTransformer := text.Transformer(func(val interface{}) string {
		return text.WrapHard(val.(string), longestNonWrap)
	})

	percentTransformer := text.Transformer(func(val interface{}) string {
		if val.(float64) < 1 {
			return fmt.Sprintf("< 1.0 %%")
		} else {
			return fmt.Sprintf("%.01f %%", val)
		}
	})

	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Number:      1,
			Transformer: nameTransformer,
		},
		{
			Number: 2,
			VAlign: text.VAlignBottom,
		},
		{
			Number: 3,
			VAlign: text.VAlignBottom,
		},
		{
			Number:      4,
			Align:       text.AlignRight,
			VAlign:      text.VAlignBottom,
			Transformer: percentTransformer,
		},
	})

	t.SetOutputMirror(os.Stdout)

	t.Render()
}
