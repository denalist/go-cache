package main

// package main is the executable entry point. All .go files in the same directory sharing
// package main form one namespace — functions in api.go are callable here without importing.
// Named packages (e.g. package cache) are libraries; capitalised identifiers are exported,
// lowercase are unexported. Dependencies flow one way: main → cache, never the reverse.

import (
	"go-cache/cache"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	dc := cache.NewSharedCache(4, 100)

	// metrics := c.GetMetrics()
	// fmt.Printf("hits=%d misses=%d evictions=%d size=%d\n",
	// 	metrics.Hits,
	// 	metrics.Misses,
	// 	metrics.Evictions,
	// 	metrics.Size,
	// )

	srv, err := cache.NewServer(":8080", dc)

	if err != nil {
		log.Fatal(err)
	}
	go srv.Serve()
	log.Println("cache server listening on :8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // blocks until Ctrl+C or kill signal

	log.Println("shutting down")
	srv.Close()
}
