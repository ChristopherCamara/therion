package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Syncer interface {
	getFiles() []file
	uploadFile(*file)
	downloadFile(*file)
}

func cleanFilePath(s string) string {
	split := strings.SplitAfter(s, fmt.Sprintf("%s%s", BASE_DIR, string(filepath.Separator)))
	return split[len(split)-1]
}

func findFile(files []file, f *file) (*file, bool) {
	for _, file := range files {
		if file.path == f.path {
			return &file, true
		}
	}
	return &file{}, false
}

func sync(path string, syncer Syncer) {
	var localFiles []file
	err := filepath.WalkDir(path, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			localFiles = append(localFiles, file{path: cleanFilePath(path), modTime: info.ModTime(), data: []byte{}})
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	backupFiles := syncer.getFiles()
	for _, file := range localFiles {
		backup, found := findFile(backupFiles, &file)
		if (found && backup.modTime.Unix() < file.modTime.Unix()) || !found {
			bytes, err := os.ReadFile(filepath.Join(path, file.path))
			if err != nil {
				panic(err)
			}
			file.data = bytes
			syncer.uploadFile(&file)
		} else if found && backup.modTime.Unix() > file.modTime.Unix() {
			syncer.downloadFile(backup)
		}
	}
}
