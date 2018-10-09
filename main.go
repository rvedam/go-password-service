package main

import (
	"fmt"
	"os"

	"github.com/rvedam/go-password-service/hashlib"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("USAGE: ./go-password-service <password>")
		return
	}
	fmt.Println(hashlib.Hash512AndEncodeBase64(os.Args[1]))
}
