package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"os"
)

type visitUrl struct {
	domain, fullPath string
}

type job struct {
	title    string
	url      string
	company  string
	location string
	since    string
}

func main() {
	visitingUrl := &visitUrl{
		domain:   "larajobs.com",
		fullPath: "https://larajobs.com/",
	}
	var jobs []job

	c := colly.NewCollector(colly.AllowedDomains(visitingUrl.domain))
	c.AllowURLRevisit = true

	c.OnHTML("a.job-link", func(e *colly.HTMLElement) {
		job := job{}

		job.title = e.ChildText("p.text-lg")
		job.url = e.Attr("data-url")
		job.company = e.ChildText("p.text-sm:first-of-type")
		job.location = e.ChildText("div.flex.items-center.mr-4.mb-1")
		job.since = e.ChildText("div.flex.items-center.mr-4.mb-1 + div.flex.items-center")
		jobs = append(jobs, job)
	})

	c.OnScraped(func(response *colly.Response) {
		file, err := os.Create("job.csv")
		if err != nil {
			fmt.Println(err.Error())
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)
		writer := csv.NewWriter(file)

		headers := []string{
			"title",
			"url",
			"company",
			"location",
			"since",
		}
		writer.Write(headers)

		for _, job := range jobs {
			record := []string{
				job.title,
				job.url,
				job.company,
				job.location,
				job.since,
			}
			writer.Write(record)
		}
		defer writer.Flush()

		fmt.Println("jobs has been scraped successfully")
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit(visitingUrl.fullPath)
	if err != nil {
		panic(err.Error())
	}
}
