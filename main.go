package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	// Step 1: Define command line flags
	var (
		port       = flag.Int("port", 8080, "Port to run the proxy server on")
		origin     = flag.String("origin", "", "Origin server URL to proxy requests to")
		clearCache = flag.Bool("clear-cache", false, "Clear the cache and exit")
	)
	flag.Parse()

	// Step 2: Handle cache clearing
	if *clearCache {
		if err := clearCacheData(); err != nil {
			log.Fatalf("Failed to clear cache: %v", err)
		}
		fmt.Println("Cache cleared successfully")
		return
	}

	// Step 3: Validate required parameters
	if *origin == "" {
		log.Fatal("Origin server URL is required. Use --origin flag")
	}

	// Step 4: Create and start proxy server
	server := NewProxyServer(*origin)

	fmt.Printf("Starting caching proxy server on port %d\n", *port)
	fmt.Printf("Proxying requests to: %s\n", *origin)

	if err := server.Start(*port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
