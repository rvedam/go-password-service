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

func ComputePasswordHash(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("Another hash request...")
		time.Sleep(5 * time.Second)
		r.ParseForm()
		fmt.Fprintln(w, hashlib.Hash512AndEncodeBase64(r.Form.Get("password")))
	}
}
