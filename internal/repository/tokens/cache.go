package tokens

import (
	"sync"
)

type MemoryTokenCache struct {
	tokens map[string]struct{}
	mx     *sync.RWMutex
}

func NewTokenCache() *MemoryTokenCache {
	return &MemoryTokenCache{
		tokens: map[string]struct{}{},
		mx:     new(sync.RWMutex),
	}
}

func (m *MemoryTokenCache) SetRecalculationChan(deleted chan []string) {
	go func() {
		for {
			removeTokens := <-deleted
			m.mx.Lock()
			for _, token := range removeTokens {
				delete(m.tokens, token)
			}
			m.mx.Unlock()
		}
	}()
}

func (m *MemoryTokenCache) Exists(token string) bool {
	m.mx.RLock()
	defer m.mx.RUnlock()

	_, exists := m.tokens[token]
	return exists
}

func (m *MemoryTokenCache) Store(token string) {
	if m.Exists(token) {
		return
	}
	m.mx.Lock()
	m.tokens[token] = struct{}{}
	m.mx.Unlock()
}
