package cmap

import "hash/fnv"
import "sync"

// Map is sharded concurrent map[string]string with Get/Set methods
type Map struct {
	shards []*Shard
}

func New(nShards int) Map {
	shards := make([]*Shard, nShards)
	for i := 0; i < nShards; i++ {
		shards[i] = NewShard()
	}
	return Map{shards: shards}
}

func (c Map) hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32() % uint32(len(c.shards))
}

func (c Map) Get(key string) string {
	return c.shards[c.hash(key)].Get(key)
}

func (c Map) Set(key, value string) {
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
