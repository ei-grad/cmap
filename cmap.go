package cmap

import "hash/fnv"
import "sync"

// CMap is concurrent sharded map[string]string with Get/Set methods
type CMap struct {
	shards []*Shard
}

func New(nShards int) CMap {
	shards := make([]*Shard, nShards)
	for i := 0; i < nShards; i++ {
		shards[i] = NewShard()
	}
	return CMap{shards: shards}
}

func (c CMap) hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32() % uint32(len(c.shards))
}

func (c CMap) Get(key string) string {
	return c.shards[c.hash(key)].Get(key)
}

func (c CMap) Set(key, value string) {
	c.shards[c.hash(key)].Set(key, value)
}

// Shard is concurrent map[string]string with Get/Set methods
type Shard struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewShard() *Shard {
	return &Shard{
		data: make(map[string]string),
	}
}

func (c *Shard) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}

func (c *Shard) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}
