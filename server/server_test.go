package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestComputePasswordHash(t *testing.T) {
	req, err := http.NewRequest("POST", "/hash", nil)
	req.Header.Set("Content-Type", "application/x-www-form-url-encoded")
	req.Form = url.Values{}
	req.Form.Set("password", "angryMonkey")

	if err != nil {
		t.Errorf("Request Generation Error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	stop := make(chan bool, 1)
	s := NewServer(stop)
	s.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}
	want := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
	got := strings.TrimSpace(rr.Body.String())
	if got != want {
		t.Errorf("Incorrect response: got %v, want %v", got, want)
	}
	stop <- true
}

func TestComputeStats(t *testing.T) {
	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Errorf("Request Generation Error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	stop := make(chan bool, 1)
	s := NewServer(stop)
	s.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v, want %v\n", status, http.StatusOK)
	}
}
