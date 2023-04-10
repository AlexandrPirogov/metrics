package crypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"memtracker/internal/config/agent"
)

func Hash(toHash string) string {
	cfg := agent.ClientCfg
	if cfg.Hash == "" {
		return ""
	}
	key := cfg.Hash
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(toHash))
	return string(h.Sum(nil))

}
