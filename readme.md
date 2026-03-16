Go In-Memory Cache

A simple high-performance in-memory cache written in Go.
This project is designed as a learning exercise to understand concurrency, data structures, and cache design patterns commonly used in backend systems.

The cache supports:
	•	Concurrent access
	•	Key-value storage
	•	TTL expiration
	•	Background cleanup
	•	Thread-safe operations

⸻

Features
	•	Thread-safe operations using sync.RWMutex
	•	Key-value storage
	•	TTL support for automatic expiration
	•	Background cleanup worker
	•	Simple and extensible architecture

Future enhancements may include:
	•	LRU eviction
	•	Cache sharding
	•	Metrics collection
	•	Persistence
	•	Distributed cache nodes

⸻

Project Structure

cache/
 ├── cache.go        # Cache implementation
 ├── item.go         # Cache item structure
 ├── cleanup.go      # Expiration worker
 └── shard.go        # Optional shard implementation

main.go              # Example usage


⸻

Architecture Overview

The cache stores items in a Go map and protects concurrent access using a read-write mutex.

Client
   │
   ▼
Cache API
   │
   ▼
Map Storage (key -> Item)
   │
   ▼
Background TTL Cleaner

Each item contains:
	•	stored value
	•	expiration timestamp

⸻

Cache Item Structure

type Item struct {
    Value      interface{}
    Expiration int64
}

Expiration stores a Unix timestamp used to determine whether an item has expired.

⸻

Cache Structure

type Cache struct {
    data  map[string]Item
    mutex sync.RWMutex
}

The cache uses:
	•	map for O(1) lookup
	•	RWMutex for safe concurrent reads and writes

⸻

Installation

Clone the repository:

git clone https://github.com/yourusername/go-cache.git
cd go-cache

Run the program:

go run main.go


⸻

Usage Example

cache := NewCache()

cache.Set("user", "Jason", 5*time.Minute)

value, found := cache.Get("user")
if found {
    fmt.Println(value)
}


⸻

API

Set

Set(key string, value interface{}, ttl time.Duration)

Stores a value with a time-to-live.

⸻

Get

Get(key string) (interface{}, bool)

Returns the cached value if it exists and has not expired.

⸻

Delete

Delete(key string)

Removes a key from the cache.

⸻

Clear

Clear()

Removes all entries from the cache.

⸻

TTL Expiration

Expired keys are cleaned by a background worker.

Example logic:

Every N seconds:
    scan cache
    delete expired keys

This prevents stale entries from accumulating.

⸻

Concurrency Model

The cache allows:
	•	multiple concurrent readers
	•	exclusive writers

This is implemented with sync.RWMutex:

Read  → RLock / RUnlock
Write → Lock / Unlock

This design provides good performance for read-heavy workloads.

⸻

Future Improvements

Possible improvements include:

LRU Eviction

Limit cache size and evict least recently used items.

Sharded Cache

Split cache into multiple shards to reduce lock contention.

Metrics

Track:
	•	cache hits
	•	cache misses
	•	evictions

Persistence

Save cache state to disk.

Distributed Cache

Run multiple cache nodes communicating over TCP.

⸻

Learning Goals

This project helps developers understand:
	•	Go concurrency primitives
	•	Mutex and RWMutex
	•	Goroutines
	•	Background workers
	•	Cache eviction strategies
	•	High-performance backend design

⸻

License

MIT License
:::


Future enhancement 14/03/2026

You already have a solid base cache (map + RWMutex + TTL). The features you listed are exactly the steps that evolve a simple cache into something closer to systems like Redis or Memcached.

I’ll explain what each feature solves, how it works conceptually, and give you implementation prompts so you can build it yourself rather than just copy code.

⸻

1️⃣ Cleanup Worker (TTL Expiration Engine)

Problem it solves

Right now expired items remain in memory until someone calls Get.

This causes:
	•	memory leak
	•	stale keys
	•	inefficient scans

Idea

Run a background goroutine that periodically removes expired keys.

Architecture

Cache
 ├── map[string]Item
 ├── RWMutex
 └── cleanup goroutine

How to build it

Add a worker that runs forever:

startCleanup()
    every N seconds
        scan cache
        delete expired keys

Prompt for implementation

1️⃣ Start worker in constructor

func NewCache() *Cache {
    c := &Cache{...}
    go c.startCleanup()
    return c
}

2️⃣ Use a ticker

ticker := time.NewTicker(1 * time.Minute)

3️⃣ Iterate keys

for key, item := range c.data
    if expired
        delete

Important learning
	•	goroutines
	•	background workers
	•	memory lifecycle

⸻

2️⃣ LRU Eviction (Limit Cache Size)

Problem it solves

Without limits the cache will grow forever.

Example:

1M keys
10M keys
100M keys

Eventually memory explodes.

Idea

When cache reaches capacity:

evict least recently used key

Data Structure

Classic design:

Hash Map + Doubly Linked List

Map gives O(1) lookup
List tracks usage order

Most recent ←→ ... ←→ Least recent

Example

Set(A)
Set(B)
Set(C)

Access(A)

Order becomes:
A → C → B

If capacity exceeded:

remove B

How to build it

Use Go’s container package:

container/list

Structure:

map[key] → pointer to list node
list node → key + value

Prompt

Implement:

type entry struct {
    key string
    value interface{}
}

Then maintain:

map[string]*list.Element

Operations:

Get → move element to front
Set → insert front
Evict → remove back


⸻

3️⃣ Sharding (Concurrency Scaling)

Problem

Right now you have one mutex.

If 1000 goroutines access cache:

everyone waits on same lock

This becomes a huge bottleneck.

Idea

Split cache into multiple shards.

Cache
 ├── shard1
 ├── shard2
 ├── shard3
 └── shard4

Each shard has its own:

map
mutex

Key routing

Choose shard using hash:

shardIndex = hash(key) % shardCount

Example

key = "user123"

hash → 938423
938423 % 8 = shard 7

Architecture

Client
   │
   ▼
ShardRouter
   │
   ├── shard 1 (map + lock)
   ├── shard 2
   ├── shard 3
   └── shard N

Prompt

Create:

type Shard struct {
    data map[string]Item
    mutex sync.RWMutex
}

type Cache struct {
    shards []*Shard
}

Routing:

func (c *Cache) getShard(key string) *Shard


⸻

4️⃣ Metrics (Observability)

Problem

You don’t know if your cache is working well.

Important metrics:

cache hits
cache misses
evictions
entries

Example:

hit rate = hits / (hits + misses)

How to build

Add counters:

type Metrics struct {
    Hits uint64
    Misses uint64
    Evictions uint64
}

Use atomic operations:

atomic.AddUint64()

Prompt

Update metrics inside:

Get()
Set()
Delete()
Evict()

Example:

if found
    Hits++
else
    Misses++


⸻

5️⃣ Persistence (Cache Restart Survival)

Problem

When your process restarts:

cache is lost

Sometimes that’s fine. Sometimes not.

Idea

Periodically write cache state to disk.

Two approaches:

Snapshot

Every N minutes:

serialize map
write file

Append log

Every write:

append SET key value

Similar to Redis AOF.

Prompt

Use Go encoding:

encoding/gob
encoding/json

Example flow:

saveSnapshot()
    serialize map
    write file

On startup:

loadSnapshot()
    read file
    rebuild map


⸻

6️⃣ Distributed Cache Nodes

Problem

One machine can only store limited memory.

To scale:

multiple cache nodes

Example cluster:

Node A
Node B
Node C
Node D

Key distribution

Use consistent hashing.

Concept:

key → hash ring → node

Example:

key1 → node B
key2 → node A
key3 → node C

Communication

Nodes talk via:

TCP
HTTP
gRPC

Example commands:

SET key value
GET key
DEL key

Architecture

Client
   │
   ▼
Hash Router
   │
   ├── Cache Node 1
   ├── Cache Node 2
   └── Cache Node 3

Prompt

Build:

TCP server
command parser
cache execution


⸻

🚀 Recommended Development Order

Don’t build everything at once.

Order:

Phase 1

TTL cleanup worker

Phase 2

LRU eviction

Phase 3

Sharded cache

Phase 4

Metrics

Phase 5

Persistence

Phase 6

Distributed cache


⸻

🧠 What You’ll Learn From This Project

This single project teaches:
	•	concurrency
	•	memory management
	•	lock contention
	•	distributed systems
	•	eviction algorithms
	•	system observability

These concepts appear in real systems like:
	•	Redis
	•	Memcached
	•	Hazelcast

⸻

If you’d like, I can also show you something extremely valuable for interviews and system design:

How high-performance caches avoid scanning the entire map for expiration.

This trick dramatically improves scalability.