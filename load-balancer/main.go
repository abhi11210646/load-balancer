package main

import (
	"log"
)

func main() {
	lb := NewLoadBalancer()
	err := lb.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
