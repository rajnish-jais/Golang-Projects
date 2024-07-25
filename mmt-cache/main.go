// Makemy trip

/*
· get <key>

· put <key> <attributeKey1> <attributeValue1> <attributeKey2> <attributeValue2>....

· delete <key>

· search <attributeKey> <attributeValue>

· keys

· exit

{"sde_bootcamp": { "title": "SDE-Bootcamp", "price": 30000.00, "enrolled": false, "estimated_time": 30 }*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
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

func handleGet(cache *Cache, parts []string) {
	if len(parts) < 2 {
		fmt.Println("Usage: get <key>")
		return
	}
	key := parts[1]
	if entry, exists := cache.Get(key); exists {
		fmt.Printf("%s: %v\n", key, entry)
	} else {
		fmt.Printf("Key %s not found\n", key)
	}
}

func handlePut(cache *Cache, parts []string) {
	if len(parts) < 3 || len(parts)%2 != 0 {
		fmt.Println("Usage: put <key> <attributeKey1> <attributeValue1> <attributeKey2> <attributeValue2>...")
		return
	}
	key := parts[1]
	attributes := make(map[string]interface{})
	for i := 2; i < len(parts); i += 2 {
		attrKey := parts[i]
		attrValue := parts[i+1]
		attributes[attrKey] = attrValue
	}
	cache.Put(key, attributes)
	fmt.Printf("Put %s: %v\n", key, attributes)
}

func handleDelete(cache *Cache, parts []string) {
	if len(parts) < 2 {
		fmt.Println("Usage: delete <key>")
		return
	}
	key := parts[1]
	cache.Delete(key)
	fmt.Printf("Deleted %s\n", key)
}

func handleSearch(cache *Cache, parts []string) {
	if len(parts) < 3 {
		fmt.Println("Usage: search <attributeKey> <attributeValue>")
		return
	}
	attrKey := parts[1]
	attrValue := parts[2]
	keys := cache.Search(attrKey, attrValue)
	fmt.Printf("Keys with %s=%s: %v\n", attrKey, attrValue, keys)
}

func handleKeys(cache *Cache) {
	keys := cache.Keys()
	fmt.Println("Keys:", keys)
}

func main() {
	cache := NewCache()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]

		switch command {
		case "get":
			handleGet(cache, parts)
		case "put":
			handlePut(cache, parts)
		case "delete":
			handleDelete(cache, parts)
		case "search":
			handleSearch(cache, parts)
		case "keys":
			handleKeys(cache)
		case "exit":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Unknown command:", command)
		}
	}
}

// TestConcurrentAccess tests the cache for concurrent reads and writes.
func TestConcurrentAccess(t *testing.T) {
	cache := NewCache()
	var wg sync.WaitGroup

	// Number of concurrent goroutines
	numGoroutines := 100

	// Writer goroutines
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			value := CacheEntry{"value": i}
			cache.Put(key, value)
		}(i)
	}

	// Reader goroutines
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			cache.Get(key)
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Verify that all keys were written correctly
	for i := 0; i < numGoroutines; i++ {
		key := fmt.Sprintf("key%d", i)
		if entry, exists := cache.Get(key); !exists || entry["value"] != i {
			t.Errorf("Key %s has incorrect value %v", key, entry)
		}
	}
}

// TestConcurrentReadWrite tests the cache for concurrent reads and writes.
func TestConcurrentReadWrite(t *testing.T) {
	cache := NewCache()
	var wg sync.WaitGroup

	// Number of concurrent goroutines
	numGoroutines := 100

	// Writer goroutines
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			value := CacheEntry{"value": i}
			cache.Put(key, value)
		}(i)
	}

	// Reader and writer goroutines
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			cache.Get(key)
			cache.Put(key, CacheEntry{"value": i + 1})
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Verify that all keys were written correctly
	for i := 0; i < numGoroutines; i++ {
		key := fmt.Sprintf("key%d", i)
		if entry, exists := cache.Get(key); !exists || entry["value"] != i+1 {
			t.Errorf("Key %s has incorrect value %v", key, entry)
		}
	}
}

func TestMain(m *testing.M) {
	// Run tests
	os.Exit(m.Run())
}
