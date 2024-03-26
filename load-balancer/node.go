package main

type Server struct {
	url    string
	active bool
}

func (s *Server) markActive() {
	s.active = true
}
func (s *Server) markInactive() {
	s.active = false
}
