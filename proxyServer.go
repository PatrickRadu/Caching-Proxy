package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
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
	cacheKey := p.createCacheKey(r)
	if entry, found := p.cache.Get(cacheKey); found {
		fmt.Printf("Cache HIT for %s\n", cacheKey)
		p.serveCachedResponse(w, entry)
		return
	}

	fmt.Printf("Cache MISS for %s\n", cacheKey)
	p.forwardRequest(w, r, cacheKey)
}

func (p *ProxyServer) serveCachedResponse(w http.ResponseWriter, entry *CacheEntry) {
	for key, value := range entry.Headers {
		w.Header().Set(key, value)
	}
	w.Header().Set("X-Cache", "HIT")
	w.WriteHeader(entry.StatusCode)
	w.Write(entry.Data)
}

func (p *ProxyServer) forwardRequest(w http.ResponseWriter, r *http.Request, cacheKey string) {
	targetURL := fmt.Sprintf("%s%s", p.origin, r.URL.Path)
	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	originUrl, _ := url.Parse(p.origin)
	proxyReq.Host = originUrl.Host

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Failed to reach origin server", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("X-Cache", "MISS")

	if (r.Method == "GET" || r.Method == "HEAD") && resp.StatusCode == http.StatusOK {
		p.cacheResponse(cacheKey, resp, body)
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func (p *ProxyServer) cacheResponse(key string, resp *http.Response, body []byte) {
	headers := make(map[string]string)

	// Copy headers we want to cache
	for k, v := range resp.Header {
		if len(v) > 0 && p.shouldCacheHeader(k) {
			headers[k] = strings.Join(v, ", ")
		}
	}

	// Create cache entry
	entry := &CacheEntry{
		Data:        body,
		Headers:     headers,
		StatusCode:  resp.StatusCode,
		Timestamp:   time.Now(),
		ContentType: resp.Header.Get("Content-Type"),
	}

	// Store in cache
	p.cache.Set(key, entry)
}

func (p *ProxyServer) createCacheKey(r *http.Request) string {
	// Include method, origin, and path in key
	return fmt.Sprintf("%s:%s%s", r.Method, p.origin, r.URL.Path)
}

func (p *ProxyServer) shouldCacheHeader(header string) bool {
	// Don't cache these headers
	skipHeaders := map[string]bool{
		"date":          true,
		"server":        true,
		"x-cache":       true,
		"set-cookie":    true,
		"authorization": true,
	}

	return !skipHeaders[strings.ToLower(header)]
}
