package main

import (
	"container/heap"
	"fmt"
)

// LFUCache represents the LFU cache structure
type LFUCache struct {
	capacity int
	nodes    map[int]*Node
	freq     *MinHeap
}

// Node represents a single node in the cache
type Node struct {
	key       int
	value     int
	frequency int
	index     int // The index of the item in the heap
}

// MinHeap represents a min-heap of nodes based on their frequency
type MinHeap []*Node

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].frequency < h[j].frequency }
func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *MinHeap) Push(x interface{}) {
	node := x.(*Node)
	node.index = len(*h)
	*h = append(*h, node)
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	node := old[n-1]
	node.index = -1 // for safety
	*h = old[0 : n-1]
	return node
}

// NewLFUCache initializes a new LFU cache
func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{
		capacity: capacity,
		nodes:    make(map[int]*Node),
		freq:     &MinHeap{},
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
		heap.Push(c.freq, node)
		c.nodes[key] = node
	}
}

// incrementFrequency increments the frequency of the node
func (c *LFUCache) incrementFrequency(node *Node) {
	node.frequency++
	heap.Fix(c.freq, node.index)
}

// evict removes the least frequently used item from the cache
func (c *LFUCache) evict() {
	node := heap.Pop(c.freq).(*Node)
	delete(c.nodes, node.key)
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
