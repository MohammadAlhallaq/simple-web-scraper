package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"net/url"
	"os"
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

	var uniqueHrefs []string

	c := colly.NewCollector(colly.AllowedDomains(visitingUrl.domain))
	c.AllowURLRevisit = true

	c.OnHTML("body", func(e *colly.HTMLElement) {
		reader := bytes.NewReader(e.Response.Body)

		doc, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			fmt.Println("Error parsing the body:", err)
			return
		}

		doc.Find("a").Each(func(i int, selection *goquery.Selection) {
			link, _ := selection.Attr("href")
			parsedURL, err := url.Parse(link)
			if err != nil {
				fmt.Println("Error parsing URL:", err)
				return
			}
			if parsedURL.Hostname() == "" && parsedURL.Path != "" {
				link = visitingUrl.fullPath + parsedURL.Path
			}
			uniqueHrefs = append(uniqueHrefs, link)

			err = os.MkdirAll("Links", os.ModePerm)
			file, err := os.Create("./links/home.csv")
			defer file.Close()

			if err != nil {
				return
			}

			writer := csv.NewWriter(file)
			defer writer.Flush()

			for _, record := range uniqueHrefs {
				if err := writer.Write([]string{record}); err != nil {
					log.Fatalln("error writing record to file", err)
				}
			}
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit(visitingUrl.fullPath)
	if err != nil {
		panic(err.Error())
	}
}
