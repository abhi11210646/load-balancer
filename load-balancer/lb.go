package main

import (
	"log"
	"net/http"

	algo "github.com/load-balancer/load-balancer/internal/algorithm"
	node "github.com/load-balancer/load-balancer/internal/server"
)

type LoadBalancer struct {
	Port    string
	servers []node.Server
	algo    algo.RoutingAlgorithm
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		Port: Config.Port,
		servers: []node.Server{
			{Url: "http://localhost:3001", Active: true},
			{Url: "http://localhost:3002", Active: true},
		},
		algo: &algo.RoundRobin{},
	}
}

func (lb *LoadBalancer) getServer() node.Server {
	return lb.algo.GetNextServer(lb.servers)
}

func (lb *LoadBalancer) ListenAndServe() error {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World\n"))
	})
	log.Println("Load balancer is listening on", lb.Port)
	return http.ListenAndServe(lb.Port, nil)
}
