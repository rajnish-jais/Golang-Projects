package commands

import (
	"fmt"
	"mmt-cache/cache"
)

func HandleSearch(cache *cache.Cache, parts []string) {
	if len(parts) < 3 {
		fmt.Println("Usage: search <attributeKey> <attributeValue>")
		return
	}
	attrKey := parts[1]
	attrValue := parts[2]
	keys := cache.Search(attrKey, attrValue)
	fmt.Printf("Keys with %s=%s: %v\n", attrKey, attrValue, keys)
}
