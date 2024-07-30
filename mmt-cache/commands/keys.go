package commands

import (
	"fmt"
	"mmt-cache/cache"
)

func HandleKeys(cache *cache.Cache) {
	keys := cache.Keys()
	fmt.Println("Keys:", keys)
}
