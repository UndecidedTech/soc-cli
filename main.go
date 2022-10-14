package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
	// "example.com/m/v2/fetch"
)

func main() {
	// fetch.Scrape()
	// fmt.Println("Hello, World!!")
	// https://reditt1.soccerstreams.net/
	c := colly.NewCollector(
		colly.AllowedDomains("books.toscrape.com"),
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

	c.Visit("https://books.toscrape.com")

	fmt.Println("ran")
}
