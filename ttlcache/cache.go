package ttlcache

import (
	"context"
	"fmt"

	"github.com/graph-gophers/dataloader/v7"
	cache "github.com/patrickmn/go-cache"
)

// TTLCache implements the dataloader.Cache interface
type TTLCache[K comparable, V any] struct {
	c *cache.Cache
}

// Get gets a value from the cache
func (c *TTLCache[K, V]) Get(_ context.Context, key K) (dataloader.Thunk[V], bool) {
	k := fmt.Sprintf("%v", key) // convert the key to string because the underlying library doesn't support Generics yet
	v, ok := c.c.Get(k)
	if ok {
		return v.(dataloader.Thunk[V]), ok
	}
	return nil, ok
}

// Set sets a value in the cache
func (c *TTLCache[K, V]) Set(_ context.Context, key K, value dataloader.Thunk[V]) {
	k := fmt.Sprintf("%v", key) // convert the key to string because the underlying library doesn't support Generics yet
	c.c.SetDefault(k, value)
}

// Delete deletes and item in the cache
func (c *TTLCache[K, V]) Delete(_ context.Context, key K) bool {
	k := fmt.Sprintf("%v", key) // convert the key to string because the underlying library doesn't support Generics yet
	if _, found := c.c.Get(k); found {
		c.c.Delete(k)
		return true
	}
	return false
}

// Clear clears the cache
func (c *TTLCache[K, V]) Clear() {
	c.c.Flush()
}
