package main

type LoadBalancer struct {
	servers []Server
	algo    string
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		servers: []Server{
			{url: "http://localhost:3001", active: true},
			{url: "http://localhost:3002", active: true},
		},
		algo: "www",
	}
}

func (lb *LoadBalancer) getServer() {
	// return lb.algo.getNextServer()

}
