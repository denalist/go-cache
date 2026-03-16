// The cache stores items in a Go map and protects concurrent access using a read-write mutex.
package cache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	data     map[string]*list.Element
	lruList  *list.List
	capacity int
	mutex    sync.RWMutex
}

// LRU for eviction to maintain cache capacity, most frequent used -> .. -> least frequent used
// Hash map (for search ) + Doubly linked list
type entry struct {
	Key        string
	Value      interface{}
	Expiration int64
}

// Constructor- Go doesnt have build in constructor like python def __init__()
// the function returbs a pointer to a cache construct instead of copying one
// &Cache {...} - create a cache construct and return its pointer

// maps must be initialised before use , so we create a empty map: map[string]item.Item, string
// is the datatype of key and item.Item is the value type
func NewCache() *Cache {
	c := &Cache{
		data:     make(map[string]*list.Element),
		lruList:  list.New(),
		capacity: 100,
	}

	go c.startCleanup()
	return c
}

// set values: lock, compute expire and store, interface{} -> "Any"
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	expiration := time.Now().Add(ttl).UnixNano()
	entry := &entry{
		Key:        key,
		Value:      value,
		Expiration: expiration,
	}

	element, found := c.data[key]

	if found {
		c.lruList.MoveToFront(element)
		c.data[key] = element
	} else {
		element := c.lruList.PushFront(entry)
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
		}
	}

}

// get values
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.Lock() // dont use RLock() as we are deleting expired items
	defer c.mutex.Unlock()

	res, found := c.data[key]
	if !found {
		return nil, false
	}

	if res.IsExpired() {
		delete(c.data, key)
		return nil, false
	}

	return res.Value, true
}

// delete values
func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

// clear values
func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]Item)
}

// cleanup worker/ TTL expiration engine runs in background
func (c *Cache) startCleanup() {
	ticker := time.NewTicker(1 * time.Minute)

	for range ticker.C {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		for key, item := range c.data {
			if item.IsExpired() {
				delete(c.data, key)
			}

			fmt.Println("Clean up worker job completed.")
		}
	}
}
