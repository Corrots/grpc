package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	redisSentinelTest()
}

func redisSentinelTest() {
	sentinels := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName: "mymaster",
		SentinelAddrs: []string{
			"192.168.56.26:26379",
			"192.168.56.26:26380",
			"192.168.56.26:26381",
		},
	})
	ctx := context.Background()
	if _, err := sentinels.Ping(ctx).Result(); err != nil {
		log.Fatalf("sentinel client new err: %v\n", err)
	}
	//
	rand.Seed(time.Now().UnixNano())
	counter := 0
	for {
		counter++
		n := strconv.Itoa(rand.Intn(100000))
		key := "k" + n
		_, err := sentinels.Set(ctx, key, "v"+n, 0).Result()
		if err != nil {
			if counter%100 == 0 {
				fmt.Printf("SET err: %v\n", err)
			}
		}
		result, err := sentinels.Get(ctx, key).Result()
		if counter%100 == 0 {
			if err != nil {
				fmt.Printf("GET err: %v\n", err)
			}
			fmt.Printf("%s = %s\n", key, result)
		}
		time.Sleep(time.Millisecond * 100)
	}
}
