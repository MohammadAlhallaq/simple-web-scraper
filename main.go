package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"simpleWebScraper/elements"
)

type visitUrl struct {
	domain   string
	fullPath string
}

func main() {
	visitingUrl := &visitUrl{
		domain:   "my.majarra.com",
		fullPath: "https://my.majarra.com",
	}

	c := colly.NewCollector(colly.AllowedDomains(visitingUrl.domain))
	c.AllowURLRevisit = true

	c.OnHTML("body", func(e *colly.HTMLElement) {

		reader := bytes.NewReader(e.Response.Body)

		doc, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			fmt.Println("Error parsing the body:", err)
			return
		}

		anchor := elements.AnchorTag{
			Links:    doc,
			FullPath: visitingUrl.fullPath,
		}
		err = anchor.Collect()
		if err != nil {
			fmt.Println("Error parsing the body:", err)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit(visitingUrl.fullPath)
	if err != nil {
		panic(err.Error())
	}
}
