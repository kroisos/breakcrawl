package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Replace with CLI arg
const baseURL = "https://breakit.se"

// const articleLinkPattern = "/artiklar/"
const maxCrawlDepth = 1

/*
Find the items:
- Url
- Publishing date
- H1 title
- H4 title
- First paragraph
*/
type ArticleInfo struct {
	url           string
	datePublished string
	titleH1       string
	titleH4       string
	paragraph     string
}

type ScrapeResult struct {
	articles []*ArticleInfo
}

func main() {
	result := &ScrapeResult{}

	err := ScrapeArticles(baseURL, 0, result)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(result.articles); i++ {
		fmt.Println(result.articles[i])
	}

	return
}

func ScrapeArticles(url string, curCrawlDepth int, res *ScrapeResult) error {
	// Validate and adjust url.
	if len(url) < 7 || (len(url) > 3 && url[:4] != "http") {
		url = baseURL + url
	}

	doc, err := FetchDoc(url)
	if err != nil {
		// Don't cancel current scraping because lower level scrapes fail. Simply log then continue.
		if curCrawlDepth == 0 {
			return err
		} else {
			log.Println(err)
		}
	}

	article, err := ParseArticle(doc)
	if err != nil {
		log.Println(err)
	} else {
		res.articles = append(res.articles, article)
	}

	// If crawl recursion depth not reached continue scraping.
	if curCrawlDepth < maxCrawlDepth {
		links := ParseArticleLinks(doc)
		for i := 0; i < len(links); i++ {
			err := ScrapeArticles(links[i], curCrawlDepth+1, res)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}

func FetchDoc(url string) (*goquery.Document, error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	return goquery.NewDocumentFromReader(res.Body)
}

func ParseArticleLinks(doc *goquery.Document) []string {
	var result = []string{}
	// matchArticleLink := fmt.Sprintf(`a[href^="%s"]`, articleLinkPattern)
	matchArticleLink := fmt.Sprintf(`a[href]`)

	// Find article links
	doc.Find(matchArticleLink).Each(func(_ int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			result = append(result, link)
		}
	})
	return result
}

func ParseArticle(doc *goquery.Document) (*ArticleInfo, error) {
	articleInfo := ArticleInfo{}

	articleInfo.titleH1 = doc.Find("h1").First().Text()
	articleInfo.titleH4 = doc.Find("h4").First().Text()
	articleInfo.datePublished = doc.Find("time").First().Text()
	articleInfo.paragraph = doc.Find("p").First().Text()

	return &articleInfo, nil
}
