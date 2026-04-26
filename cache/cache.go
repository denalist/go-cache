// The cache stores items in a Go map and protects concurrent access using a read-write mutex.
package cache

import (
	"container/list"
	"log"
	"sync"
	"time"
)

type Cache struct {
	data     map[string]*list.Element
	lruList  *list.List
	capacity int
	mutex    sync.RWMutex

	// Metrics
	hits      int64
	misses    int64
	evictions int64
}

// Constructor- Go doesnt have build in constructor like python def __init__()
// the function returbs a pointer to a cache construct instead of copying one
// &Cache {...} - create a cache construct and return its pointer

// maps must be initialised before use , so we create a empty map: map[string]item.Item, string
// is the datatype of key and item.Item is the value type
func NewCache(capacity int) *Cache {
	c := &Cache{
		data:     make(map[string]*list.Element),
		lruList:  list.New(),
		capacity: capacity,
	}

	go c.startCleanup()
	return c
}

// set values: lock, compute expire and store, interface{} -> "Any"
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}

	newEntry := &entry{
		Key:        key,
		Value:      value,
		Expiration: expiration,
	}

	element, found := c.data[key]

	if found {
		// if found, update the entry val, update EXP and move to the front
		existing := element.Value.(*entry) // type assertion, converts interface{} → *entry
		existing.Value = value             // struct field assignment use "="
		existing.Expiration = expiration
		c.lruList.MoveToFront(element)

	} else {
		element := c.lruList.PushFront(newEntry)
		c.data[key] = element // the value expects List.element
	}

	// check capacity
	if len(c.data) > c.capacity {
		last := c.lruList.Back()
		if last != nil {
			c.lruList.Remove(last)

			// list.Element.Value stores *entry
			evicted := last.Value.(*entry) // first convert the value back to *entry.
			delete(c.data, evicted.Key)

			c.evictions++
		}
	}

}

// get values
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.Lock() // dont use RLock() as we are deleting expired items
	defer c.mutex.Unlock()

	res, found := c.data[key]
	if !found {
		c.misses++
		return nil, false
	}

	// convert the res back to entry
	rec := res.Value.(*entry)

	if rec.IsExpired() {
		c.lruList.Remove(res)
		delete(c.data, key)
		c.misses++
		return nil, false
	}

	c.lruList.MoveToFront(res)
	c.hits++
	return rec.Value, true
}

// delete values
func (c *Cache) Delete(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	element, found := c.data[key]

	if !found {
		return false
	}

	c.lruList.Remove(element)
	delete(c.data, key)
	return true
}

// clear values
func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]*list.Element)
	c.lruList = list.New()
}

// cleanup worker/ TTL expiration engine runs in background
func (c *Cache) startCleanup() {
	ticker := time.NewTicker(1 * time.Minute)

	for range ticker.C {
		c.mutex.Lock()

		for key, element := range c.data {
			entry := element.Value.(*entry)
			if entry.IsExpired() {
				c.lruList.Remove(element)
				delete(c.data, key)
			}

		}

		c.mutex.Unlock() // once exited, unlock to avoid deadlock issue
		log.Println("cleanup worker job completed")
	}
}

// TODO
// Multi sharfds to reduce lock contention, dsitributed cache nodes over TCP
// tests
