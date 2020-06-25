package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/html"
)

type (
	Link struct {
		Text string
		HRef string
	}
)

func main() {
	path := flag.String("path", "", "path to html file")
	flag.Parse()
	if path == nil || *path == "" {
		fmt.Println("path is required")
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}

	node, err := html.Parse(f)
	if err != nil {
		log.Fatal(err)
	}

	link := make(chan Link)
	go func() {
		extractAnchors(node, link)
		close(link)
	}()

	for l := range link {
		fmt.Println(l)
	}
}

func extractText(node *html.Node) string {
	if node == nil {
		return ""
	}

	if node.Type == html.TextNode {
		return node.Data + extractText(node.NextSibling)
	}

	return extractText(node.FirstChild) + extractText(node.NextSibling)
}

func extractHRef(node *html.Node) string {
	if node == nil {
		return ""
	}

	for _, attr := range node.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}

func extractAnchors(root *html.Node, links chan<- Link) {
	if root == nil {
		return
	}

	if root.Type == html.ElementNode && root.Data == "a" {
		links <- Link{
			Text: extractText(root.FirstChild),
			HRef: extractHRef(root),
		}
	}

	extractAnchors(root.FirstChild, links)
	extractAnchors(root.NextSibling, links)
}
