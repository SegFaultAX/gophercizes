package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/segfaultax/gophercizes/link/walk"
	"golang.org/x/net/html"
)

type (
	Seen map[string]bool

	Page struct {
		URL   string
		Links []walk.Link
	}
)

func main() {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}

	fmt.Println(Walk("http://example.com", c))
}

func Walk(root string, client *http.Client) ([]Page, error) {
	var out []Page
	work := []string{root}
	seen := make(Seen)

	for len(work) > 0 {
		cur := work[0]
		work = work[1:]

		if seen[cur] {
			continue
		}
		fmt.Println("Fetching: ", cur)

		seen[cur] = true
		resp, err := client.Get(cur)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		node, err := html.Parse(resp.Body)
		if err != nil {
			return nil, err
		}

		links := walk.ExtractSync(node)
		out = append(out, Page{
			URL:   cur,
			Links: links,
		})

		u, _ := url.Parse(cur)
		for _, l := range links {
			u2, _ := u.Parse(l.HRef)
			work = append(work, u2.String())
		}
	}

	return out, nil
}
