package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
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

type CategoryTree struct {
	Name          string
	Url           string
	SubCategories []CategoryTree
}

func GetSubCategories(url string) []CategoryTree {
	return nil
}

func GetCategoriesTree2() CategoryTree {
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

	homeUrl := "https://www.avito.ru"
	ret := CategoryTree{Name: "All", Url: url}

	doc.Find(".category-map-title").Each(func(i int, s *goquery.Selection) {
		selection := s.Find("a")

		name := selection.Text()
		name = strings.Join(strings.Fields(name), "") // removing all whitespaces
		
		url, exist := selection.Attr("href")
		if !exist {
			log.Println("category url not found, skip...")
			return
		}

		tCat := CategoryTree{Name: name, Url: homeUrl + url, SubCategories: GetSubCategories(homeUrl + url)}
		ret.SubCategories = append(ret.SubCategories, tCat)
	})

	return ret
}
