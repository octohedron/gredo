package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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

func logPanic(err error, msg string) {
	if err != nil {
		log.Printf("ERROR: %v %s", err, msg)
		panic(err)
	}
}

// Set up environment
func init() {
	LOCAL = os.Getenv("local")
	POOL = newPool("127.0.0.1:6379")
}

func main() {
	flag.Parse()
	argAmnt := len(flag.Args())
	if argAmnt == 2 {
		log.Printf("Redis command: %s %s %s", "SRANDMEMBER", flag.Args()[0], flag.Args()[1])
		log.Print("Is that OKAY? [y/n] ")
		var resp string
		fmt.Scanln(&resp)
		log.Printf("You chose: %s", resp)
		switch resp {
		case "y":
			conn := POOL.Get()
			defer conn.Close()
			data, _ := redis.Strings(conn.Do("SRANDMEMBER", flag.Args()[0], flag.Args()[1]))
			// string for appending all the values
			res := ""
			// print and build string for saving
			for _, u := range data {
				res += u + "\n"
			}
			// write the whole body at once
			err := ioutil.WriteFile(flag.Args()[0]+".txt", []byte(res), 0777)
			logPanic(err, "Writing file")
			log.Printf("Operation completed. Exported %d items", len(data))
		case "n":
			break
		default:
			log.Print("Only [y/n] allowed")
			break
		}
	} else {
		log.Print("Available operations:")
		log.Print("./gredo {SET_IDENTIFIER} {AMOUNT}")
	}
}
