package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

type SizeInfo struct {
	path  string
	bytes int64
}

func main() {
	fmt.Println("sizeup")

	infos := []SizeInfo{}

	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		infos = append(infos, SizeInfo{path, info.Size()})
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
		return
	}

	for i, info := range infos {
		fmt.Printf("PRINT %v: %+v\n", i, info)
	}
}
