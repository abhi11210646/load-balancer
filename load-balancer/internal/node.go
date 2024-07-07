package lb

import (
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	Url    string
	Active bool
	mu     sync.Mutex
}

func NewServer(url string) *Server {
	return &Server{
		Url:    url,
		Active: true,
	}
}

func (s *Server) HealthCheck(wg *sync.WaitGroup) {
	defer wg.Done()
	req, err := http.NewRequest("GET", s.Url+"/health", nil)
	if err != nil {
		log.Println("error in http.NewRequest", err)
		return
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		s.MarkInactive()
		return
	}
	defer res.Body.Close()

	if _, err := io.ReadAll(res.Body); err != nil {
		log.Println("error in reading response", err)
	}

	if res.StatusCode == http.StatusOK {
		s.MarkActive()
		return
	}
	s.MarkInactive()
}

func (s *Server) MarkActive() {
	s.mu.Lock()
	s.Active = true
	s.mu.Unlock()
}
func (s *Server) MarkInactive() {
	s.mu.Lock()
	s.Active = false
	s.mu.Unlock()
}
