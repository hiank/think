package main

import (
	"flag"
	"fmt"

	"github.com/hiank/think/db"
)

func main() {

	flag.Parse()

	status := db.RedisMaster().Set("testInt", 1)
	if status.Err() != nil {
		panic(status.Err())
	}

	val, err := db.RedisMaster().Get("testInt").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("result : ", val)


	val, err = db.RedisSlave().Get("testInt").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("result from slave: ", val)
}
