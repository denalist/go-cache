package cache

// Metrics data struct
type Metrics struct {
	Hits      int64
	Misses    int64
	Evictions int64
	Size      int
}

// function to return the current metrics stats
func (c *Cache) GetMetrics() Metrics {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return Metrics{
		Hits:      c.hits,
		Misses:    c.misses,
		Evictions: c.evictions,
		Size:      len(c.data),
	}
}
