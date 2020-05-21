package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var problems string
var limit int

type (
	validator func(string) bool

	question struct {
		query    string
		response validator
	}
)

func init() {
	flag.StringVar(&problems, "file", "quiz/problems.csv", "path to problems csv")
	flag.IntVar(&limit, "limit", 20, "time limit (in seconds)")
}

func anyOf(s []string) validator {
	return func(resp string) bool {
		for _, w := range s {
			if w == resp {
				return true
			}
		}
		return false
	}
}

func anyOfLowercase(s []string) validator {
	return func(resp string) bool {
		for _, w := range s {
			if strings.ToLower(w) == strings.ToLower(resp) {
				return true
			}
		}
		return false
	}
}

func loadQuestions(path string) ([]question, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var qs []question
	rows := csv.NewReader(f)
	for {
		row, err := rows.Read()
		if err == io.EOF {
			break
		}
		if err != nil && !errors.Is(err, csv.ErrFieldCount) {
			return nil, err
		}
		qs = append(qs, question{
			query:    row[0],
			response: anyOf(row[1:]),
		})
	}

	return qs, nil
}

func main() {
	flag.Parse()

	if problems == "" {
		printUsageAndExit(1, "file is required")
	}

	qs, err := loadQuestions(problems)
	if err != nil {
		fmt.Println("error while loading questions:", err)
		os.Exit(1)
	}

	t := time.After(time.Duration(limit) * time.Second)

	correct := 0
	ans := make(chan string)
	for i, q := range qs {
		go func() {
			fmt.Printf("Question %d: %s - ", i+1, q.query)
			var resp string
			fmt.Scanf("%s", &resp)
			ans <- resp
		}()

		select {
		case <-t:
			fmt.Printf("You got %d correct!\n", correct)
			return
		case a := <-ans:
			if q.response(a) {
				correct++
			}
		}
	}

	fmt.Printf("You got %d correct!\n", correct)
}

func printUsageAndExit(code int, message string) {
	fmt.Println(message)
	flag.PrintDefaults()
	os.Exit(code)
}
