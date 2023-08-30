package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reapedjuggler/url-shortener/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceUrl struct {
	Urls    string
	LongUrl string
}

type ErrorMessage struct {
	Message    string
	StatusCode int
}

func ShortenService(ctx *gin.Context, urls *ServiceUrl) string {
	var redisClient *redis.Client = utils.GetClient()
	// log.Print("Getting mongodb connection")
	var mongoClient *mongo.Client = utils.GetMongoClient()

	nextid, err := redisClient.Get("nextid").Result()
	if err == redis.Nil {
		log.Print(err, " err in service")
		nextid = "1"
	}

	nextidint, _ := strconv.ParseInt(nextid, 10, 64)
	base64encoded := utils.ConvertToBase64(nextidint)
	log.Print(base64encoded)
	_, err = redisClient.Get(base64encoded).Result()

	if err != redis.Nil {
		fmt.Println(err, " err")
		ctx.JSON(http.StatusBadRequest, ErrorMessage{"URL already exists", 400})
		panic(err)
	}

	shorturl := base64encoded
	shorturl = utils.CompleteShortUrl(shorturl)

	// Caching starts here
	// Enter into cache first and then into mongodb first so that consistency is there in cache and mongodb
	InsertIntoRedis(redisClient, shorturl, nextidint, *urls)
	// P.S - Study about cache policies

	// log.Print("reached the db layer ", mongoClient)
	_, err = InsertIntoMongodb(mongoClient, shorturl, urls)
	if err != nil {
		// Internal server error;
		removeFromRedis(redisClient, shorturl)
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{"Internal server error", 500})
	}
	shorturl = "http://localhost:3000/resolve?shorturl=" + shorturl
	return shorturl
}

func InsertIntoRedis(redisClient *redis.Client, shorturl string, nextIdInt int64, urls ServiceUrl) *redis.StatusCmd {
	status := redisClient.Set(shorturl, urls.Urls, 3600*1e9)
	log.Print(status)
	redisClient.Set("nextid", nextIdInt+1, 0)
	return status
}

func InsertIntoRedisWithoutNextId(redisClient *redis.Client, shorturl string, urls ServiceUrl) *redis.StatusCmd {
	status := redisClient.Set(shorturl, urls.Urls, 3600*1e9)
	log.Print(status)
	return status
}

func InsertIntoMongodb(mongoClient *mongo.Client, shorturl string, urls *ServiceUrl) (bool, error) {
	db := mongoClient.Database(utils.GetKeyFromEnv(utils.DatabaseName))
	log.Print(db, " db")
	coll := db.Collection(utils.GetKeyFromEnv(utils.CollectionName))
	// log.Print(coll, " coll")
	doc := ServiceUrl{Urls: shorturl, LongUrl: urls.Urls}
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return false, err
	} else {
		log.Printf("Inserted document with _id: %v\n", result)
		return true, nil
	}
}

// Function to remove a key from redis
func removeFromRedis(redisClient *redis.Client, key string) (int64, error) {
	status := redisClient.Del(key)
	valFromDeleteKey, err := status.Result()
	if err != redis.Nil {
		log.Print(err, " Error while deleting the key from redis")
	}
	log.Print(valFromDeleteKey, " valFromDeleteKey")
	return valFromDeleteKey, err
}
