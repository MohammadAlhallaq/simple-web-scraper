package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(colly.AllowedDomains(
		"my.majarrastg.com",
	))

	c.OnHTML("a[href]", func(element *colly.HTMLElement) {
		err := element.Request.Visit(element.Attr("href"))
		if err != nil {
			panic(err.Error())
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit("https://my.majarrastg.com/")
	if err != nil {
		panic(err.Error())
	}
}
