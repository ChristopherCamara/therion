package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

const SERVER_PORT = 3000

//go:embed assets/*.*
var assets embed.FS

//go:embed assets/templates/*
var templatesFS embed.FS

type server struct {
	templates *template.Template
	devMode   bool
}

func newServer(devMode bool) *server {
	server := new(server)
	server.devMode = devMode
	return server
}

func (s *server) loadTemplates() {
	if s.devMode {
		path, err := filepath.Abs("assets/templates/*.html")
		if err != nil {
			panic(err)
		}
		s.templates, err = template.ParseGlob(path)
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		s.templates, err = template.ParseFS(templatesFS, "assets/templates/*.html")
		if err != nil {
			panic(err)
		}
	}

}

func (s *server) start() {
	s.loadTemplates()
	if s.devMode {
		http.Handle("/assets/", http.FileServer(http.Dir(".")))
	} else {
		http.Handle("/assets/", http.FileServerFS(assets))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buffer := &bytes.Buffer{}
		err := s.templates.ExecuteTemplate(buffer, "index.html", s)
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, buffer)
	})
	//http.HandleFunc("/sidebar", )

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", SERVER_PORT), nil))
}
