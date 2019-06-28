package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type qGen struct {
	chars []int32
}

func (q *qGen) next() string {
	ret := string(q.chars)
	last := len(q.chars) - 1

	if last < 0 {
		return ""
	}

	if q.chars[last] == 'я' {
		q.chars = q.chars[:last]
		last--
		if last < 0 {
			return ""
		}
	}
	q.chars[last]++

	return ret
}

func (q *qGen) add() {
	q.chars = append(q.chars, 'а')
}

func (q *qGen) back() {
	last := len(q.chars) - 1
	if last < 0 {
		return
	}
	if q.chars[last] == 'а' {
		return
	}
	q.chars[last]--
}

type location struct {
	Id     int               `json:"id"`
	Names  map[string]string `json:"names"`
	Parent *location         `json:"parent"`
}

func getLocations() []location {
	reqF := "https://www.avito.ru/web/1/slocations?limit=10000&q=%s"
	generator := qGen{}
	generator.add()
	client := http.DefaultClient

	ret := []location{}

	stored := make(map[int]bool)

	maxErrors := 10
	errOccured := 0

	for {
		c := generator.next()
		if c == "" {
			break
		}

		time.Sleep(100 * time.Millisecond)
		fmt.Println("requesting ", string(c))

		resp, err := client.Get(fmt.Sprintf(reqF, string(c))) // INFO: выдает максимум 1000 локаций без возможности указания смещения или страниц
		if err != nil {
			log.Println(err, "_1")
			errOccured++
			if errOccured > maxErrors {
				os.Exit(1)
			}
			time.Sleep(time.Second)
			generator.back()
			continue
		}
		raw, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()

		var locs map[string]map[string][]location

		err = json.Unmarshal(raw, &locs)
		if err != nil {
			log.Println(err, "_2")
			errOccured++
			if errOccured > maxErrors {
				os.Exit(1)
			}
			time.Sleep(time.Second)
			generator.back()
			continue
		}

		num := len(locs["result"]["locations"])
		fmt.Println("\tgot", num, "locations")

		if num >= 1000 {
			generator.back()
			generator.add()
			continue
		}

		newLocs := 0
		for _, v := range locs["result"]["locations"] {
			if _, ok := stored[v.Id]; !ok {
				ret = append(ret, v)
				stored[v.Id] = true
				newLocs++
			}
		}
		fmt.Println("\tadded", newLocs, "new locations")
	}
	return ret
}

func loadOrParseLocs() []location {
	var ret []location
	data, err := ioutil.ReadFile("locations_save")
	if err != nil {
		fmt.Println("Couldn't read/find saved locations, parsing from site:", err)
		ret = getLocations()
		dataToSave, err := json.Marshal(ret)
		if err != nil {
			fmt.Println("Couldn't marshal parsed locations:", err)
			return ret
		}
		_ = ioutil.WriteFile("locations_save", dataToSave, 644)
		return ret
	}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		log.Print("Couldn't unmarshal loaded locations:", err)
	}
	return ret
}

type locTree struct {
	Id    int       `json:"id"`
	Name  string    `json:"name"`
	Child []locTree `json:"locations"`
}

func locsToJsonTree(locs []location) {
	var tree []locTree
	var withParents []location

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

	_ = ioutil.WriteFile("locs.json", data, 644)
}

func main() {
	locs := loadOrParseLocs()

	locsToJsonTree(locs)

	//buffer := bytes.NewBufferString("")
	//for _, l := range locs {
	//	buffer.WriteString(l.Names["1"] + "\n")
	//}
	//_ = ioutil.WriteFile("locations.txt", buffer.Bytes(), 644)

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

	doc.Find(".js-search-form-region").Find("*").Each(func(i int, s *goquery.Selection) {
		band := s.Find("a").Text()
		title := s.Find("i").Text()
		text := s.Text()
		fmt.Printf("Review %d: %s - %s %s\n", i, band, title, text)
	})

}
