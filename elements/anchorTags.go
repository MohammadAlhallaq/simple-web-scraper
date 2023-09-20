package elements

import (
	"encoding/csv"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"os"
)

type Elements interface {
	Collect()
}

type AnchorTag struct {
	Links    *goquery.Document
	FullPath string
}

func (a AnchorTag) Collect() error {

	links := a.Links.Find("a")

	if links.Length() > 0 {

		uniqueHrefs := map[string]bool{}

		err := os.MkdirAll("Links", os.ModePerm)
		file, err := os.Create("./links/home.csv")
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				panic(err.Error())
			}
		}(file)

		if err != nil {
			panic(err.Error())
		}

		writer := csv.NewWriter(file)
		defer writer.Flush()

		links.Each(func(i int, selection *goquery.Selection) {
			link, _ := selection.Attr("href")
			parsedURL, _ := url.Parse(link)

			if parsedURL.Hostname() == "" && parsedURL.Path != " " {
				link = a.FullPath + parsedURL.Path
			}
			if !uniqueHrefs[link] {
				uniqueHrefs[link] = true
			}
		})
		for record, _ := range uniqueHrefs {
			if err := writer.Write([]string{record}); err != nil {
				return err
			}
		}
	}
	return nil
}
