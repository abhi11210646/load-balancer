package main

type RoutingAlgorithm interface {
	getNextServer(servers []Server) Server
}

type RoundRobin struct {
	current_index int
}

func (r *RoundRobin) getNextServer(servers []Server) Server {
	r.current_index = (r.current_index + 1) % len(servers)
	return servers[r.current_index]
}
