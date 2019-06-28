package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

type CategoryGroup struct {
	Name       string   `json:"name"`
	Categories []string `json:"categories"`
}

func GetCategoriesTree() []CategoryGroup {
	url := "https://www.avito.ru/rossiya"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	ret := []CategoryGroup{{"All", []string{}}}

	doc.Find(".js-search-form-category").Find("*").Each(func(i int, s *goquery.Selection) {
		isGroup := s.HasClass("opt-group")
		name := s.Text()

		if isGroup {
			ret = append(ret, CategoryGroup{Name: name, Categories: []string{}})
			return
		}

		last := len(ret) - 1
		ret[last].Categories = append(ret[last].Categories, name)
	})

	return ret
}
