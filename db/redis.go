package db

import "github.com/go-redis/redis/v8"

func ConnectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		// Addr: "docker.for.mac.localhost:6379",
		Addr: "localhost:6379",
		// Addr: "redis:6379",

		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
