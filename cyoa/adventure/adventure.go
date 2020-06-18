package adventure

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

type (
	Adventure map[string]Arc

	Arc struct {
		Title   string
		Story   []string
		Options []ArcOption
	}

	ArcOption struct {
		Text string
		Arc  string
	}

	Game struct {
		Adventure  Adventure
		CurrentArc string
	}

	WebGame struct {
		Adventure  Adventure
		DefaultArc string
	}
)

func LoadAdventure(path string) (Adventure, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var adv Adventure = make(map[string]Arc)
	err = json.NewDecoder(f).Decode(&adv)
	if err != nil {
		return nil, err
	}

	return adv, nil
}

func (g *Game) PlayCli() {
	playing := true
	for playing {
		arc, ok := g.Adventure[g.CurrentArc]
		if !ok {
			fmt.Println("unknown adventure arc, how did you get here?")
			os.Exit(1)
		}

		fmt.Println(arc.Title, "\n")
		for _, p := range arc.Story {
			fmt.Println(p, "\n")
		}

		for i, o := range arc.Options {
			fmt.Printf("%d: %s\n", i+1, o.Text)
		}
		fmt.Println("(q to quit)")

		playing = g.handleInput()
	}
}

func (g *Game) handleInput() bool {
	for {
		var input string
		fmt.Printf("What will you do? ")
		fmt.Scanf("%s", &input)
		switch input {
		case "q", "quit", "exit":
			return false
		default:
			i, err := strconv.Atoi(input)
			if err != nil || i <= 0 || i > len(g.Adventure[g.CurrentArc].Options) {
				fmt.Println("Bad choice...")
				break
			}
			g.CurrentArc = g.Adventure[g.CurrentArc].Options[i-1].Arc

			return true
		}
	}
}

func (g *WebGame) HandleDefaultArc() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "story.html", g.Adventure[g.DefaultArc])
	}
}

func (g *WebGame) HandleArc() echo.HandlerFunc {
	return func(c echo.Context) error {
		arc := c.Param("arc")
		if _, ok := g.Adventure[arc]; !ok {
			return echo.ErrNotFound
		}
		return c.Render(http.StatusOK, "story.html", g.Adventure[arc])
	}
}
