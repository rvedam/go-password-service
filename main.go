package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rvedam/go-password-service/hashlib"
)

func computePasswordHash(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {

	} else {
		r.ParseForm()
		fmt.Fprintln(w, hashlib.Hash512AndEncodeBase64(r.Form.Get("password")))
		time.Sleep(5 * time.Second)
	}
}

func main() {
	http.HandleFunc("/hash", computePasswordHash)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
