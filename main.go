package main

// package main is the executable entry point. All .go files in the same directory sharing
// package main form one namespace — functions in api.go are callable here without importing.
// Named packages (e.g. package cache) are libraries; capitalised identifiers are exported,
// lowercase are unexported. Dependencies flow one way: main → cache, never the reverse.

import (
	"fmt"
	"go-cache/cache"

	"github.com/gin-gonic/gin"
)

func main() {

	c := cache.NewCache(100)

	metrics := c.GetMetrics()
	fmt.Printf("hits=%d misses=%d evictions=%d size=%d\n",
		metrics.Hits,
		metrics.Misses,
		metrics.Evictions,
		metrics.Size,
	)

	router := gin.Default()

	router.POST("/cache/:key", SetCacheValue(c))      // Create
	router.PUT("/cache/:key", SetCacheValue(c))       // Update
	router.GET("/cache/:key", GetCacheValue(c))       // Read
	router.DELETE("/cache/:key", DeleteCacheValue(c)) // Delete

	router.Run(":8080")

}
