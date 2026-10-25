package helpers

import (
	"sync"
	"time"
)

var revokedTokens = struct {
	sync.Mutex
	expiresAt map[string]time.Time
}{expiresAt: make(map[string]time.Time)}

func RevokeToken(tokenID string, expiresAt time.Time) {
	revokedTokens.Lock()
	defer revokedTokens.Unlock()

	now := time.Now()
	for id, expiry := range revokedTokens.expiresAt {
		if !expiry.After(now) {
			delete(revokedTokens.expiresAt, id)
		}
	}
	revokedTokens.expiresAt[tokenID] = expiresAt
}

func IsTokenRevoked(tokenID string) bool {
	revokedTokens.Lock()
	defer revokedTokens.Unlock()
	if expiresAt, ok := revokedTokens.expiresAt[tokenID]; ok && expiresAt.After(time.Now()) {
		return true
	}
	delete(revokedTokens.expiresAt, tokenID)
	return false
}
