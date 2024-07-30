package commands

import (
	"fmt"
	"mmt-cache/cache"
)

func HandleDelete(cache *cache.Cache, parts []string) {
	if len(parts) < 2 {
		fmt.Println("Usage: delete <key>")
		return
	}
	key := parts[1]
	cache.Delete(key)
	fmt.Printf("Deleted %s\n", key)
}
