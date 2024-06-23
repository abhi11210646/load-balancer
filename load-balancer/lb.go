package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	algo "github.com/load-balancer/load-balancer/internal/algorithm"
	node "github.com/load-balancer/load-balancer/internal/server"
)

type LoadBalancer struct {
	Port    string
	ready   bool
	servers []node.Server
	algo    algo.RoutingAlgorithm
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		Port:  Config.Port,
		ready: true,
		servers: []node.Server{
			{Url: "http://localhost:3002", Active: true},
			{Url: "http://localhost:3003", Active: true},
			{Url: "http://localhost:3004", Active: true},
		},
	}
}

func (lb *LoadBalancer) setRoutingAlgorithm(algo algo.RoutingAlgorithm) {
	lb.algo = algo
}

func (lb *LoadBalancer) getServer() node.Server {
	return lb.algo.GetNextServer(lb.servers)
}

func (lb *LoadBalancer) ListenAndServe() error {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		status := "healthy"
		if !lb.ready {
			status = "unhealthy"
		}
		fmt.Fprintf(w, "Load balancer is %s", status)
	})
	http.HandleFunc("/", lb.handleConnection)
	log.Println("Load balancer is listening on", lb.Port)
	return http.ListenAndServe(lb.Port, nil)
}

func (lb *LoadBalancer) handleConnection(w http.ResponseWriter, r *http.Request) {
	server := lb.getServer()
	req, err := http.NewRequest(r.Method, server.Url+r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		log.Println("error in http.NewRequest", err)
		return
	}
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 3 * time.Second,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		log.Println("error in client.Do", err.Error())
		return
	}
	defer resp.Body.Close()
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, "Failed to write response body", http.StatusInternalServerError)
		log.Println("error in io.Copy", err)
		return
	}
}
