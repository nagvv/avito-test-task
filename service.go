package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var mainPage = `<!doctype html>
		<html>
		<body>
		<button class="tablinks" onclick="window.location.href='/locations';">получить локации</button>
		<button class="tablinks" onclick="window.location.href='/categories';">получить категории</button>
		<button class="tablinks" onclick="window.location.href='/cityMoscow';">получить категории Москвы</button>
		</form>
		</body>
		</html>`

var backPageF = `<!doctype html>
		<html>
		<body>
		<p>%s</p>
		<button class="tablinks" onclick="window.location.href='/';">назад</button>
		</form>
		</body>
		</html>`

func mainHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(mainPage))
	if err != nil {
		log.Fatal(err)
	}
}

func locationsHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("locationsTree.json")
	if err != nil {
		_, err = w.Write([]byte(fmt.Sprintf(backPageF, "wait until locations tree is parsed")))
		return
	}
	defer file.Close()

	fileHeader := make([]byte, 512)
	file.Read(fileHeader)
	fileContentType := http.DetectContentType(fileHeader)

	fileStat, _ := file.Stat()
	fileSize := strconv.FormatInt(fileStat.Size(), 10)

	w.Header().Set("Content-Disposition", "attachment; filename=locationsTree.json")
	w.Header().Set("Content-Type", fileContentType)
	w.Header().Set("Content-Length", fileSize)

	file.Seek(0, 0)
	_, err = io.Copy(w, file)
	if err != nil {
		log.Fatal("Couldn't send locations tree", err)
	}
}

func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("categories.json")
	if err != nil {
		_, err = w.Write([]byte(fmt.Sprintf(backPageF, "wait until categories tree is parsed")))
		return
	}
	defer file.Close()

	fileHeader := make([]byte, 512)
	file.Read(fileHeader)
	fileContentType := http.DetectContentType(fileHeader)

	fileStat, _ := file.Stat()
	fileSize := strconv.FormatInt(fileStat.Size(), 10)

	w.Header().Set("Content-Disposition", "attachment; filename=categories.json")
	w.Header().Set("Content-Type", fileContentType)
	w.Header().Set("Content-Length", fileSize)

	file.Seek(0, 0)
	_, err = io.Copy(w, file)
	if err != nil {
		log.Fatal("Couldn't send locations tree", err)
	}
}

func categoriesMoscowHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("categoriesMoscow.json")
	if err != nil {
		_, err = w.Write([]byte(fmt.Sprintf(backPageF, "wait until Moscow categories is parsed")))
		return
	}
	defer file.Close()

	fileHeader := make([]byte, 512)
	file.Read(fileHeader)
	fileContentType := http.DetectContentType(fileHeader)

	fileStat, _ := file.Stat()
	fileSize := strconv.FormatInt(fileStat.Size(), 10)

	w.Header().Set("Content-Disposition", "attachment; filename=categoriesMoscow.json")
	w.Header().Set("Content-Type", fileContentType)
	w.Header().Set("Content-Length", fileSize)

	file.Seek(0, 0)
	_, err = io.Copy(w, file)
	if err != nil {
		log.Fatal("Couldn't send locations tree", err)
	}
}

func startService() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/locations", locationsHandler)
	http.HandleFunc("/categories", categoriesHandler)
	http.HandleFunc("/cityMoscow", categoriesMoscowHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("couldn't start service", err)
	}
}
