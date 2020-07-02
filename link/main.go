package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/segfaultax/gophercizes/link/walk"
	"golang.org/x/net/html"
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

	link := make(chan walk.Link)
	go func() {
		walk.ExtractAnchors(node, link)
		close(link)
	}()

	for l := range link {
		fmt.Println(l)
	}
}
