package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

type CategoryTree struct {
	Name          string
	Url           string
	SubCategories []CategoryTree
}

func buildSubTree(tree *CategoryTree) func(i int, s *goquery.Selection) {
	return func(i int, s *goquery.Selection) {
		homeUrl := "https://www.avito.ru"

		if s.Is("li") {
			var subTree CategoryTree

			aElem := s.Find("a")
			name, exist := aElem.Attr("title")
			if !exist {
				log.Println("couldn't find item name")
				return
			}
			subTree.Name = name

			url, exist := aElem.Attr("href")
			if !exist {
				log.Println("couldn't find item url")
				return
			}
			subTree.Url = homeUrl + url

			s.ChildrenFiltered("ul").Each(buildSubTree(&subTree))
			tree.SubCategories = append(tree.SubCategories, subTree)
		}

		if s.Is("ul") {
			s.ChildrenFiltered("li").Each(buildSubTree(tree))
		}
	}
}

func getSubCategories(url string) []CategoryTree {
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

	var temp CategoryTree

	doc.Find("*[class^='rubricator-list']").Each(buildSubTree(&temp))

	return temp.SubCategories
}

func GetCategoriesTree() CategoryTree {
	fmt.Println("Getting the main categories")

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

		fmt.Println("Getting subcategories for", name)
		ret.SubCategories = append(ret.SubCategories, getSubCategories(homeUrl+url)...)
	})

	return ret
}
