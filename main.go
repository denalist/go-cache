package main

import (
	"fmt"
	"time"

	"go-cache/cache"
)

func main() {

	c := cache.NewCache()

	fmt.Println("Setting key...")
	c.Set("user", "Jason", 5*time.Second)

	value, found := c.Get("user")

	if found {
		fmt.Println("Value:", value)
	} else {
		fmt.Println("Key not found")
	}

	fmt.Println("Waiting 6 seconds for expiration...")
	time.Sleep(6 * time.Second)

	value, found = c.Get("user")

	fmt.Println("Checking cache after expiration...")

	if found {
		fmt.Println("Value:", value)
	} else {
		fmt.Println("Key expired")
	}

}


