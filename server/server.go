package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rvedam/go-password-service/hashlib"
)

type Stats struct {
	Total   int
	Average int64
}

type Server struct {
	totalrequests  chan int
	totaltime      chan int64
	stop           chan bool
	mux            *http.ServeMux
	stats_request  chan int
	incoming_stats chan Stats
}

func NewServer(stop chan bool) *Server {
	mux := http.NewServeMux()
	total_request_chan := make(chan int, 100)
	total_time_chan := make(chan int64, 100)
	stats_request_chan := make(chan int, 100)
	incoming_stats_chan := make(chan Stats, 100)

	s := &Server{
		totalrequests:  total_request_chan,
		totaltime:      total_time_chan,
		stats_request:  stats_request_chan,
		incoming_stats: incoming_stats_chan,
		stop:           stop,
		mux:            mux,
	}

	go func() {
		total_passwd_requests := 0
		totalTime := int64(0)
		for {
			select {
			case c := <-total_request_chan:
				fmt.Println(c)
				total_passwd_requests += c
			case request_time := <-total_time_chan:
				totalTime += request_time
			case <-stats_request_chan:
				incoming_stats_chan <- Stats{Total: total_passwd_requests, Average: totalTime}
			case <-stop:
				return
			}
		}
	}()
	s.mux.HandleFunc("/hash", s.computePasswordHash)
	s.mux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, http.StatusOK)
		stop <- true
	})
	s.mux.HandleFunc("/stats", s.computeStats)
	return s
}

func (s *Server) computePasswordHash(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		s.totalrequests <- 1
		fmt.Println("Another hash request...")
		time.Sleep(5 * time.Second)
		r.ParseForm()
		fmt.Fprintln(w, hashlib.Hash512AndEncodeBase64(r.Form.Get("password")))
	}
}

func (s *Server) computeStats(w http.ResponseWriter, r *http.Request) {
	s.stats_request <- 1
	data := <-s.incoming_stats
	fmt.Println(data)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
