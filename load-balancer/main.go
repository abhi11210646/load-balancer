package main

import (
	"log"

	algo "github.com/load-balancer/load-balancer/internal/algorithm"
)

func main() {
	lb := NewLoadBalancer()
	lb.setRoutingAlgorithm(&algo.RoundRobin{})
	err := lb.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
