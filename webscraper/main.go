package webscraper

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Scraper struct {
	Articles              []*ArticleInfo
	MaxCrawlDepth         int
	BaseURL               string
	MatchingRoute         string
	MaxConcurrentRequests int
	curConcurrentRequests int
	wg                    sync.WaitGroup
	lock                  sync.Mutex
	NumRequests           int
}

func (sr *Scraper) addNew(a *ArticleInfo) {
	sr.lock.Lock()
	defer sr.lock.Unlock()

	for i := 0; i < len(sr.Articles); i++ {
		if sr.Articles[i].url == a.url {
			return
		}
	}
	sr.Articles = append(sr.Articles, a)
}

func (sr *Scraper) ScrapePage(url string, curCrawlDepth int) error {
	doc, err := fetchDoc(url)
	if err != nil {
		return err
	}

	links := ParseArticleLinks(doc)
	for i := 0; i < len(links); i++ {
		sr.scrape(links[i], 1)
	}

	// Only first crawl waits for children.
	sr.wg.Wait()

	return nil
}

func (sr *Scraper) scrape(url string, curCrawlDepth int) {
	// Validate and adjust url.
	if len(url) < 7 || (len(url) > 3 && url[:4] != "http") {
		url = sr.BaseURL + url
	}

	sr.NumRequests++
	doc, err := fetchDoc(url)
	if err != nil {
		return
	}

	var article *ArticleInfo

	// Validate url is article page. Else continue scraping.
	if sr.validateUrl(url) {
		article, err = ParseArticle(doc)
		if err != nil {
			log.Println(err)
		} else {
			article.url = url
			sr.addNew(article)
		}
	}

	// If crawl recursion depth not reached continue scraping.
	if curCrawlDepth <= sr.MaxCrawlDepth {
		links := ParseArticleLinks(doc)
		for i := 0; i < len(links); i++ {
			if sr.MaxConcurrentRequests > 1 && sr.curConcurrentRequests < sr.MaxConcurrentRequests {
				sr.curConcurrentRequests++
				sr.wg.Add(1)
				go func() {
					defer sr.wg.Done()
					defer func() { sr.curConcurrentRequests-- }()
					sr.scrape(links[i], curCrawlDepth+1)
				}()
				return
			}
			sr.scrape(links[i], curCrawlDepth+1)
		}
	}
	return
}

func fetchDoc(url string) (*goquery.Document, error) {
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

func (sr *Scraper) validateUrl(url string) bool {
	matchString := fmt.Sprintf("%s%s", sr.BaseURL, sr.MatchingRoute)

	match, err := regexp.MatchString(matchString, url)
	if err != nil {
		return false
	}
	return match
}
