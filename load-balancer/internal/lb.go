package lb

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type LoadBalancer struct {
	port              int
	ready             bool
	servers           []*Server
	algo              RoutingAlgorithm
	heartBeatInterval time.Duration
	sync.Mutex
}

func NewLoadBalancer(port int) *LoadBalancer {
	return &LoadBalancer{
		port:              port,
		ready:             true,
		heartBeatInterval: 5,
		servers:           []*Server{},
	}
}

func (lb *LoadBalancer) SetRoutingAlgorithm(algo RoutingAlgorithm) {
	lb.algo = algo
}
func (lb *LoadBalancer) markUnhealthy() {
	lb.Lock()
	defer lb.Unlock()
	lb.ready = false
}
func (lb *LoadBalancer) markHealthy() {
	lb.Lock()
	defer lb.Unlock()
	lb.ready = true
}
func (lb *LoadBalancer) SetServers(nodes []string) {
	for _, n := range nodes {
		lb.servers = append(lb.servers, NewServer(n))
	}
}

func (lb *LoadBalancer) getServer() (*Server, error) {
	c := 0
	for c < len(lb.servers) {
		server := lb.algo.GetNextServer(lb.servers)
		if server.Active {
			return server, nil
		}
		c += 1
	}
	lb.markUnhealthy()
	return nil, errors.New("no server is available to serve request")
}

func (lb *LoadBalancer) healthCheck() {
	wg := &sync.WaitGroup{}
	for {
		wg.Add(len(lb.servers))
		for _, server := range lb.servers {
			go server.HealthCheck(wg)
		}
		wg.Wait()
		Inactive := true
		for _, server := range lb.servers {
			if server.Active {
				Inactive = false
				lb.markHealthy()
				break
			}
		}
		if Inactive {
			lb.markUnhealthy()
		}
		time.Sleep(lb.heartBeatInterval * time.Second)
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
	log.Println("Load balancer is listening on", lb.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", lb.port), nil)
}

func (lb *LoadBalancer) handleConnection(w http.ResponseWriter, r *http.Request) {
	if !lb.ready {
		http.Error(w, "Load balancer is offlie", http.StatusServiceUnavailable)
		return
	}
	server, err := lb.getServer()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
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
		// server.MarkInactive()
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
