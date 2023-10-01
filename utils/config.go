package utils

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var redisClient *redis.Client = nil
var mongoClient *mongo.Client = nil

const CollectionName = "CollectionName"
const MONGODB_URI = "MONGODB_URI"
const DatabaseName = "DatabaseName"

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}

func InitMongodb() {
	client, err := mongo.NewClient(options.Client().ApplyURI(GetKeyFromEnv("MONGODB_URI")))
	log.Print(client, " . ", " client ")
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	mongoClient = client
	// defer client.Disconnect(ctx)docker-compose logs redis

}

func GetClient() *redis.Client {
	if redisClient == nil {
		InitRedis()
		log.Print("Redis client initialized successfully")
	}
	return redisClient
}

func GetKeyFromEnv(key string) string {
	return os.Getenv(key)
}

// Basically this function will bind the variables into os from env file
func GoDotEnvVariable() {
	// log.Print("Loading the .env file")
	err := godotenv.Load(".env")
	log.Print("loaded the .env file")
	if err != nil {
		log.Print("not loaded the .env file")
		log.Fatalf("Error loading .env file ", err)
	}
	log.Print(os.Getenv("MONGODB_URI"))
}

func GetMongoClient() *mongo.Client {
	if mongoClient == nil {
		InitMongodb()
	}
	return mongoClient
}
