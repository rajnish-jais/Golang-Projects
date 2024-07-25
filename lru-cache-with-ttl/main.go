package main

import (
	"container/list"
	"fmt"
	"time"
)

type CacheEntry struct {
	key        int
	value      int
	expiryTime time.Time
}

type IteratorsContainer struct {
	cacheIterator           *list.Element
	keyIteratorInTimeBucket *list.Element
}

type LRUCache struct {
	capacity        int
	cache           *list.List
	timeBuckets     map[time.Time]*list.List
	keyIteratorsMap map[int]IteratorsContainer
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity:        capacity,
		cache:           list.New(),
		timeBuckets:     make(map[time.Time]*list.List),
		keyIteratorsMap: make(map[int]IteratorsContainer),
	}
}

func (c *LRUCache) updateCache(key, value int, ttl time.Time) {
	temp := c.keyIteratorsMap[key]

	//currentValue := temp.cacheIterator.Value.(CacheEntry).value
	currentExpiryTime := temp.cacheIterator.Value.(CacheEntry).expiryTime
	newExpiryTime := ttl

	c.cache.Remove(temp.cacheIterator)
	c.cache.PushFront(CacheEntry{key, value, newExpiryTime})

	c.timeBuckets[currentExpiryTime].Remove(temp.keyIteratorInTimeBucket)
	if c.timeBuckets[currentExpiryTime].Len() == 0 {
		delete(c.timeBuckets, currentExpiryTime)
	}

	if _, exists := c.timeBuckets[newExpiryTime]; !exists {
		c.timeBuckets[newExpiryTime] = list.New()
	}
	c.timeBuckets[newExpiryTime].PushFront(key)

	c.keyIteratorsMap[key] = IteratorsContainer{cacheIterator: c.cache.Front(), keyIteratorInTimeBucket: c.timeBuckets[newExpiryTime].Front()}
}

func (c *LRUCache) get(key int) int {

	if temp, exists := c.keyIteratorsMap[key]; exists {
		currentTime := time.Now()
		expiryTime := temp.cacheIterator.Value.(CacheEntry).expiryTime
		value := temp.cacheIterator.Value.(CacheEntry).value

		if expiryTime.Before(currentTime) {
			return -1
		}

		c.updateCache(key, value, expiryTime)
		return value
	}
	return -1
}

func (c *LRUCache) evictLRU() {
	evictionCandidate := c.cache.Back().Value.(CacheEntry)
	temp := c.keyIteratorsMap[evictionCandidate.key]

	c.timeBuckets[evictionCandidate.expiryTime].Remove(temp.keyIteratorInTimeBucket)
	if c.timeBuckets[evictionCandidate.expiryTime].Len() == 0 {
		delete(c.timeBuckets, evictionCandidate.expiryTime)
	}

	c.cache.Remove(temp.cacheIterator)
	//delete(c.timeBuckets,evictionCandidate.expiryTime)
	delete(c.keyIteratorsMap, evictionCandidate.key)
}

func (c *LRUCache) put(key, value int, ttl time.Time) {

	if _, exists := c.keyIteratorsMap[key]; exists {
		fmt.Println("Key exists. Update cache")
		c.updateCache(key, value, ttl)
		return
	}

	if len(c.keyIteratorsMap) == c.capacity {
		c.evictLRU()
	}

	expiryTime := ttl
	c.cache.PushFront(CacheEntry{key, value, expiryTime})

	if _, exists := c.timeBuckets[expiryTime]; !exists {
		c.timeBuckets[expiryTime] = list.New()
	}
	c.timeBuckets[expiryTime].PushFront(key)

	c.keyIteratorsMap[key] = IteratorsContainer{cacheIterator: c.cache.Front(), keyIteratorInTimeBucket: c.timeBuckets[expiryTime].Front()}
}

func main() {
	cache := NewLRUCache(3)

	cache.put(1, 1, time.Now().Add(5*time.Second))

	val := cache.get(1)
	fmt.Println(val)
	val = cache.get(2)
	fmt.Println(val)

	time.Sleep(6 * time.Second)
	val = cache.get(1)
	fmt.Println(val)
	fmt.Println(cache.cache.Len())
	fmt.Println(len(cache.keyIteratorsMap), len(cache.timeBuckets))

}
