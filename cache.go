package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type CacheEntry struct {
	Data        []byte            `json:"data"`
	Headers     map[string]string `json:"headers"`
	StatusCode  int               `json:"status_code"`
	Timestamp   time.Time         `json:"timestamp"`
	ContentType string            `json:"content_type"`
}

type MemoryCache struct {
	data map[string]*CacheEntry
	mu   sync.RWMutex
}

type Cache interface {
	Get(key string) (*CacheEntry, bool)
	Set(key string, entry *CacheEntry)
	Clear() error
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: make(map[string]*CacheEntry),
	}
}

func (c *MemoryCache) Get(key string) (*CacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, found := c.data[key]
	return entry, found
}

func (c *MemoryCache) Set(key string, entry *CacheEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = entry
	c.persistToDisk(key, entry)
}

func (c *MemoryCache) persistToDisk(key string, entry *CacheEntry) {
	cacheDir := getCacheDir()
	os.MkdirAll(cacheDir, 0755)
	filename := filepath.Join(cacheDir, hashKey(key)+".json")
	data, _ := json.Marshal(entry)
	// fmt.Printf("Persisting cache to disk: %s\n", filename)
	// fmt.Println("Data:", string(data))
	os.WriteFile(filename, data, 0644)
}

func (c *MemoryCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]*CacheEntry)
	return clearCacheData()
}

func hashKey(key string) string {
	HASH := md5.Sum([]byte(key))
	return hex.EncodeToString(HASH[:])
}

func getCacheDir() string {
	return filepath.Join(os.TempDir(), "caching-proxy")
}
func clearCacheData() error {
	cacheDir := getCacheDir()
	return os.RemoveAll(cacheDir)
}
