package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
)

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		fn(w, r, m[2])
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadAllPages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderPages(w, p, "list")
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderPage(w, p, "view")
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	renderPage(w, p, "edit")
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func getTmplFile(getPath string) string {
	return filepath.Join(TmplPath, getPath)
}

func getDataFile(getPath string) string {
	return filepath.Join(DataPath, getPath)
}

var templates = template.Must(template.ParseFiles(getTmplFile("list.html"), getTmplFile("edit.html"), getTmplFile("view.html")))
var validPath = regexp.MustCompile("^/(edit|view|save)/([a-zA-Z0-9]+)$")

const TmplPath = "tmpl/"
const DataPath = "data/"

func main() {
	http.HandleFunc("/", listHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
