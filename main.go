package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

const BASE_DIR = "therion"
const DAILY_DIR = "daily"

type dailyTemplateFields struct {
	Date string
}

//go:embed templates
var templatesFS embed.FS

func createDir(dir string) string {
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
	createDir(baseDir)
	dailyDir := filepath.Join(baseDir, DAILY_DIR)

	templates, err := template.ParseFS(templatesFS, "templates/*.md")
	if err != nil {
		panic(err)
	}

	for _, arg := range args {
		if arg == "today" {

			createDir(dailyDir)
			if err != nil {
				panic(err)
			}
			today := time.Now()
			todayExists := false
			err = filepath.WalkDir(dailyDir, func(name string, _ os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if name == dailyDir {
					return nil
				}
				date, err := time.Parse(time.DateOnly, filepath.Base(name))
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
				dailyTemplate := dailyTemplateFields{Date: time.Now().Format(time.RFC3339)}
				err = templates.ExecuteTemplate(file, "daily.md", dailyTemplate)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
