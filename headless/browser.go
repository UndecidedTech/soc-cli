package headless

import (
	"fmt"

	"github.com/go-rod/rod"
)

func Example() {
	page := rod.New().MustConnect().MustPage("https://tinyurl.is/3Ka9?sport=soccer")

	url, error := page.MustElement(".btn-secondary").Property("href")

	if error != nil {
		fmt.Println(error)
	} else {
		fmt.Println(url)
	}

	// resp, error := http.Get(url.String())
	// defer resp.Body.Close()
	// body, err := io.ReadAll(resp.Body)
	// fmt.Println(string(body), err)

}
