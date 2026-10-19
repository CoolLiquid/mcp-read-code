package gitlab

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
)

// ClientPool returns isolated GitLab clients keyed by URL + token hash.
type ClientPool struct {
	mu      sync.Mutex
	clients map[string]*Client
}

func NewClientPool() *ClientPool {
	return &ClientPool{clients: make(map[string]*Client)}
}

func (p *ClientPool) poolKey(baseURL, token string) string {
	sum := sha256.Sum256([]byte(baseURL + "\x00" + token))
	return hex.EncodeToString(sum[:])
}

// Get returns a cached client for the given GitLab URL and token.
func (p *ClientPool) Get(baseURL, token string) *Client {
	key := p.poolKey(baseURL, token)
	p.mu.Lock()
	defer p.mu.Unlock()
	if c, ok := p.clients[key]; ok {
		return c
	}
	c := NewClient(baseURL, token)
	p.clients[key] = c
	return c
}

// ClientFromContext resolves credentials on ctx and returns the matching client.
func (p *ClientPool) ClientFromContext(ctx context.Context) (*Client, error) {
	creds, ok := CredsFromContext(ctx)
	if !ok || creds.Token == "" {
		return nil, fmt.Errorf("GitLab token missing: set X-Gitlab-Token or Authorization Bearer header")
	}
	if creds.BaseURL == "" {
		return nil, fmt.Errorf("GitLab URL missing: set X-Gitlab-Url header or GITLAB_URL on server")
	}
	return p.Get(creds.BaseURL, creds.Token), nil
}
