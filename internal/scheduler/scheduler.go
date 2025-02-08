package scheduler

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func saveItem(item string, id string) {
	fmt.Println("Loading environment variables...")
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env")
	}

	redAddr := os.Getenv("RED_ADDR")

	client := redis.NewClient(&redis.Options{
		Addr:     redAddr,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	ctx := context.Background()

	err = client.Set(ctx, id, item, 0).Err()

	if err != nil {
		panic(err)
	}
}

func getValue(id string) string {
	fmt.Println("Loading environment variables...")
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env")
	}

	redAddr := os.Getenv("RED_ADDR")

	client := redis.NewClient(&redis.Options{
		Addr:     redAddr,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	ctx := context.Background()

	val, err := client.Get(ctx, "foo").Result()

	if err != nil {
		panic(err)
	}

	return val
}
