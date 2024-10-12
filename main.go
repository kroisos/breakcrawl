package main

import (
	"breakcrawl/webscraper"
	"fmt"
	"log"
)

// Replace with CLI arg
const baseURL = "https://breakit.se"
const maxCrawlDepth = 1

func main() {
	scraper := &webscraper.Scraper{
		BaseURL:       baseURL,
		MaxCrawlDepth: maxCrawlDepth,
	}

	err := scraper.ScrapeArticles(baseURL, 0)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(scraper.Articles); i++ {
		fmt.Println(scraper.Articles[i])
	}

	return
}
