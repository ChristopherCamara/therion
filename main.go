package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/google/uuid"
)

const BASE_DIR = ".therion"

type operation string

const ()

type file struct {
	Id      uuid.UUID `json:"id"`
	Path    string    `json:"path"`
	modTime time.Time
	data    []byte
	op      operation
}

type TemplateFields struct {
	Name string
	Date string
}

//go:embed templates
var templatesFS embed.FS

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.UTC().Date()
	y2, m2, d2 := date2.UTC().Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

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

	referenceFile, err := os.OpenFile(filepath.Join(baseDir, "files.json"), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer referenceFile.Close()

	templates, err := template.ParseFS(templatesFS, "templates/*.md")
	if err != nil {
		panic(err)
	}

	for _, arg := range args {
		if arg == "--today" {
			if err != nil {
				panic(err)
			}
			today := time.Now()
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
		} else if arg == "--sync" {
			s3 := NewS3Sync()
			sync(baseDir, s3)
		} else {
			err := os.MkdirAll(notesDir, os.ModePerm)
			if err != nil {
				panic(err)
			}
			id, err := uuid.NewV7()
			if err != nil {
				panic(err)
			}
			file, err := os.Create(filepath.Join(notesDir, fmt.Sprintf("%s.md", id.String())))
			if err != nil {
				panic(err)
			}
			defer file.Close()
			noteTemplate := TemplateFields{Name: arg}
			err = templates.ExecuteTemplate(file, "note.md", noteTemplate)
			if err != nil {
				panic(err)
			}
		}
	}
}
