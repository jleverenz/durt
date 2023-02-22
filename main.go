package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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

func main() {
	args := []PathStat{{path: "."}}

	if len(os.Args) > 1 {
		args = []PathStat{}
		for _, path := range os.Args[1:] {
			args = append(args, PathStat{path: path})
		}
	}

	if len(args) == 1 {
		stat, _ := os.Stat(args[0].path)
		args[0].stat = &stat

		if stat.IsDir() {

			entries, _ := os.ReadDir(args[0].path)

			expanding := args[0]
			args = []PathStat{}
			for _, entry := range entries {
				// fmt.Printf("Adding %v\n", path.Join(expanding.path, entry.Name()))
				args = append(args, PathStat{path: path.Join(expanding.path, entry.Name())})
			}
		}
	}

	// Create and sort the list

	rows := []*Node{}

	for _, pathStat := range args {
		node := getSize(&pathStat)
		rows = append(rows, node)
	}

	sort.Sort(ByBytes(rows))

	displaySortedResults(rows)
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
		return node.(*Node)
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
	countAncestors(topNode)

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

	for _, value := range nodes {
		pct := float64(value.sumBytes) / float64(reportTotalBytes) * 100
		t.AppendRows([]table.Row{
			{value.path, value.ancestorCount, ByteSize(value.sumBytes), fmt.Sprintf("%.01f %%", pct)},
		})
	}

	t.AppendFooter(table.Row{"TOTALS", reportTotalFiles, ByteSize(reportTotalBytes), "100.0 %"})

	nameTransformer := text.Transformer(func(val interface{}) string {
		return text.WrapHard(val.(string), longestNonWrap)
	})

	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Number:      1,
			Transformer: nameTransformer,
		},
		{
			Number:      4,
			Align:       text.AlignRight,
			AlignHeader: text.AlignRight,
		},
	})

	t.SetOutputMirror(os.Stdout)

	t.Render()
}
