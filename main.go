package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
)

const BASE_DIR = "therion"
const DAILY_DIR = "daily"
const NOTES_DIR = "notes"

type TemplateFields struct {
	Name string
	Date string
}

//go:embed templates
var templatesFS embed.FS

func ensureDirectoryExists(dir string) string {
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	return dir
}

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
	ensureDirectoryExists(baseDir)
	dailyDir := filepath.Join(baseDir, DAILY_DIR)
	notesDir := filepath.Join(baseDir, NOTES_DIR)

	templates, err := template.ParseFS(templatesFS, "templates/*.md")
	if err != nil {
		panic(err)
	}

	for _, arg := range args {
		if arg == "--today" {
			ensureDirectoryExists(dailyDir)
			today := time.Now()
			todayExists := false
			err = filepath.WalkDir(dailyDir, func(name string, _ os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if name == dailyDir {
					return nil
				}
				name = filepath.Base(name)
				name = strings.Replace(name, filepath.Ext(name), "", 1)
				date, err := time.Parse(time.DateOnly, name)
				if err != nil {
					return nil
				}
				if DateEqual(today, date) {
					todayExists = true
					return filepath.SkipAll
				}
				return nil
			})
			if err != nil {
				panic(err)
			}
			if !todayExists {
				file, err := os.Create(filepath.Join(dailyDir, fmt.Sprintf("%s.md", today.Format(time.DateOnly))))
				if (err) != nil {
					panic(err)
				}
				defer file.Close()
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
			ensureDirectoryExists(notesDir)
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
