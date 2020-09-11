package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/hiank/think/db"
)

func main() {

	flag.Parse()

	ctx := context.Background()

	status := db.TryRedisMaster().Set(ctx, "testInt", 1, 0)
	if status.Err() != nil {
		panic(status.Err())
	}

	val, err := db.TryRedisMaster().Get(ctx, "testInt").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("result : ", val)

	val, err = db.TryRedisSlave().Get(ctx, "testInt").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("result from slave: ", val)
}
