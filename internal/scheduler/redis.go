package scheduler

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func auth() *redis.Client {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	// Construct the full path to .env
	envPath := cwd + "/.env"

	if len(cwd) == 0 {
		cwd = "root/lari-go"
	}

	fmt.Println("Loading environment variables...")
	err = godotenv.Load(envPath)

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

	return client
}

func saveList(id string, items []string) {
	client := auth()
	ctx := context.Background()

	for _, item := range items {
		err := client.RPush(ctx, id, item).Err()
		if err != nil {
			log.Fatalf("Failed to push item to Redis list. %v", err)
		}
	}

	fmt.Println("Item saved to Redis list.")

}

func Validate(patientId string, timeslotId string) bool {
	client := auth()
	ctx := context.Background()

	index, err := client.LPos(ctx, timeslotId, patientId, redis.LPosArgs{}).Result()

	if err == redis.Nil {
		fmt.Println("List does not exist")
		return false
	} else if err != nil {
		fmt.Println("Error")
		return false
	} else {
		fmt.Printf("Element found at index %v", index)
		return true
	}
}

func Remove(timeslotId string) {
	client := auth()
	ctx := context.Background()

	deletedCount, err := client.Del(ctx, timeslotId).Result()

	if err != nil {
		log.Fatalf("Error deleting key %q: %v", timeslotId, err)
	}

	fmt.Printf("Deleted %d key(s).\n", deletedCount)
}

func saveItem(id string, item string) {
	client := auth()

	ctx := context.Background()

	err := client.Set(ctx, id, item, 0).Err()

	if err != nil {
		panic(err)
	}
}

func getValue(id string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	// Construct the full path to .env
	envPath := cwd + "/.env"

	if len(cwd) == 0 {
		cwd = "root/lari-go"
	}

	fmt.Println("Loading environment variables...")
	err = godotenv.Load(envPath)

	if err != nil {
		log.Fatal("Error loading .env")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RED_ADDR"),
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	ctx := context.Background()

	val, err := client.Get(ctx, id).Result()

	if err != nil {
		panic(err)
	}

	return val
}
