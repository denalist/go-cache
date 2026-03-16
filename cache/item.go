package cache

import "time"

// import "fmt"

type Item struct {
	Value      interface{}
	Expiration int64
}

func (item Item) IsExpired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}
