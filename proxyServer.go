package main

import (
	"fmt"
	"net/http"
)

type ProxyServer struct {
	origin string
	cache  Cache
}

func NewProxyServer(origin string) *ProxyServer {
	return &ProxyServer{
		origin: origin,
		cache:  NewMemoryCache(),
	}
}

func (p *ProxyServer) Start(port int) error {
	http.HandleFunc("/", p.handleRequest)
	fmt.Printf("Server listening on port %d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (p *ProxyServer) handleRequest(w http.ResponseWriter, r *http.Request) {

}

func (p *ProxyServer) createCacheKey(r *http.Request) string {
	// Include method, origin, and path in key
	return fmt.Sprintf("%s:%s%s", r.Method, p.origin, r.URL.Path)
}
