package webscraper

import (
	"fmt"

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

func ParseArticleLinks(doc *goquery.Document) []string {
	var result = []string{}
	matchArticleLink := fmt.Sprintf(`[href]`)

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
	if doc == nil {
		return nil, fmt.Errorf("Invalid goquery document!")
	}

	articleInfo := ArticleInfo{}

	articleInfo.titleH1 = doc.Find("h1").First().Text()
	articleInfo.titleH4 = doc.Find("h4").First().Text()
	articleInfo.datePublished = doc.Find("time").First().Text()
	articleInfo.paragraph = doc.Find("p").First().Text()

	return &articleInfo, nil
}
