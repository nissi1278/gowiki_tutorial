package main

import (
	"fmt"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) Save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, err
}

func main() {
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.Save()

	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
}
