package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Node struct {
	path     string
	bytes    int64
	children []*Node
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

	// 	// Build a config map:
	// confMap := map[string]string{}
	// for _, v := range myconfig {
	//     confMap[v.Key] = v.Value
	// }

	// // And then to find values by key:
	// if v, ok := confMap["key1"]; ok {
	//     // Found
	// }

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

	var printFlatTree func(*Node)
	printFlatTree = func(node *Node) {
		fmt.Println(node.path)
		for _, n := range node.children {
			printFlatTree(n)
		}
	}
	printFlatTree(topNode)
}
