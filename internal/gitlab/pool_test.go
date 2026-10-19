package gitlab

import (
	"context"
	"testing"
)

func TestClientPoolIsolation(t *testing.T) {
	pool := NewClientPool()
	c1 := pool.Get("http://gitlab.example.com", "token-a")
	c2 := pool.Get("http://gitlab.example.com", "token-b")
	c1Again := pool.Get("http://gitlab.example.com", "token-a")
	if c1 == c2 {
		t.Fatal("expected different clients for different tokens")
	}
	if c1 != c1Again {
		t.Fatal("expected same client instance for same url+token")
	}
}

func TestClientFromContextMissingToken(t *testing.T) {
	pool := NewClientPool()
	ctx := ContextWithCreds(context.Background(), Creds{BaseURL: "http://x", Token: ""})
	if _, err := pool.ClientFromContext(ctx); err == nil {
		t.Fatal("expected error for empty token")
	}
}
