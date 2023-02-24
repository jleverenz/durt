package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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
	Head       bool
	Exclusions []*regexp.Regexp
	MyVal      string
}

var GlobalOpts ProgramOptions

func MainAction(cliArgs []string) {
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
	for _, re := range GlobalOpts.Exclusions {
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
		node := list[0]
		if node != nil {
			return node
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

func collectSizes(path string) []*Node {
	var dirNode *Node
	dirNode = nil
	var topNode *Node
	topNode = nil

	dirMap := map[string]*Node{}

	list := []*Node{}

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		for _, re := range GlobalOpts.Exclusions {
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
			list = append(list, &node)
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

		list = append(list, &node)
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
