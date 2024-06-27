package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	algo "github.com/load-balancer/load-balancer/internal/algorithm"
	node "github.com/load-balancer/load-balancer/internal/server"
)

type LoadBalancer struct {
	Port    string
	ready   bool
	servers []*node.Server
	algo    algo.RoutingAlgorithm
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		Port:  Config.Port,
		ready: true,
		servers: []*node.Server{
			node.NewServer("http://localhost:3002"),
			node.NewServer("http://localhost:3003"),
			node.NewServer("http://localhost:3004"),
		},
	}
}

func (lb *LoadBalancer) setRoutingAlgorithm(algo algo.RoutingAlgorithm) {
	lb.algo = algo
}

func (lb *LoadBalancer) getServer() *node.Server {
	c := 0
	for c < len(lb.servers) {
		server := lb.algo.GetNextServer(lb.servers)
		if server.Active {
			return server
		}
		c += 1
	}
	lb.ready = false
	return nil
}

func (lb *LoadBalancer) healthCheck() {
	wg := &sync.WaitGroup{}
	for {
		for _, server := range lb.servers {
			wg.Add(1)
			go server.HealthCheck(wg)
		}
		Inactive := true
		for _, server := range lb.servers {
			if server.Active {
				Inactive = false
				lb.ready = true
				break
			}
		}
		if Inactive {
			lb.ready = false
		}
		wg.Wait()
		time.Sleep(5 * time.Second)
	}
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
	go lb.healthCheck()
	log.Println("Load balancer is listening on", lb.Port)
	return http.ListenAndServe(lb.Port, nil)
}

func (lb *LoadBalancer) handleConnection(w http.ResponseWriter, r *http.Request) {
	if !lb.ready {
		http.Error(w, "Load balancer is offlie", http.StatusServiceUnavailable)
		return
	}
	server := lb.getServer()
	if server == nil {
		http.Error(w, "No server is available to serve request", http.StatusServiceUnavailable)
		return
	}
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
		Timeout: 5 * time.Second,
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
