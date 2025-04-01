package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	EmployeeDetailsRedisClient *redis.Client
)

type RedisConfig struct {
	Address  string
	Password string
	Timeout  time.Duration
}

func InitRedisConnections() {

	config := RedisConfig{
		Address:  os.Getenv("redis_host"),
		Password: os.Getenv("redis_password"),
		Timeout:  5 * time.Second, // Connection timeout
	}

	EmployeeDetailsRedisClient = NewRedisClient(config, 0)
	log.Println("Employee Details Redis connection established")

}

func NewRedisClient(config RedisConfig, dbNumber int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       dbNumber,
	})

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis DB %d: %v", dbNumber, err)
	}

	return rdb
}

// Close all Redis connections
func CloseRedisConnections() {
	if err := EmployeeDetailsRedisClient.Close(); err != nil {
		log.Printf("Error closing Employee Redis connection: %v", err)
	}
}
