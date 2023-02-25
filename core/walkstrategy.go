package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var WalkStrategy = struct {
	Execute func([]PathStat) []PathResult
}{Execute: executeWalkStrategy}

type walkNode struct {
	path          string
	bytes         int64
	isFile        bool
	ancestorCount int
	sumBytes      int
	children      []*walkNode
}

func executeWalkStrategy(pathStats []PathStat) []PathResult {
	results := []PathResult{}

	for _, pathStat := range pathStats {
		if !checkPathExclusion(pathStat.Path) {
			node := getSize(&pathStat)
			if node != nil {

				results = append(results, PathResult{
					Path:       node.path,
					TotalBytes: node.sumBytes,
					TotalFiles: node.ancestorCount,
				})
			}
		}
	}

	return results
}

func checkPathExclusion(path string) bool {
	for _, re := range GlobalOpts.Exclusions {
		if re.MatchString(path) {
			return true
		}
	}

	return false
}

func getSize(pathStat *PathStat) *walkNode {
	if pathStat.Stat == nil {
		fileInfo, err := os.Stat(pathStat.Path)
		if err != nil {
			fmt.Println(err)
		}
		pathStat.Stat = &fileInfo
	}

	if (*pathStat.Stat).IsDir() {
		node := collectSizes(pathStat.Path)
		if node != nil {
			return node
		} else {
			return nil
		}
	} else {
		node := walkNode{path: pathStat.Path, isFile: true}
		node.bytes = (*pathStat.Stat).Size()
		node.sumBytes = int(node.bytes)
		return &node
	}
}

func collectSizes(path string) *walkNode {
	var dirNode *walkNode
	dirNode = nil
	var topNode *walkNode
	topNode = nil

	dirMap := map[string]*walkNode{}

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

		node := walkNode{path: path, isFile: !info.IsDir()}
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

		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
		return nil
	}

	countAncestors(topNode)

	return topNode
}

func countAncestors(node *walkNode) (int, int64) {
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
