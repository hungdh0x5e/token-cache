package tokencache

import (
	"context"
	"sync"

	"golang.org/x/oauth2"
)

// TokenGetter is the interface that wraps the FetchToken function.
type TokenGetter interface {
	// FetchToken syncs the token from the remote
	FetchToken(ctx context.Context) (*oauth2.Token, error)
}

func NewTokenCache(getter TokenGetter) *TokenCache {
	return &TokenCache{getter: getter}
}

type TokenCache struct {
	getter   TokenGetter
	inflight *inflight
	token    *oauth2.Token
	mu       sync.Mutex // guard token
}

// GetToken return token from cache or fetch from remote.
func (c *TokenCache) GetToken(ctx context.Context) (*oauth2.Token, error) {
	token := c.getFromCache()
	if token != nil && token.Valid() {
		return token, nil
	}

	return c.fetchFromRemote(ctx)
}

func (c *TokenCache) getFromCache() *oauth2.Token {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.token
}

// fetchFromRemote syncs the token from the remote, records the values in the
// cache, and returns the token.
func (c *TokenCache) fetchFromRemote(ctx context.Context) (*oauth2.Token, error) {
	// Need to lock to inspect the inflight request field.
	c.mu.Lock()
	// If there's not a current inflight request, create one.
	if c.inflight == nil {
		c.inflight = newInflight()

		// This goroutine has exclusive ownership over the current inflight
		// request. It releases the resource by nil'ing the inflight field
		// once the goroutine is done.
		go func() {
			// Sync token and finish inflight when that's done.
			token, err := c.getter.FetchToken(ctx)

			c.inflight.done(token, err)

			// Lock to update the token and indicate that there is no longer an
			// inflight request.
			c.mu.Lock()
			defer c.mu.Unlock()

			if err == nil {
				c.token = token
			} else { // Clear token
				c.token = nil
			}

			// Free inflight so a different request can run.
			c.inflight = nil
		}()
	}
	inflight := c.inflight
	c.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-inflight.wait():
		return inflight.result()
	}
}
