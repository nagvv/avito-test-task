package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
		fmt.Println("Parent not found", v.Names["1"], v.Parent.Names["1"])
	}

	data, err := json.MarshalIndent(tree, "", "\t")
	if err != nil {
		fmt.Println("Couldn't marshal locations tree:", err)
		return
	}

	_ = ioutil.WriteFile("locationsTree.json", data, 644)
}

func catsToJson(cats CategoryTree) {
	data, err := json.MarshalIndent(cats, "", "\t")
	if err != nil {
		fmt.Println("Couldn't marshal categories tree:", err)
		return
	}

	_ = ioutil.WriteFile("categories.json", data, 644)
}

func moscowCats() {
	cats, err := GetCategoriesTree("moskva", true)
	if err != nil {
		fmt.Println("Couldn't get categories tree:", err)
		return
	}

	data, err := json.MarshalIndent(cats, "", "\t")
	if err != nil {
		fmt.Println("Couldn't marshal categories tree:", err)
		return
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
			fmt.Println(err)
			continue
		}
		n, err := strconv.Atoi(in)
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch n {
		case 1:
			locs := LoadOrParseLocs()
			locsToJsonTree(locs)
		case 2:
			cats, err := GetCategoriesTree("rossiya", false)
			if err != nil {
				fmt.Println("Не удалось получить дерево категорий")
				continue
			}
			catsToJson(cats)
		case 3:
			moscowCats()
		case 0:
			return
		default:
			fmt.Println("Недопустимый номер")
		}
	}
}

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "cli" {
			cli()
			return
		}
	}

	go func() {
		locs := LoadOrParseLocs()
		locsToJsonTree(locs)
		fmt.Println("Parsing locations finished")

		cats, err := GetCategoriesTree("rossiya", false)
		if err != nil {
			log.Fatal("Не удалось получить дерево категорий:", err)
		}
		catsToJson(cats)
		fmt.Println("Parsing categories finished")

		moscowCats()
		fmt.Println("Parsing Moscow categories finished")
	}()
	startService()
}
