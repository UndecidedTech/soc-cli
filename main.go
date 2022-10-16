package main

import (
	"fmt"

	"example.com/m/v2/cli"
	"example.com/m/v2/fetch"
)

func main() {
	// headless.Example()
	fetch.GetJson()
	fmt.Println("ran")
	cli.Start()
}
