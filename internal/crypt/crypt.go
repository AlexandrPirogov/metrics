package crypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"memtracker/internal/config/agent"
)

func Hash(toHash, key string) string {
	cfg := agent.ClientCfg
	if cfg.Hash == "" {
		return ""
	}
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(toHash))
	return string(h.Sum(nil))

}
