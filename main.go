package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/jedib0t/go-pretty/table"
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

			args = []PathStat{}
			for _, entry := range entries {
				args = append(args, PathStat{path: entry.Name()})
			}
		}
	}

	// Sort the list

	rows := arraylist.New()

	for _, pathStat := range args {
		node := getSize(&pathStat)
		rows.Add(node)
	}

	rows.Sort(func(a, b interface{}) int { return b.(*Node).sumBytes - a.(*Node).sumBytes })

	// Output

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Path", "Files", "Bytes"})

	var reportTotalFiles int
	var reportTotalBytes int

	it := rows.Iterator()
	for it.Begin(); it.Next(); {
		value := it.Value().(*Node)
		t.AppendRows([]table.Row{
			{value.path, value.ancestorCount, value.sumBytes},
		})
		reportTotalFiles += value.ancestorCount
		reportTotalBytes += value.sumBytes
	}

	t.AppendFooter(table.Row{"TOTALS", reportTotalFiles, reportTotalBytes})
	t.Render()
}

func getSize(pathStat *PathStat) *Node {
	if pathStat.stat == nil {
		fileInfo, _ := os.Stat(pathStat.path)
		pathStat.stat = &fileInfo
	}

	if (*pathStat.stat).IsDir() {
		fmt.Printf("Get DIR %v\n", pathStat.path)
		list := collectSizes(pathStat.path)
		node, _ := list.Get(0)
		return node.(*Node)
	} else {
		fmt.Printf("Get FILE %v\n", pathStat.path)
		node := New(pathStat.path, true)
		node.bytes = (*pathStat.stat).Size()
		node.sumBytes = int(node.bytes)
		return &node
	}
}

func collectSizes(path string) *arraylist.List {
	fmt.Printf("Walking %v\n", path)
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
