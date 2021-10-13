package utils

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

var client *redis.Client

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func InitialiseRedis() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6380"
	}
	client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Error connecting redis %v", err.Error())
	}

	log.Println("Connected to redis!")
}
