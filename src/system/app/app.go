package app

import (
	"log"
	"net/http"
)

// Server object
type Server struct {
	port string
}

// NewServer returns a new Server
func NewServer() Server {
	return Server{}
}

// Init all vals
func (s *Server) Init(port string) {
	log.Println("Initializing server...")
	s.port = ":" + port
}

// Start the server
func (s *Server) Start() {
	log.Println("Starting server on port " + s.port)
	http.ListenAndServe(s.port, nil)
}
