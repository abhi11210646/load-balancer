package algo

import (
	node "github.com/load-balancer/load-balancer/internal/server"
)

type RoutingAlgorithm interface {
	GetNextServer(servers []node.Server) node.Server
}

type RoundRobin struct {
	current_index int
}

func (r *RoundRobin) GetNextServer(servers []node.Server) node.Server {
	server := servers[r.current_index]
	r.current_index = (r.current_index + 1) % len(servers)
	return server
}
