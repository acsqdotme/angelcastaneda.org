package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	htmlDir     = "./html"
	staticDir   = "./static"
	tmplFileExt = ".tmpl.html"
)

func doesFileExist(pathToFile string) bool {
	info, err := os.Stat(filepath.Clean(pathToFile))
	if err != nil || info.IsDir() {
		return false
	}
	return true
}

func bindTMPL(files ...string) (*template.Template, error) {
	for _, f := range files {
		if !doesFileExist(f) {
			return nil, errors.New("Template file missing " + f)
		}
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	path := strings.Split(r.URL.Path, "/")
	page := path[1]
	if r.URL.Path == "/" {
		page = "index"
	}

	if !doesFileExist(filepath.Join(htmlDir, "pages", page+tmplFileExt)) {
		http.Error(w, "page not found", 404)
		return
	}

	if len(path) > 2 && path[2] == "" {
		http.Redirect(w, r, scheme+"://"+r.Host+"/"+page, 302)
		return
	}

	if len(path) > 2 {
		http.Error(w, "page not found", 404)
		return
	}

	tmpl, err := bindTMPL(
		filepath.Join(htmlDir, "base"+tmplFileExt),
		filepath.Join(htmlDir, "pages", page+tmplFileExt),
	)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server errors", 500)
		return
	}

	data := make(map[string]interface{})
	data["Path"] = r.URL.Path
	tmpl.ExecuteTemplate(w, "base", data)
}

func cbtsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(staticDir, "files", "the_claim_before_the_storm.pdf"))
}

func pgpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.ServeFile(w, r, filepath.Join(staticDir, "files", "angelcastaneda.asc"))
}

func cvHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(staticDir, "files", "cv.pdf"))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(staticDir, "favicon.ico"))
}
