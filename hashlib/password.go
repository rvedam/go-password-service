package hashlib

import (
	"crypto/sha512"
	"encoding/base64"
)

//
// Hash512AndEncodeBase64 takes a password as input and returns the base64 encoded version of the hash.
//
func Hash512AndEncodeBase64(passwd string) string {
	hasher := sha512.New()
	hasher.Write([]byte(passwd))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}
