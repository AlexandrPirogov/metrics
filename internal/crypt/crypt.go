// Package crypt contains functions for Hash
package crypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

// Hash hashes given string with given key
func Hash(toHash, key string) string {
	if key == "" {
		return ""
	}
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(toHash))
	res := h.Sum(nil)
	return fmt.Sprintf("%x", res)
}
