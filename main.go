package main

import (
	"fmt"

	"example.com/m/v2/cli"
	"example.com/m/v2/fetch"
	"example.com/m/v2/headless"
)

func main() {
	fetch.Scrape()
	headless.Example()
	fmt.Println("ran")
	cli.Run()
}
