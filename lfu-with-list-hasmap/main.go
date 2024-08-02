package main

import (
	"container/list"
	"fmt"
)

// LFUCache represents the LFU cache structure
type LFUCache struct {
	capacity int
	minFreq  int
	nodes    map[int]*Node
	freq     map[int]*list.List
}

// Node represents a single node in the cache
type Node struct {
	key       int
	value     int
	frequency int
	element   *list.Element
}

// NewLFUCache initializes a new LFU cache
func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{
		capacity: capacity,
		nodes:    make(map[int]*Node),
		freq:     make(map[int]*list.List),
	}
}

// Get retrieves the value of the key if it exists in the cache
func (c *LFUCache) Get(key int) int {
	if node, exists := c.nodes[key]; exists {
		c.incrementFrequency(node)
		return node.value
	}
	return -1
}

// Put sets or inserts the value if the key is not already present
func (c *LFUCache) Put(key int, value int) {
	if c.capacity == 0 {
		return
	}

	if node, exists := c.nodes[key]; exists {
		node.value = value
		c.incrementFrequency(node)
	} else {
		if len(c.nodes) >= c.capacity {
			c.evict()
		}
		node := &Node{key: key, value: value, frequency: 1}
		if c.freq[1] == nil {
			c.freq[1] = list.New()
		}
		node.element = c.freq[1].PushFront(node)
		c.nodes[key] = node
		c.minFreq = 1
	}
}

// incrementFrequency increments the frequency of the node
func (c *LFUCache) incrementFrequency(node *Node) {
	freq := node.frequency
	c.freq[freq].Remove(node.element)

	if c.freq[freq].Len() == 0 {
		delete(c.freq, freq)
		if c.minFreq == freq {
			c.minFreq++
		}
	}

	node.frequency++
	if c.freq[node.frequency] == nil {
		c.freq[node.frequency] = list.New()
	}
	node.element = c.freq[node.frequency].PushFront(node)
}

// evict removes the least frequently used item from the cache
func (c *LFUCache) evict() {
	list := c.freq[c.minFreq]
	node := list.Back().Value.(*Node)
	list.Remove(node.element)
	delete(c.nodes, node.key)
	if list.Len() == 0 {
		delete(c.freq, c.minFreq)
	}
}

func main() {
	cache := NewLFUCache(2)
	cache.Put(1, 1)
	cache.Put(2, 2)
	fmt.Println(cache.Get(1)) // returns 1
	cache.Put(3, 3)           // evicts key 2
	fmt.Println(cache.Get(2)) // returns -1 (not found)
	fmt.Println(cache.Get(3)) // returns 3
	cache.Put(4, 4)           // evicts key 1
	fmt.Println(cache.Get(1)) // returns -1 (not found)
	fmt.Println(cache.Get(3)) // returns 3
	fmt.Println(cache.Get(4)) // returns 4
}
