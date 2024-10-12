package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

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

type Scraper struct {
	articles []*ArticleInfo
	wg       sync.WaitGroup
	lock     sync.Mutex
}

func (sr *Scraper) append(a *ArticleInfo) {
	sr.lock.Lock()
	sr.articles = append(sr.articles, a)
	sr.lock.Unlock()
}

func (sr *Scraper) ScrapeArticles(url string, curCrawlDepth int) error {

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
		sr.append(article)
	}

	// If crawl recursion depth not reached continue scraping.
	if curCrawlDepth < maxCrawlDepth {
		links := ParseArticleLinks(doc)
		for i := 0; i < len(links); i++ {
			sr.wg.Add(1)
			go func() {
				defer sr.wg.Done()
				err := sr.ScrapeArticles(links[i], curCrawlDepth+1)
				if err != nil {
					log.Println(err)
				}
			}()
		}
	}

	// Condition to avoid every recursive crawl to wait for every other. Only first crawl waits for children.
	if curCrawlDepth == 0 {
		sr.wg.Wait()
	}

	return nil
}

func FetchDoc(url string) (*goquery.Document, error) {
	// Request the HTML page.
	page, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer page.Body.Close()
	if page.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", page.StatusCode, page.Status)
	}

	// Load the HTML document
	return goquery.NewDocumentFromReader(page.Body)
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
