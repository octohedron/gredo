package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

// POOL is a global variable to store the Redis connection pool.
var POOL *redis.Pool

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

func logErr(err error, msg string) {
	if err != nil {
		log.Printf("ERROR: %v %s", err, msg)
	}
}

// Set up redis connection pool
func init() {
	POOL = newPool("127.0.0.1:6379")
}

// dump is a func to dump random members of a redis set to a file
func dump(setID string, amount string) {
	log.Printf("Redis command: %s %s %s", "SRANDMEMBER", setID, amount)
	log.Print("Is that OKAY? [y/n] ")
	var resp string
	fmt.Scanln(&resp)
	log.Printf("You chose: %s", resp)
	switch resp {
	case "y":
		conn := POOL.Get()
		defer conn.Close()
		data, _ := redis.Strings(conn.Do("SRANDMEMBER", setID, amount))
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
}

// load is a func to load items from a file to a redis set
func load(setID string, amount string) {
	amounToAdd, err := strconv.Atoi(amount)
	logPanic(err, "COULDN'T CONVERT STRING TO INT")
	conn := POOL.Get()
	defer conn.Close()
	file, err := os.Open("./" + setID + ".txt")
	logErr(err, "ERROR READING FILE")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	duplicates := 0
	amountAdded := 0
	for scanner.Scan() {
		s := scanner.Text()
		ok, err := redis.Bool(conn.Do("SADD", setID, s))
		logErr(err, "ERROR ADDING TO REDIS SET")
		if ok {
			amountAdded++
		} else {
			duplicates++
		}
		// only load amount specified
		if amounToAdd == amountAdded {
			break
		}
	}
	log.Printf("Loaded %d items to set %s DUPLICATES: %d", amountAdded, setID, duplicates)
	log.Printf("Operation completed.")
}

func main() {
	flag.Parse()
	argAmnt := len(flag.Args())
	// 3 arguments = {dump|load} {SET_IDENTIFIER} {AMOUNT}
	if argAmnt == 3 {
		switch flag.Args()[0] {
		case "dump":
			dump(flag.Args()[1], flag.Args()[2])
		case "load":
			load(flag.Args()[1], flag.Args()[2])
		}

	} else {
		log.Print("Available operations:")
		log.Print("./gredo {dump|load} {SET_IDENTIFIER} {AMOUNT}")
	}
}
