package commands

import (
	"fmt"
	"mmt-cache/cache"
)

func HandleGet(cache *cache.Cache, parts []string) {
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
