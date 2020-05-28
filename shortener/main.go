package main

import (
	"fmt"
	"log"
	"net/http"
)

var links map[string]string = map[string]string{
	"/hello": "http://www.example.com",
}

func MapHandler(links map[string]string, fallback http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		target, ok := links[r.URL.Path]
		if !ok {
			fallback(w, r)
			return
		}
		http.Redirect(w, r, target, http.StatusTemporaryRedirect)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("fallback: not found!"))
}

func main() {
	fmt.Println("Hello, world!")

	http.HandleFunc("/", MapHandler(links, notFoundHandler))
	log.Fatalf("server %s:", http.ListenAndServe(":8080", nil))
}
