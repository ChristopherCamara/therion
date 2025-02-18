package main

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/google/uuid"
)

const BASE_DIR = ".therion"
const REPLICA_FILE = "replica.txt"
const DAILY_DIR = "daily"

type TemplateFields struct {
	Name string
	Date string
}

//go:embed templates
var templatesFS embed.FS

func main() {
	args := os.Args[1:]
	var ft FileTree
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	baseDir := filepath.Join(homeDir, BASE_DIR)
	ft.AddDirectory(baseDir)

	fmt.Println(ft)

	replicaFile, err := os.OpenFile(filepath.Join(baseDir, REPLICA_FILE), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer replicaFile.Close()
	replicaInfo, err := replicaFile.Stat()
	if err != nil {
		panic(err)
	}
	var replicaID uuid.UUID
	if replicaInfo.Size() == 0 {
		replicaID, err = uuid.NewV7()
		if err != nil {
			panic(err)
		}
		size, err := replicaFile.WriteString(replicaID.String())
		if err != nil || size != len(replicaID.String()) {
			panic("error writing replicaID to file")
		}
	} else {
		id, err := io.ReadAll(replicaFile)
		if err != nil {
			panic(err)
		}
		replicaID, err = uuid.ParseBytes(id)
		if err != nil {
			panic(err)
		}
	}

	templates, err := template.ParseFS(templatesFS, "templates/*.md")
	if err != nil {
		panic(err)
	}

	for _, arg := range args {
		if arg == "--today" {
			dailyDir := filepath.Join(baseDir, DAILY_DIR)
			ft.AddDirectory(dailyDir)
			name := filepath.Join(dailyDir, fmt.Sprintf("%s.md", time.Now().Format(time.DateOnly)))
			file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, os.ModePerm)
			if (err) != nil {
				panic(err)
			}
			defer file.Close()
			info, err := file.Stat()
			if err != nil {
				panic(err)
			}
			if info.Size() == 0 {
				dailyTemplate := TemplateFields{Date: time.Now().Format(time.RFC3339)}
				err = templates.ExecuteTemplate(file, "daily.md", dailyTemplate)
				if err != nil {
					panic(err)
				}
			}
		}
		//    else if arg == "--sync" {
		//		s3 := NewS3Sync()
		//		sync(baseDir, s3)
		//	} else {
		//		err := os.MkdirAll(notesDir, os.ModePerm)
		//		if err != nil {
		//			panic(err)
		//		}
		//		id, err := uuid.NewV7()
		//		if err != nil {
		//			panic(err)
		//		}
		//		file, err := os.Create(filepath.Join(notesDir, fmt.Sprintf("%s.md", id.String())))
		//		if err != nil {
		//			panic(err)
		//		}
		//		defer file.Close()
		//		noteTemplate := TemplateFields{Name: arg}
		//		err = templates.ExecuteTemplate(file, "note.md", noteTemplate)
		//		if err != nil {
		//			panic(err)
		//		}
		//	}
	}
}
