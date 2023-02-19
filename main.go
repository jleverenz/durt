package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/emirpasic/gods/lists/arraylist"
)

type Node struct {
	path          string
	bytes         int64
	ancestorCount int
	children      []*Node
}

func New(path string) Node {
	return Node{
		path:     path,
		bytes:    0,
		children: []*Node{},
	}
}

func main() {
	fmt.Println("sizeup")

	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	var dirNode *Node
	dirNode = nil
	var topNode *Node
	topNode = nil

	dirMap := map[string]*Node{}

	infos := []Node{}

	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		node := New(path)
		node.bytes = info.Size()

		// Initial case
		if dirNode == nil {
			dirNode = &node
			topNode = &node

			if info.IsDir() {
				dirMap[path] = &node
			}

			infos = append(infos, node)
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

		infos = append(infos, node)
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
		return
	}

	for i, info := range infos {
		fmt.Printf("PRINT %v: %+v\n", i, info)
	}

	// walk the tree
	var printTree func(*Node, int)
	printTree = func(node *Node, depth int) {
		fmt.Printf("%v%v\n", strings.Repeat(" ", depth), node.path)
		for _, n := range node.children {
			printTree(n, depth+2)
		}
	}
	//	printTree(topNode, 0)

	var countAncestors func(*Node) int
	countAncestors = func(node *Node) int {
		count := 1
		for _, child := range node.children {
			count = count + countAncestors(child)
		}
		node.ancestorCount = count
		return count
	}
	countAncestors(topNode)

	list := arraylist.New()

	var buildList func(*Node)
	buildList = func(node *Node) {
		list.Add(node)
		for _, child := range node.children {
			buildList(child)
		}
	}
	buildList(topNode)

	list.Sort(func(a, b interface{}) int { return b.(*Node).ancestorCount - a.(*Node).ancestorCount })

	it := list.Iterator()
	for it.Begin(); it.Next(); {
		value := it.Value().(*Node)
		fmt.Printf("%v %v\n", value.ancestorCount, value.path)
	}

	var printFlatTree func(*Node)
	printFlatTree = func(node *Node) {
		fmt.Printf("[%v] %v\n", node.ancestorCount, node.path)
		for _, n := range node.children {
			printFlatTree(n)
		}
	}
	// printFlatTree(topNode)
}
