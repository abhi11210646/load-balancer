package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting Load Balancer on PORT", Config.Port)

	// if len(Config.Servers) == 0 {
	// 	fmt.Println("No nodes present to balance load.Exiting...")
	// 	os.Exit(0)
	// }
	// fmt.Println("Heart Beat Interval:", Config.HeartBeatInterval, "ms")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	log.Fatal(http.ListenAndServe(Config.Port, nil))
}
