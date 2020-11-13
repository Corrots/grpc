package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	TestRedisSentinel()
}

func TestRedisSentinel() {
	cluster := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"192.168.56.26:7000",
			"192.168.56.26:7001",
			"192.168.56.26:7002",
			"192.168.56.26:7003",
			"192.168.56.26:7004",
		},
	})
	ctx := context.Background()
	if _, err := cluster.Ping(ctx).Result(); err != nil {
		log.Fatalf("redis cluster new err: %v\n", err)
	}
	//
	rand.Seed(time.Now().UnixNano())
	file, err := os.Create("redis-cluster-test.log")
	if err != nil {
		log.Fatal(err)
	}
	for {
		n := strconv.Itoa(rand.Intn(100000))
		key := "k" + n
		_, err := cluster.Set(ctx, key, "v"+n, 0).Result()
		if err != nil {
			file.WriteString(fmt.Sprintf("SET err: %v\n", err))
			continue
		}
		result, err := cluster.Get(ctx, key).Result()
		if err != nil {
			//fmt.Printf("GET err: %v\n", err)
			file.WriteString(fmt.Sprintf("GET err: %v\n", err))
			continue
		}
		fmt.Printf("%s = %s\n", key, result)
		time.Sleep(time.Second)
	}
	defer file.Close()
	defer cluster.Close()
}
