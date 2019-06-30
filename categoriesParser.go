package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var gClient http.Client

func getCustom(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("couldn't create new request:", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:67.0) Gecko/20100101 Firefox/67.0")

	return gClient.Do(req)
}

type CategoryTree struct {
	Name          string         `json:"name"`
	Url           string         `json:"url"`
	Count         int            `json:"count"`
	SubCategories []CategoryTree `json:"subCategories"`
}

func getCountFromUrl(url string) int {
	fmt.Println("Getting the number of announcements from:", url)

	du := time.Duration(rand.Intn(11))
	time.Sleep((du*100 + 1000) * time.Millisecond) //rand sleep from 1 to 2 second with step 0.1 second
	resp, err := getCustom(url)
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

	countStr := doc.Find(".page-title-count").Text()
	countStr = strings.Join(strings.Fields(countStr), "")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		log.Println("Couldn't convert '"+countStr+"':", err)
	}

	return count
}

func buildSubTree(tree *CategoryTree, doCount bool) func(i int, s *goquery.Selection) {
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

			if doCount {
				subTree.Count = getCountFromUrl(subTree.Url)
			}

			s.ChildrenFiltered("ul").Each(buildSubTree(&subTree, doCount))
			tree.SubCategories = append(tree.SubCategories, subTree)
		}

		if s.Is("ul") {
			s.ChildrenFiltered("li").Each(buildSubTree(tree, doCount))
		}
	}
}

func getSubCategories(url string, doCount bool) []CategoryTree {
	time.Sleep(100 * time.Millisecond)
	resp, err := getCustom(url)
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

	doc.Find("*[class^='rubricator-list']").Each(buildSubTree(&temp, doCount))

	return temp.SubCategories
}

func GetCategoriesTree(region string, doCount bool) CategoryTree {
	fmt.Println("Getting the main categories")

	url := "https://www.avito.ru/" + region

	time.Sleep(100 * time.Millisecond)
	resp, err := getCustom(url)
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

	if doCount {
		countStr := doc.Find(".breadcrumbs-link-count").Text()
		countStr = strings.Join(strings.Fields(countStr), "")
		count, err := strconv.Atoi(countStr)
		if err != nil {
			log.Println("Couldn't convert '"+countStr+"':", err)
		}
		ret.Count = count
	}

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
		ret.SubCategories = append(ret.SubCategories, getSubCategories(homeUrl+url, doCount)...)
	})

	return ret
}
