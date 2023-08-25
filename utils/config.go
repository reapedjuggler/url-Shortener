package utils

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *redis.Client = nil
var mongoClient *mongo.Client = nil

func InitRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func InitMongodb() {
	client, err := mongo.NewClient(options.Client().ApplyURI(GetKeyFromEnv("MONGO_URI")))
	log.Print(client, " . ", " client ")
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
}

func GetClient() *redis.Client {
	if client == nil {
		InitRedis()
	}
	return client
}

func GetKeyFromEnv(key string) string {
	return os.Getenv(key)
}

// Basically this function will bind the variables into os from env file
func GoDotEnvVariable() {
	// load .env file
	log.Print("Loading the .env file")
	err := godotenv.Load(filepath.Join(path_dir, ".env"))
	if err != nil {
		log.Fatalf("Error loading .env file ", err)
	}
	log.Print()
}

func GetMongoClient() *mongo.Client {
	if mongoClient == nil {
		InitMongodb()
	}
	return mongoClient
}
