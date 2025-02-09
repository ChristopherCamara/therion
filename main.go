package main

import (
	"embed"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

const BASE_DIR = ".therion"
const DAILY_DIR = "daily"

type operation string

const ()

type file struct {
	path string
	data []byte
}

type TemplateFields struct {
	Name string
	Date string
}

//go:embed templates
var templatesFS embed.FS

func main() {
	args := os.Args[1:]

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	baseDir := filepath.Join(homeDir, BASE_DIR)
	err = os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	templates, err := template.ParseFS(templatesFS, "templates/*.md")
	if err != nil {
		panic(err)
	}

	for _, arg := range args {
		if arg == "--today" {
			today := time.Now()
			name := ""
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
		} // else if arg == "--sync" {
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
