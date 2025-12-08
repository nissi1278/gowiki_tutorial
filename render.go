package main

import "net/http"

type TemplateData struct {
	Page  *Page
	Pages []*Page
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
