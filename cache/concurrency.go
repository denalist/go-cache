package cache

import (
	"hash/fnv"
	"time"
)

type SharedCache struct {
	shards []*Cache
	count  int
}

// key -> hash(key) % N -> shard[i].mutex
func NewSharedCache(shards, capacityPerShard int) *SharedCache {
	s := &SharedCache{
		shards: make([]*Cache, shards),
		count:  shards,
	}
	for i := range s.shards {
		s.shards[i] = NewCache(capacityPerShard)
	}
	return s
}

func (s *SharedCache) shard(key string) *Cache {
	h := fnv.New32a()
	h.Write([]byte(key))
	return s.shards[h.Sum32()%uint32(s.count)]
}

func (s *SharedCache) Set(key string, value interface{}, ttl time.Duration) {
	s.shard(key).Set(key, value, ttl)
}

func (s *SharedCache) Get(key string) (interface{}, bool) {
	return s.shard(key).Get(key)
}

func (s *SharedCache) Delete(key string) bool {
	return s.shard(key).Delete(key)
}
