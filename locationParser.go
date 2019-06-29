package main

import (
	"encoding/json"
	"fmt"
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

type Location struct {
	Id     int               `json:"id"`
	Names  map[string]string `json:"names"`
	Parent *Location         `json:"parent"`
}

const MAX_ERRORS = 10

func GetLocations() []Location {
	reqF := "https://www.avito.ru/web/1/slocations?limit=10000&q=%s"
	generator := qGen{}
	generator.add()
	client := http.DefaultClient

	var ret []Location

	stored := make(map[int]bool)

	errOccurred := 0

	for {
		query := generator.next()
		if query == "" {
			break
		}

		time.Sleep(100 * time.Millisecond)
		fmt.Println("requesting ", string(query))

		resp, err := client.Get(fmt.Sprintf(reqF, string(query))) // INFO: выдает максимум 1000 локаций без возможности указания смещения или страниц
		if err != nil {
			log.Println("couldn't get response:", err)
			errOccurred++
			if errOccurred > MAX_ERRORS {
				os.Exit(1)
			}
			time.Sleep(time.Second)
			generator.back()
			continue
		}

		raw, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()

		var locs map[string]map[string][]Location

		err = json.Unmarshal(raw, &locs)
		if err != nil {
			log.Println("couldn't unmarshal response:", err)
			errOccurred++
			if errOccurred > MAX_ERRORS {
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

func LoadOrParseLocs() []Location {
	var ret []Location
	fmt.Println("Looking for locations_save file")

	data, err := ioutil.ReadFile("locations_save")
	if err != nil {
		fmt.Println("couldn't read/find saved locations, parsing from site:", err)
		ret = GetLocations()
		dataToSave, err := json.Marshal(ret)
		if err != nil {
			fmt.Println("couldn't marshal parsed locations:", err)
			return ret
		}
		_ = ioutil.WriteFile("locations_save", dataToSave, 644)
		return ret
	}

	fmt.Println("Loading locations from locations_save file")
	err = json.Unmarshal(data, &ret)
	if err != nil {
		log.Print("couldn't unmarshal loaded locations:", err)
	}

	return ret
}
