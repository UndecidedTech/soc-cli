package fetch

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

func Scrape() {
	c := colly.NewCollector(
		colly.AllowedDomains("https://reditt1.soccerstreams.net/"),
	)
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})

	c.Visit("https://reditt1.soccerstreams.net/")

	fmt.Println("ran")
}
