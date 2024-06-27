package node

import (
	"fmt"
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
	fmt.Println(s.Url)
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
		s.markInactive()
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		s.markActive()
		return
	}
	s.markInactive()
}

func (s *Server) markActive() {
	s.mu.Lock()
	s.Active = true
	s.mu.Unlock()
}
func (s *Server) markInactive() {
	s.mu.Lock()
	s.Active = false
	s.mu.Unlock()
}
