package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/segfaultax/gophercizes/cyoa/adventure"
)

func main() {
	path := flag.String("adventure", "", "path to adventure file")
	flag.Parse()
	if path == nil || *path == "" {
		fmt.Println("path is required")
		os.Exit(1)
	}

	adv, err := adventure.LoadAdventure(*path)
	if err != nil {
		fmt.Println("error while loading adventure:", err)
		os.Exit(1)
	}

	// game := &adv.Game{
	// 	Adventure:  adv,
	// 	CurrentArc: "intro",
	// }

	// game.handleInput()

	// game.PlayCli()

	game := &adventure.WebGame{
		Adventure:  adv,
		DefaultArc: "intro",
	}

	http.Handle("/", game)
	fmt.Println("Server running on 8080...")
	log.Fatalf("server %s:", http.ListenAndServe(":8080", nil))
}
