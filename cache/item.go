package cache

import "time"

// LRU for eviction to maintain cache capacity, most frequent used -> .. -> least frequent used
// Hash map (for search ) + Doubly linked list
type entry struct {
	Key        string
	Value      interface{}
	Expiration int64
}

func (e entry) IsExpired() bool {
	if e.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > e.Expiration
}
