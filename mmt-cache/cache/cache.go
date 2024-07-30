package cache

import (
	"sync"
)

// CacheEntry represents a cache entry with dynamic attributes.
type CacheEntry map[string]interface{}

// Cache is a thread-safe cache structure.
type Cache struct {
	mu    sync.RWMutex
	store map[string]CacheEntry
}

// NewCache creates a new Cache instance.
func NewCache() *Cache {
	return &Cache{
		store: make(map[string]CacheEntry),
	}
}

// Get retrieves a cache entry by key.
func (c *Cache) Get(key string) (CacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, exists := c.store[key]
	return entry, exists
}

// Put adds or updates a cache entry by key and attributes.
func (c *Cache) Put(key string, attributes map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = attributes
}

// Delete removes a cache entry by key.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

// Search finds keys by attribute key and value.
func (c *Cache) Search(attributeKey string, attributeValue interface{}) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var result []string
	for key, entry := range c.store {
		if val, exists := entry[attributeKey]; exists && val == attributeValue {
			result = append(result, key)
		}
	}
	return result
}

// Keys returns all keys in the cache.
func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var keys []string
	for key := range c.store {
		keys = append(keys, key)
	}
	return keys
}
