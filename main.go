package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
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

func catsToJson(cats CategoryTree) {
	data, err := json.MarshalIndent(cats, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	_ = ioutil.WriteFile("categories.json", data, 644)
}

func moscowCats() {
	cats := GetCategoriesTree("moskva", true)

	data, err := json.MarshalIndent(cats, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	_ = ioutil.WriteFile("categoriesMoscow.json", data, 644)
}

func cli() {
	text := "Введите номер команды: " +
		"\n\t1) Получить список локаций" +
		"\n\t2) Получить дерево категорий" +
		"\n\t3) Получить дерево категорий Москвы" +
		"\n\t4) Получить дерево категорий другого города(not implemented)" +
		"\n\t0) Выход\n"

	for {
		fmt.Print(text)
		var in string
		_, err := fmt.Scanln(&in)
		if err != nil {
			log.Println(err)
			continue
		}
		n, err := strconv.Atoi(in)
		if err != nil {
			log.Println(err)
			continue
		}
		switch n {
		case 1:
			locs := LoadOrParseLocs()
			locsToJsonTree(locs)
		case 2:
			cats := GetCategoriesTree("rossiya", false)
			catsToJson(cats)
		case 3:
			moscowCats()
		case 0:
			return
		}
	}
}

func main() {
	cli()
}
