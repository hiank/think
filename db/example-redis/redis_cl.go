package main

import (
	"flag"
)

func main() {

	flag.Parse()

	// ctx := context.Background()

	// ar := db.NewAutoRedis(ctx, &db.RedisConf{
	// 	// MasterURL: k8s.TryServiceURL(ctx, k8s.TypeKubIn, "redis-master", ""),
	// 	// SlaveURL:  k8s.TryServiceURL(ctx, k8s.TypeKubIn, "redis-slave", ""),
	// 	Password: os.Getenv("REDIS_PASSWORD"),
	// 	DB:       0,
	// })
	// status := ar.TryMaster().Set(ctx, "testInt", 1, 0)
	// if status.Err() != nil {
	// 	panic(status.Err())
	// }

	// val, err := ar.TryMaster().Get(ctx, "testInt").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("result : ", val)

	// val, err = ar.TrySlave().Get(ctx, "testInt").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("result from slave: ", val)
}
