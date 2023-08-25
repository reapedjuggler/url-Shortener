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
	Urls string
}

type ErrorMessage struct {
	Message    string
	StatusCode int
}

func ShortenService(ctx *gin.Context, urls *ServiceUrl) string {
	var client *redis.Client = utils.GetClient()
	log.Print("Getting mongodb connection")
	var mongoClient *mongo.Client = utils.GetMongoClient()

	nextid, err := client.Get("nextid").Result()
	if err == redis.Nil {
		client.Set("nextId", 1, 0)
		nextid = "1"
	}

	nextidint, err := strconv.ParseInt(nextid, 10, 64)
	base64encoded, err := utils.ConvertToBase64(nextidint)
	log.Print(base64encoded)
	_, err = client.Get(base64encoded).Result()

	if err != redis.Nil {
		fmt.Println(err, " err")
		ctx.JSON(http.StatusBadRequest, ErrorMessage{"URL already exists", 400})
		panic("url already exists in the database")
	}

	shorturl := base64encoded
	shorturl = utils.CompleteShortUrl(shorturl)

	// Caching starts here
	// status := client.Set(shorturl, urls.Urls, 3600*1e9)
	// log.Print(status)
	// client.Set("nextid", nextidint+1, 0)

	// Enter into cache first and then into mongodb first so that consistency is there in cache and mongodb
	// P.S - Study about cache policies
	log.Print("reached the db layer")
	coll := mongoClient.Database("shorturls").Collection("shortUrls")
	doc := ServiceUrl{Urls: shorturl}
	result, err := coll.InsertOne(context.TODO(), doc)
	log.Printf("Inserted document with _id: %v\n", result)
	return shorturl
}
