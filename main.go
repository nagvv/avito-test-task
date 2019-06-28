package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type locTree struct {
	Id    int       `json:"id"`
	Name  string    `json:"name"`
	Child []locTree `json:"locations"`
}

func locsToJsonTree(locs []Location) {
	var tree []locTree
	var withParents []Location

	for _, v := range locs {
		var t locTree
		t.Id = v.Id
		t.Name = v.Names["1"]
		if v.Parent != nil {
			withParents = append(withParents, v)
			continue
		}
		tree = append(tree, t)
	}

wp:
	for _, v := range withParents {
		parentId := v.Parent.Id
		for i := range tree {
			if parentId == tree[i].Id {
				tree[i].Child = append(tree[i].Child, locTree{Id: v.Id, Name: v.Names["1"]})
				continue wp
			}
		}
		fmt.Println("parent not found", v.Names["1"], v.Parent.Names["1"])
		///log.Fatal("parent not found", v.Id, v.Parent.Id)
	}

	data, err := json.MarshalIndent(tree, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	_ = ioutil.WriteFile("locationsTree.json", data, 644)
}



func main() {
	locs := LoadOrParseLocs()

	locsToJsonTree(locs)


	time.Sleep(1 * time.Second)

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

	doc.Find(".js-search-form-category").Find("*").Each(func(i int, s *goquery.Selection) {
		band := s.Find("a").Text()
		title := s.Find("i").Text()
		text := s.Text()
		fmt.Printf("Review %d: %s - %s %s\n", i, band, title, text)
	})

}
