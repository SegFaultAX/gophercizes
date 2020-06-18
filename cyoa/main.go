package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/segfaultax/gophercizes/cyoa/adventure"
)

type (
	Template struct {
		templates *template.Template
	}
)

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

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
	fmt.Println(game)

	// http.Handle("/", game)
	// fmt.Println("Server running on 8080...")
	// log.Fatalf("server %s:", http.ListenAndServe(":8080", nil))

	t := &Template{
		templates: template.Must(template.ParseGlob("cyoa/templates/*.html")),
	}

	e := echo.New()
	e.Renderer = t

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", game.HandleDefaultArc())
	e.GET("/:arc", game.HandleArc())

	e.Start(":8080")
}
