# go-cache

An in-memory LRU cache with TTL eviction, sharding for concurrent throughput, and a TCP server for distributed access.

## Design

### Core Cache (`cache/cache.go`)
Each cache shard is a hash map + doubly linked list implementing LRU eviction:
- **Map** provides O(1) key lookup
- **Linked list** tracks recency — most recently used at front, least at back
- On capacity overflow, the tail (LRU item) is evicted
- A background goroutine sweeps expired TTL entries every minute
- A single `sync.RWMutex` protects each shard

### Sharding (`cache/concurrency.go`)
`SharedCache` splits the keyspace across N independent `Cache` shards using FNV hash:

```
hash(key) % N  →  shard[i]
```

Each shard has its own lock, so operations on different keys run concurrently. With 16 shards, lock contention drops ~16x compared to a single global mutex.

### TCP Server (`cache/server.go`)
A lightweight line-protocol TCP server wraps `SharedCache`. Each client connection gets its own goroutine. Protocol:

```
SET <key> <ttl_ms> <value>   →  OK
GET <key>                    →  VALUE <val>  |  NIL
DEL <key>                    →  OK
```

Test with `nc`:
```bash
nc localhost 8080
SET foo 5000 bar
GET foo
DEL foo
```

### HTTP API (`api.go`)
Gin REST endpoints over the single-shard `Cache`:

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` / `PUT` | `/cache/:key` | Set a value |
| `GET` | `/cache/:key` | Get a value |
| `DELETE` | `/cache/:key` | Delete a value |

### Metrics (`cache/metrics.go`)
Each shard tracks hits, misses, and evictions accessible via `GetMetrics()`.

## TCP vs HTTP

HTTP is built on top of TCP. TCP is the transport layer — it just moves bytes reliably between two machines. HTTP is a protocol that defines a specific *format* for those bytes.

```
TCP  (raw bytes, you define the format)
 └── HTTP (structured format built on TCP)
      └── curl (speaks HTTP only)
```

When you run `curl localhost:8080`, curl connects via TCP but then sends an HTTP-formatted request:
```
GET / HTTP/1.1
Host: localhost:8080
...
```
The TCP server reads raw bytes and sees `"GET / HTTP/1.1\n..."` — not `"GET foo\n"` — so nothing matches the switch cases. Use `nc` instead, which sends raw bytes exactly as typed.

| | TCP (this server) | HTTP (Gin) |
|---|---|---|
| Format | you define it | fixed standard |
| Client | `nc`, custom client | `curl`, browsers |
| Overhead | minimal | headers, status codes, MIME types |
| Use case | internal service-to-service | public APIs |

Redis, Memcached, and the PostgreSQL wire protocol all use raw TCP for the same reason: lower overhead and full control over the format when you own both sides of the connection.

## Running

```bash
go run .
```

Server starts on `:8080`. Shutdown with `Ctrl+C`.
