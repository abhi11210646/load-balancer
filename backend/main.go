package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var PORT int

func main() {
	flag.IntVar(&PORT, "p", 3001, "Server PORT")
	flag.Parse()
	fmt.Println("Starting server on PORT", PORT)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("---------------------------------------\n")

		fmt.Printf("Received request from %s \n", r.RemoteAddr)
		fmt.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Proto)
		fmt.Printf("Host: %s\n", r.Host)
		fmt.Printf("User-Agent: %s\n", r.UserAgent())
		fmt.Printf("Accept: %s\n", r.Header.Get("Accept"))

		fmt.Printf("---------------------------------------\n")

		w.Header().Add("Content-Type", "application/json")

		responseData := map[string]string{
			"status": "OK",
			"server": r.Host,
		}
		data, _ := json.Marshal(responseData)
		w.Write(data)
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK\n")
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%+d", PORT), nil))
}
