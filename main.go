package main

import (
	"breakcrawl/webscraper"
	"flag"
	"fmt"
	"log"
	"time"
)

// Replace with CLI arg
const baseURL = "https://www.breakit.se"
const urlMatchRoute = "/artikel/"
const defaultMaxCrawlDepth = 2
const defaultMaxConcurrentRequests = 10

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v .\n", name, time.Since(start))
	}
}

func main() {
	defer timer("Webcrawler")()
	var maxCrawlDepth int
	var maxConcurrentRequests int

	flag.IntVar(&maxCrawlDepth, "d", defaultMaxCrawlDepth, "webcrawler recursion depth")
	flag.IntVar(&maxConcurrentRequests, "p", defaultMaxConcurrentRequests, "upper limit to webcrawler concurrent GET requests")
	flag.Parse()

	fmt.Printf("Max scraping recursion depth: %v\n", maxCrawlDepth)
	fmt.Printf("Max concurrent GET requests: %v\n", maxConcurrentRequests)
	fmt.Printf("\n\nScraping...\n\n")

	scraper := &webscraper.Scraper{
		BaseURL:               baseURL,
		MatchingRoute:         urlMatchRoute,
		MaxCrawlDepth:         maxCrawlDepth,
		MaxConcurrentRequests: maxConcurrentRequests,
	}

	err := scraper.ScrapePage(baseURL, 0)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(scraper.Articles); i++ {
		fmt.Println(scraper.Articles[i])
	}
	fmt.Printf("Collected %v articles.\n", len(scraper.Articles))
	fmt.Printf("Made %v requests.\n", scraper.NumRequests)

	return
}
