
# Caching Proxy

A high-performance HTTP caching proxy server written in Go that sits between clients and origin servers to cache responses and reduce backend load.

## Features

- **In-memory caching**: Fast response caching using Go's built-in data structures
- **Cache hit/miss headers**: Adds `X-Cache: HIT/MISS` headers to responses
- **Configurable origin**: Proxy requests to any HTTP server
- **Cache management**: Clear cache with command-line flag
- **Multiple HTTP methods**: Supports GET and HEAD request caching
- **Query parameter support**: Handles URL query parameters in cache keys

## Installation

```bash
cd cachingProxy
go build -o caching-proxy
```

## Usage

### Start the proxy server
```bash
# Basic usage
./caching-proxy --origin http://dummyjson.com --port 3000

# Or with go run
go run . --origin http://dummyjson.com --port 3000
```

### Command-line flags
- `--port`: Port to run the proxy server on (default: 8080)
- `--origin`: Origin server URL to proxy requests to (required)
- `--clear-cache`: Clear the cache and exit

### Clear cache
```bash
./caching-proxy --clear-cache
```

## Examples

1. **Start the proxy:**
```bash
go run . --origin http://dummyjson.com --port 3000
```

2. **Make requests through the proxy:**
```bash
# First request - Cache MISS
curl http://localhost:3000/products
# Response includes: X-Cache: MISS

# Second request - Cache HIT
curl http://localhost:3000/products  
# Response includes: X-Cache: HIT

# HEAD requests
curl -I http://localhost:3000/users

# With query parameters
curl http://localhost:3000/products?limit=5
```

## How it works

1. Client sends request to proxy server
2. Proxy checks if response exists in cache
3. If **cache hit**: Returns cached response with `X-Cache: HIT`
4. If **cache miss**: 
   - Forwards request to origin server
   - Caches successful (200 OK) responses
   - Returns response with `X-Cache: MISS`

## Cache Behavior

- **Cached methods**: GET and HEAD requests
- **Cache key**: Combination of HTTP method, origin URL, and request path
- **Cached status codes**: Only 200 OK responses
- **Headers**: Most response headers are cached (excludes Date, Server, Set-Cookie, etc.)
- **Storage**: In-memory cache (data lost on restart)

## Project Structure

```
cachingProxy/
├── main.go          # Entry point and CLI handling
├── proxyServer.go   # HTTP proxy server implementation
├── cache.go         # Cache interface and in-memory implementation
├── go.mod           # Go module definition
└── README.md        # This file
```

## Development

To run in development mode:
```bash
go run . --origin http://httpbin.org --port 8080
```

Test with different endpoints:
```bash
curl http://localhost:8080/get
curl http://localhost:8080/json