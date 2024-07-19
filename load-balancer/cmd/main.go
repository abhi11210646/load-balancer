package main

import (
	"flag"
	"log"
	"strings"

	lb "github.com/load-balancer/load-balancer/internal"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

var (
	PORT = flag.Int("p", 8080, "Load balancer PORT")
)

func main() {
	var nodes stringSlice
	flag.Var(&nodes, "n", "List of servers to balance load")
	flag.Parse()

	loadBalancer := lb.NewLoadBalancer(*PORT)
	loadBalancer.SetRoutingAlgorithm(&lb.RoundRobin{})
	loadBalancer.SetServers(nodes)

	if err := loadBalancer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
