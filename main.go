package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var MAP_URL = "https://910map.tistory.com/?page="

func main() {
	nums := getPagination()
	getMapInformations(nums)
}

func getMapInformations(pagination int) {

	for i := 1; i <= pagination; i++ {
		res, err := http.Get(MAP_URL + string(i))
		checkErr(err)
		defer res.Body.Close()
		checkCode(res)

		doc, err := goquery.NewDocumentFromReader(res.Body)
		checkErr(err)

		doc.Find(".post-header").Each(func(i int, s *goquery.Selection) {
			postItem := s.Find(".post-item")
			title := postItem.Find(".title").Text()
			fmt.Println(title)
		})
	}

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
		log.Fatal("Request failed with Status:", res.Status)
	}
}
