package repository_implement

import (
	"context"
	"fmt"
	"sync"
	"time"

	repository_interface "github.com/vando2108/sandbox_service/internal/app/sandbox/repository/interface"
)

type MemNonceCache struct {
	mem   map[string]string
	mutex sync.Mutex
}

func NewMemNonceCache() repository_interface.NonceCacheRepository {
	return &MemNonceCache{
		mem: make(map[string]string),
	}
}

func (r *MemNonceCache) SetNonce(ctx context.Context, publickey string, nonce string, expireTime time.Duration) error {
	// avoid expireTime
	r.mutex.Lock()
	r.mem[publickey] = nonce
	r.mutex.Unlock()

	return nil
}

func (r *MemNonceCache) GetNonce(ctx context.Context, publickey string) (string, error) {
	if nonce, exists := r.mem[publickey]; exists {
		return nonce, nil
	} else {
		return "", fmt.Errorf("publickey not found")
	}
}
