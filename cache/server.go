package cache

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	cache    *SharedCache
	listener net.Listener
}

func NewServer(addr string, cache *SharedCache) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{cache: cache, listener: ln}, nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn) // interface reading data

	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 0 {
			continue
		}

		switch strings.ToUpper(parts[0]) {
		case "GET":
			if len(parts) < 2 {
				fmt.Fprintf(conn, "ERR usage: GET <key>\n")
				continue
			}
			val, ok := s.cache.Get(parts[1])
			if ok {
				fmt.Fprintf(conn, "VALUE %v\n", val)
			} else {
				fmt.Fprintf(conn, "NIL\n")
			}

		case "SET":
			if len(parts) < 4 {
				fmt.Fprintf(conn, "ERR usage: SET <key> <ttl_ms> <value>\n")
				continue
			}
			ttlMs, err := strconv.ParseInt(parts[2], 10, 64)
			if err != nil {
				fmt.Fprintf(conn, "ERR invalid ttl\n")
				continue
			}
			s.cache.Set(parts[1], parts[3], time.Duration(ttlMs)*time.Millisecond)
			fmt.Fprintf(conn, "OK\n")

		case "DEL":
			if len(parts) < 2 {
				fmt.Fprintf(conn, "ERR usage: DEL <key>\n")
				continue
			}
			s.cache.Delete(parts[1])
			fmt.Fprintf(conn, "OK\n")

		default:
			fmt.Fprintf(conn, "ERR unknown command\n")
		}
	}
}

func (s *Server) Serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("server stopped: %v", err)
			return
		}
		go s.handleConn(conn) // each client gets own goroutine
	}
}

func (s *Server) Close() error {
	return s.listener.Close()
}
