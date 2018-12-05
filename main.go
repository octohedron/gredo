package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
)

// POOL is a declared a global variable to store the Redis connection pool.
var POOL *redis.Pool

// LOCAL is an environment handler
var LOCAL = "host-local"

// Redis pool
func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

// Set up environment
func init() {
	LOCAL = os.Getenv("local")
	POOL = newPool("127.0.0.1:6379")
}

func main() {
	flag.Parse()
	reader := bufio.NewReader(os.Stdin)
	log.Printf("Redis command: %s %s %s", "SRANDMEMBER", flag.Args()[0], flag.Args()[1])
	log.Print("Is that OKAY? [Y/N] ")
	resp, _ := reader.ReadString('\n')
	switch resp {
	case "Y":
		conn := POOL.Get()
		defer conn.Close()
		data, _ := redis.Strings(conn.Do("SRANDMEMBER", flag.Args()[0], flag.Args()[1]))
		for _, u := range data {
			log.Println(u)
		}
	case "N":
		break
	default:
		log.Print("Only [Y/N] allowed")
		break
	}

}
