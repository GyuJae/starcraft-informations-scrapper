package mapscrapper

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var MAP_URL = "https://910map.tistory.com"

type Map struct {
	Title string `json:"title"`
	Date  string `json:"date"`
	Url   string `json:"url"`
	Image string `json:"image"`
}

func Scraper() []Map {
	var maps []Map
	totalPage := getPagination()
	c := make(chan []Map)
	for i := 1; i <= totalPage; i++ {
		go getMapData(i, c)
	}
	for i := 1; i <= totalPage; i++ {
		mapDatas := <-c
		maps = append(maps, mapDatas...)
	}

	return maps

}

func getMapData(pageNum int, mainC chan<- []Map) {
	var maps []Map
	c := make(chan Map)
	targetURL := MAP_URL + "/?page=" + strconv.FormatInt(int64(pageNum), 10)
	fmt.Println("✅ Current:", targetURL)
	res, err := http.Get(targetURL)
	checkErr(err)
	defer res.Body.Close()
	checkCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	mapDataCart := doc.Find(".post-item")

	mapDataCart.Each(func(i int, s *goquery.Selection) {
		go extractedMapData(s, c)
	})
	for i := 0; i < mapDataCart.Length(); i++ {
		mapData := <-c
		maps = append(maps, mapData)
	}

	mainC <- maps
}

func extractedMapData(s *goquery.Selection, c chan<- Map) {
	title := s.Find(".title").Text()
	date := s.Find(".date").Text()
	url := MAP_URL + s.Find("a").AttrOr("href", "a")
	image := s.Find("img").AttrOr("src", "img")
	c <- Map{Title: title, Date: date, Url: url, Image: image}
}

func getPagination() int {
	var mapNum int
	res, err := http.Get(MAP_URL)
	checkErr(err)
	defer res.Body.Close()
	checkCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".post-header").Each(func(i int, s *goquery.Selection) {
		mapNum, err = strconv.Atoi(s.Find("em").Text())
		checkErr(err)
	})

	if (mapNum % 12) != 0 {
		return (mapNum / 12) + 1
	}

	return mapNum / 12

}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatal("❌ Request failed with Status:", res.Status)
	}
}
