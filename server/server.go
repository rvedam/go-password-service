package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rvedam/go-password-service/hashlib"
)

type Stats struct {
	Total   int   `json:"total"`
	Average int64 `json:"average"`
}

type Server struct {
	totaltime     chan time.Duration
	stop          chan bool
	mux           *http.ServeMux
	statsrequest  chan int
	incomingstats chan Stats
}

type managerChannels struct {
	totalTimeChan     <-chan time.Duration
	statsRequestChan  <-chan int
	incomingStatsChan chan<- Stats
	stop              chan bool
}

func computeStats(mgr managerChannels) {
	totalPasswordRequests := 0
	var totalTime time.Duration
	for {
		select {
		case requestTime := <-mgr.totalTimeChan:
			totalPasswordRequests++
			totalTime += requestTime
		case <-mgr.statsRequestChan:
			var avg int64
			if totalPasswordRequests > 0 {
				avg = (totalTime.Nanoseconds() / (int64(totalPasswordRequests) * 1000))
			} else {
				avg = 0
			}
			mgr.incomingStatsChan <- Stats{Total: totalPasswordRequests, Average: avg}
		case <-mgr.stop:
			mgr.stop <- true
			return
		}
	}
}

// NewServer generates a new http server with our password service
func NewServer(stop chan bool) *Server {
	mux := http.NewServeMux()
	totalTimeChan := make(chan time.Duration, 100)
	statsRequestChan := make(chan int, 100)
	incomingStatsChan := make(chan Stats, 100)
	mgr := managerChannels{
		totalTimeChan:     totalTimeChan,
		statsRequestChan:  statsRequestChan,
		incomingStatsChan: incomingStatsChan,
		stop:              stop,
	}
	s := &Server{
		totaltime:     totalTimeChan,
		statsrequest:  statsRequestChan,
		incomingstats: incomingStatsChan,
		stop:          stop,
		mux:           mux,
	}

	go computeStats(mgr)
	s.mux.HandleFunc("/hash", s.computePasswordHash)
	s.mux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			stop <- true
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	s.mux.HandleFunc("/stats", s.getStats)
	return s
}

func (s *Server) computePasswordHash(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method == "POST" {
		time.Sleep(5 * time.Second)
		r.ParseForm()
		password := r.Form.Get("password")
		if password == "" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			hash := hashlib.Hash512AndEncodeBase64(password)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintln(w, strings.TrimSpace(hash))
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
	end := time.Since(start)
	s.totaltime <- end
}

func (s *Server) getStats(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		s.statsrequest <- 1
		data := <-s.incomingstats
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
