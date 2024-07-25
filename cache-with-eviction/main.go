package main

import (
	"container/list"
	"fmt"
	"sync"
)

// CacheEvictionStrategy represents the eviction strategy for the cache.
type CacheEvictionStrategy int

const (
	LRU CacheEvictionStrategy = iota
	LFU
)

// Cache is the generic cache structure.
type Cache struct {
	capacity     int
	eviction     CacheEvictionStrategy
	accessList   *list.List
	accessMap    map[interface{}]*list.Element
	accessListMu sync.Mutex
}

type cacheItem struct {
	key      interface{}
	value    interface{}
	hitCount int
}

// NewCache creates a new Cache with the specified capacity and eviction strategy.
func NewCache(capacity int, eviction CacheEvictionStrategy) *Cache {
	return &Cache{
		capacity:   capacity,
		eviction:   eviction,
		accessList: list.New(),
		accessMap:  make(map[interface{}]*list.Element),
	}
}

// Get retrieves a value from the cache based on the key.
func (c *Cache) Get(key interface{}) (interface{}, bool) {
	c.accessListMu.Lock()
	defer c.accessListMu.Unlock()

	if elem, ok := c.accessMap[key]; ok {
		item := elem.Value.(*cacheItem)
		// Update hit count for eviction strategies
		item.hitCount++
		c.accessList.MoveToFront(elem)
		return item.value, true
	}
	return nil, false
}

// Set inserts a key-value pair into the cache.
func (c *Cache) Set(key, value interface{}) {
	c.accessListMu.Lock()
	defer c.accessListMu.Unlock()

	if elem, ok := c.accessMap[key]; ok {
		// Key already exists, update value and move to front
		item := elem.Value.(*cacheItem)
		item.value = value
		item.hitCount++
		c.accessList.MoveToFront(elem)
		return
	}

	if len(c.accessMap) >= c.capacity {
		c.evict()
	}

	item := &cacheItem{
		key:      key,
		value:    value,
		hitCount: 1,
	}
	elem := c.accessList.PushFront(item)
	c.accessMap[key] = elem
}

// SetEvictionPolicy changes the eviction policy of the cache during runtime.
func (c *Cache) SetEvictionPolicy(eviction CacheEvictionStrategy) {
	c.accessListMu.Lock()
	defer c.accessListMu.Unlock()

	c.eviction = eviction
}

// Evict removes the least recently or least frequently used item from the cache.
func (c *Cache) evict() {
	switch c.eviction {
	case LRU:
		c.evictLRU()
	case LFU:
		c.evictLFU()
	}
}

// evictLRU removes the least recently used item from the cache.
func (c *Cache) evictLRU() {
	if c.accessList.Len() > 0 {
		oldestItem := c.accessList.Back().Value.(*cacheItem)
		delete(c.accessMap, oldestItem.key)
		c.accessList.Remove(c.accessList.Back())
	}
}

// evictLFU removes the least frequently used item from the cache.
func (c *Cache) evictLFU() {
	var minHitCount = int(^uint(0) >> 1)
	var lfuItem *cacheItem

	for _, elem := range c.accessMap {
		item := elem.Value.(*cacheItem)
		if item.hitCount < minHitCount {
			minHitCount = item.hitCount
			lfuItem = item
		}
	}

	if lfuItem != nil {
		c.accessList.Remove(c.accessMap[lfuItem.key])
		delete(c.accessMap, lfuItem.key)
	}
}

func main() {
	// Example usage
	cache := NewCache(3, LRU)

	cache.Set("1", "one")
	cache.Set("2", "two")
	cache.Set("3", "three")

	val, exists := cache.Get("2")
	if exists {
		fmt.Println("Value for key '2':", val)
	}
	cache.Set("4", "three")

	// Print the current cache contents
	for key, elem := range cache.accessMap {
		item := elem.Value.(*cacheItem)
		fmt.Printf("Key: %v, Value: %v , Capacity: %v \n", key, item.value, len(cache.accessMap))
	}
	println()

	//Change the eviction policy during runtime
	cache.SetEvictionPolicy(LFU)

	// Now the cache should evict keys based on the new LFU eviction policy
	cache.Set("4", "four")
	cache.Set("4", "four")
	cache.Set("3", "three")
	cache.Set("4", "four")
	cache.Set("4", "four")
	cache.Set("3", "three")
	cache.Set("1", "one")
	cache.Set("3", "three")
	cache.Set("2", "one")
	cache.Set("2", "three")
	cache.Set("2", "one")
	cache.Set("2", "three")
	cache.Set("2", "one")

	// Print the current cache contents
	for key, elem := range cache.accessMap {
		item := elem.Value.(*cacheItem)
		fmt.Printf("Key: %v, Value: %v , Capacity: %v \n", key, item.value, len(cache.accessMap))
	}
}
