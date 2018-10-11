package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rvedam/go-password-service/hashlib"
)

type requestStats struct {
	totaltime     time.Duration
	totalrequests int
}

// Server holds our routing table (needed by http.server)
// the main purpose of this object is to eventually hold statistics for
// the /stat call
type Server struct {
	mux  *http.ServeMux
	stop chan bool
}

// New creates a new server with our custom HTTP Handlers
func New(stop chan bool) *Server {
	s := &Server{mux: http.NewServeMux(), stop: stop}
	s.mux.HandleFunc("/hash", s.computePasswordHash)
	s.mux.HandleFunc("/shutdown", s.shutdown)
	return s
}

func (s *Server) shutdown(w http.ResponseWriter, r *http.Request) {
	s.stop <- true
}

func (s *Server) computePasswordHash(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("Another hash request...")
		time.Sleep(5 * time.Second)
		r.ParseForm()
		fmt.Fprintln(w, hashlib.Hash512AndEncodeBase64(r.Form.Get("password")))
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
