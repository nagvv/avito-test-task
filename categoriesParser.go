package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
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
		fmt.Println("Couldn't create new request:", err)
		return nil, err
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
	time.Sleep((du*100 + 1000) * time.Millisecond) // rand sleep from 1 to 2 second with step 0.1 second
	resp, err := getCustom(url)
	if err != nil {
		fmt.Println("Couldn't GET from \""+url+"\":", err)
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("status code error: %d %s\n", resp.StatusCode, resp.Status)
		return 0
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Couldn't create goquery document:", err)
		return 0
	}

	countStr := doc.Find(".page-title-count").Text()
	countStr = strings.Join(strings.Fields(countStr), "")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		fmt.Println("Couldn't convert '"+countStr+"':", err)
		return 0
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
				fmt.Println("Couldn't find item name")
				return
			}
			subTree.Name = name

			url, exist := aElem.Attr("href")
			if !exist {
				fmt.Println("Couldn't find item url")
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

func getSubCategories(url string, doCount bool) ([]CategoryTree, error) {
	time.Sleep(100 * time.Millisecond)
	resp, err := getCustom(url)
	if err != nil {
		fmt.Println("Couldn't GET from \""+url+"\":", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Status code error: %d %s\n", resp.StatusCode, resp.Status)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Couldn't create goquery document:", err)
		return nil, err
	}

	var temp CategoryTree

	doc.Find("*[class^='rubricator-list']").Each(buildSubTree(&temp, doCount))

	return temp.SubCategories, nil
}

func GetCategoriesTree(region string, doCount bool) (CategoryTree, error) {
	fmt.Println("Getting the main categories")

	url := "https://www.avito.ru/" + region

	time.Sleep(100 * time.Millisecond)
	resp, err := getCustom(url)
	if err != nil {
		fmt.Println("Couldn't GET from \""+url+"\":", err)
		return CategoryTree{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Status code error: %d %s\n", resp.StatusCode, resp.Status)
		return CategoryTree{}, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Couldn't create goquery document:", err)
		return CategoryTree{}, err
	}

	homeUrl := "https://www.avito.ru"
	ret := CategoryTree{Name: "All", Url: url}

	if doCount {
		countStr := doc.Find(".breadcrumbs-link-count").Text()
		countStr = strings.Join(strings.Fields(countStr), "")
		count, err := strconv.Atoi(countStr)
		if err != nil {
			fmt.Println("Couldn't convert '"+countStr+"':", err)
		}
		ret.Count = count
	}

	doc.Find(".category-map-title").Each(func(i int, s *goquery.Selection) {
		selection := s.Find("a")

		name := selection.Text()
		name = strings.Join(strings.Fields(name), "") // removing all whitespaces

		url, exist := selection.Attr("href")
		if !exist {
			fmt.Println("Category url not found, skip...")
			return
		}

		fmt.Println("Getting subcategories for", name)
		subCats, err := getSubCategories(homeUrl+url, doCount)
		if err != nil {
			fmt.Println("Couldn't get subcategories for", name)
			return
		}
		ret.SubCategories = append(ret.SubCategories, subCats...)
	})

	return ret, nil
}
