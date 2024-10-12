package main

import (
	"fmt"
	"log"
)

// Replace with CLI arg
const baseURL = "https://breakit.se"

// const articleLinkPattern = "/artiklar/"
const maxCrawlDepth = 1

func main() {
	scraper := &Scraper{}

	err := scraper.ScrapeArticles(baseURL, 0)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(scraper.articles); i++ {
		fmt.Println(scraper.articles[i])
	}

	return
}
