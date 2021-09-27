package main

import (
	"fmt"
	"log"
	"net/http"
)

type Personne struct {
	id   int
	age  int
	name string
}

func hello(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, Personne{1, 26, "mokadem"})
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	fmt.Println("starting server on port :4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}
