package commands

import (
	"fmt"
	"mmt-cache/cache"
)

func HandlePut(cache *cache.Cache, parts []string) {
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
