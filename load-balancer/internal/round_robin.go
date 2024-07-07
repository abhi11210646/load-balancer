package lb

type RoutingAlgorithm interface {
	GetNextServer(servers []*Server) *Server
}

type RoundRobin struct {
	current_index int
}

func (r *RoundRobin) GetNextServer(servers []*Server) *Server {
	server := servers[r.current_index]
	r.current_index = (r.current_index + 1) % len(servers)
	return server
}
