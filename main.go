package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"golang.org/x/net/html"
)

const (
	baseURL = "https://www.theguardian.com/food/series/how-to-cook-the-perfect----"
	// TODO: replace with something determined from the index page.
	pageCount = 22
)

var (
	recipeURLRegex = regexp.MustCompile(`^https://www.theguardian.com/(food|lifeandstyle/wordofmouth)/(?P<Year>\d{4})/(?P<Month>[a-z]{3})/(?P<Day>\d{2})/(?P<Title>[a-z-]+)$`)
)

func main() {
	var extractLinks func(*html.Node)
	extractLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					if recipeURLRegex.MatchString(a.Val) {
						fmt.Println(a.Val)
					}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractLinks(c)
		}
	}

	for i := 0; i <= pageCount; i++ {
		rsp, err := http.Get(fmt.Sprintf("%s?page=%d", baseURL, i))
		if err != nil {
			log.Fatal(err)
		}

		doc, err := html.Parse(rsp.Body)
		if err != nil {
			log.Fatal(err)
		}
		extractLinks(doc)
	}
}
