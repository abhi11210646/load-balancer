package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var PORT string

func main() {
	flag.StringVar(&PORT, "p", ":3000", "Server PORT")
	flag.Parse()
	fmt.Println("Starting server on PORT", PORT)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Response from %s", PORT)
	})
	log.Fatal(http.ListenAndServe(PORT, nil))
}
