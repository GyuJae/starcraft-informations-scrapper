package main

import (
	"fmt"

	mapscrapper "github.com/gyujae/starcraft_scrapper/mapScrapper"
)

func main() {
	maps := mapscrapper.Scraper()
	fmt.Println(maps)
}
