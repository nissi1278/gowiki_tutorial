package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}

type TemplateData struct {
	Page  *Page
	Pages []*Page
}

func (p *Page) Save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func getDataTitleInDir() ([]string, error) {
	entries, err := os.ReadDir(dataPath)
	if err != nil {
		return nil, fmt.Errorf("ディレクトリ %s が見つかりませんでした。", dataPath)
	}

	var titles []string
	for _, entry := range entries {
		// ディレクトリはdata/に格納していないが、チェック。
		if entry.IsDir() {
			continue
		}
		title := strings.TrimSuffix(entry.Name(), ".txt")

		titles = append(titles, title)
	}

	if len(titles) == 0 {
		return nil, fmt.Errorf("ディレクトリ %s にデータファイルが見つかりませんでした。", dataPath)
	}

	return titles, nil
}

func loadAllPages() ([]*Page, error) {
	titles, err := getDataTitleInDir()
	if err != nil {
		return nil, err
	}
	var pages []*Page
	for _, title := range titles {
		page, err := loadPage(title)
		if err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}
	return pages, nil
}

func loadPage(title string) (*Page, error) {
	filename := dataPath + title + ".txt"
	body, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, err
}

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
func renderPage(w http.ResponseWriter, p *Page, tmpl string) {
	data := &TemplateData{Page: p}
	renderTemplate(w, tmpl, data)
}

func renderPages(w http.ResponseWriter, p []*Page, tmpl string) {
	data := &TemplateData{Pages: p}
	renderTemplate(w, tmpl, data)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data *TemplateData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var templates = template.Must(template.ParseFiles(tmplPath+"list.html", tmplPath+"edit.html", tmplPath+"view.html"))
var validPath = regexp.MustCompile("^/(edit|view|save)/([a-zA-Z0-9]+)$")
var tmplPath = "tmpl/"
var dataPath = "data/"

func main() {
	http.HandleFunc("/", listHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
