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

		fmt.Printf("Received request from %s \n", r.RemoteAddr)
		fmt.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Proto)
		fmt.Printf("Host: %s\n", r.Host)
		fmt.Printf("User-Agent: %s\n", r.UserAgent())
		fmt.Printf("Accept: %s\n", r.Header.Get("Accept"))

		// fmt.Fprintf(w, "Response from %s \n", PORT)
	})
	log.Fatal(http.ListenAndServe(PORT, nil))
}
