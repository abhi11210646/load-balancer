package node

type Server struct {
	Url    string
	Active bool
}

func (s *Server) markActive() {
	s.Active = true
}
func (s *Server) markInactive() {
	s.Active = false
}
