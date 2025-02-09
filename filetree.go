package main

import (
	"os"
	"path/filepath"
)

func loadFileTreeFromDir(path string) *Tree[*file] {
	tree := newTree[*file]()
	err := filepath.WalkDir(path, func(path string, entry os.DirEntry, err error) error {
		if entry.IsDir() {

		}
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return tree
}
