package scrapper

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/gyujae/starcraft_scrapper/utils"
)

var MAP_URL = "https://910map.tistory.com"
var ASL_SEASON_MAP_URL = "https://910map.tistory.com/category/ASL/ASL%20%EC%8B%9C%EC%A6%8C%20"

type Map struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Date  string `json:"date"`
	Url   string `json:"url"`
	Image string `json:"image"`
}

type ASLMap struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	Url    string `json:"url"`
	Image  string `json:"image"`
	Season string `json:"season"`
}

func ASLMapScrapper() []ASLMap {
	var maps []ASLMap
	c := make(chan []ASLMap)
	for i := 1; i <= 13; i++ {
		go getASLMapData(ASL_SEASON_MAP_URL, i, c)
	}
	for i := 1; i <= 13; i++ {
		mapDatas := <-c
		maps = append(maps, mapDatas...)
	}

	return maps
}

func MapScraper() []Map {
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

func getASLMapData(url string, season int, mainC chan<- []ASLMap) {
	var maps []ASLMap
	c := make(chan ASLMap)
	targetURL := ASL_SEASON_MAP_URL + strconv.FormatInt(int64(season), 10)
	fmt.Println("✅ Current:", targetURL)
	res, err := http.Get(targetURL)
	utils.CheckErr(err)
	defer res.Body.Close()
	utils.CheckResponseCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	utils.CheckErr(err)

	mapDataCart := doc.Find(".post-item")

	mapDataCart.Each(func(i int, s *goquery.Selection) {
		go extractedASLMapData(s, c, season)
	})
	for i := 0; i < mapDataCart.Length(); i++ {
		mapData := <-c
		maps = append(maps, mapData)
	}

	mainC <- maps
}

func extractedASLMapData(s *goquery.Selection, c chan<- ASLMap, season int) {
	id := s.Find("a").AttrOr("href", "a")
	title := s.Find(".title").Text()
	date := s.Find(".date").Text()
	url := MAP_URL + id
	image := s.Find("img").AttrOr("src", "img")
	c <- ASLMap{Title: title, Date: date, Url: url, Image: image, ID: id[1:], Season: strconv.Itoa(season)}
}

func getMapData(pageNum int, mainC chan<- []Map) {
	var maps []Map
	c := make(chan Map)
	targetURL := MAP_URL + "/?page=" + strconv.FormatInt(int64(pageNum), 10)
	fmt.Println("✅ Current:", targetURL)
	res, err := http.Get(targetURL)
	utils.CheckErr(err)
	defer res.Body.Close()
	utils.CheckResponseCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	utils.CheckErr(err)

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
	id := s.Find("a").AttrOr("href", "a")
	title := s.Find(".title").Text()
	date := s.Find(".date").Text()
	url := MAP_URL + id
	image := s.Find("img").AttrOr("src", "img")
	c <- Map{Title: title, Date: date, Url: url, Image: image, ID: id[1:]}
}

func getPagination() int {
	var mapNum int
	res, err := http.Get(MAP_URL)
	utils.CheckErr(err)
	defer res.Body.Close()
	utils.CheckResponseCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	utils.CheckErr(err)

	doc.Find(".post-header").Each(func(i int, s *goquery.Selection) {
		mapNum, err = strconv.Atoi(s.Find("em").Text())
		utils.CheckErr(err)
	})

	if (mapNum % 12) != 0 {
		return (mapNum / 12) + 1
	}

	return mapNum / 12

}
