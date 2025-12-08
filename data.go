package main

import (
	"fmt"
	"os"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
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
