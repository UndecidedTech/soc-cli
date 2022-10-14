package headless

import (
	"fmt"

	"github.com/go-rod/rod"
)

// This example opens https://github.com/, searches for "git",
// and then gets the header element which gives the description for Git.
func Example() {
	// Launch a new browser with default options, and connect to it.
	browser := rod.New().MustConnect()

	// Even you forget to close, rod will close it after main process ends.
	defer browser.MustClose()

	// Create a new page
	page := browser.MustPage("https://tinyurl.is/xhTP?sport=soccer")

	// We use css selector to get the search input element and input "git"
	fmt.Printf("Here")
	button := page.MustElement(".btn")

	fmt.Printf(button)
	// page.MustElement("input").MustInput("git").MustType(input.Enter)

	// Wait until css selector get the element then get the text content of it.
	// text := page.MustElement(".codesearch-results p").MustText()

	// fmt.Println(text)

	// Get all input elements. Rod supports query elements by css selector, xpath, and regex.
	// For more detailed usage, check the query_test.go file.
	// fmt.Println("Found", len(page.MustElements("input")), "input elements")

	// Eval js on the page
	// page.MustEval(`() => console.log("hello world")`)

	// Pass parameters as json objects to the js function. This MustEval will result 3
	// fmt.Println("1 + 2 =", page.MustEval(`(a, b) => a + b`, 1, 2).Int())

	// When eval on an element, "this" in the js is the current DOM element.
	// fmt.Println(page.MustElement("title").MustEval(`() => this.innerText`).String())

	// Output:
	// Git is the most widely used version control system.
	// Found 5 input elements
	// 1 + 2 = 3
	// Search · git · GitHub
}
