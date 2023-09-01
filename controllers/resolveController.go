package controllers

import (
	"context"
	"log"
	"net/http"
	"reapedjuggler/url-shortener/services"
	"reapedjuggler/url-shortener/utils"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResolvedRequest struct {
	ActualUrl string `json:"actualUrl"`
}

type ResultFromMongoDB struct {
	Urls    string             `json:"urls" bson:"urls"`
	Longurl string             `json:"longurl" bson:"longurl"`
	Id      primitive.ObjectID `bson:"_id" json:"id"`
}

func Resolve(ctx *gin.Context) {
	code := ctx.Request.URL.Query().Get("shorturl")
	var client *redis.Client = utils.GetClient()
	var wg *sync.WaitGroup = &sync.WaitGroup{}
	log.Print(code, " code")
	val, err := client.Get(code).Result()

	// Cache lookup
	if err == nil {
		log.Print("Found in cache")
		log.Print(val, " Corresponding Resolved URL")
		ctx.Redirect(http.StatusMovedPermanently, val)
		// ctx.JSON(http.StatusAccepted, "Resolved")
		log.Print("Redirected")
		return
	}

	log.Print("Cache miss occurred for shorturl = ", code)

	// Cache miss, lookup in mongodb
	mongoClient := utils.GetMongoClient()
	db := mongoClient.Database(utils.GetKeyFromEnv(utils.DatabaseName))
	coll := db.Collection(utils.GetKeyFromEnv(utils.CollectionName))
	filter := bson.D{{Key: "urls", Value: code}}
	var correspondingUrl ResultFromMongoDB = ResultFromMongoDB{}

	errFromFindOne := coll.FindOne(context.TODO(), filter).Decode(&correspondingUrl)
	if errFromFindOne != nil {
		ctx.JSON(http.StatusNotFound, "The given short url is invalid")
		panic(errFromFindOne)
	}

	// Add in the cache as well, I think this should be done by a goroutine
	// And yes it will be done by it, cause we don't care even if it fails, as we already persisted it
	log.Print(correspondingUrl, " correspondingUrl")

	// Read about this, aisa to nahi ho raha ki before inserting into redis y program exit kar ja raha hai
	// Answer: It won't because we are already listening on a server and hence the main file is never existing.
	// Even though just add a wait group just for learning.
	wg.Add(1)
	go services.InsertIntoRedisWithoutNextId(client, code, services.ServiceUrl{Urls: code, LongUrl: correspondingUrl.Longurl}, wg)

	// ctx.JSON(http.StatusAccepted, "Resolved")
	ctx.Redirect(http.StatusMovedPermanently, correspondingUrl.Longurl)

	// ToDo:
	// Rate Limit it
	// Rate Limit it
	// Rate Limit it
	// Rate Limit it
	// Rate Limit it
	// Rate Limit it
	// Rate Limit it
	// Rate Limit it
	// Rate Limit it
	// Rate Limit it
	// Rate Limit it
	// Rate Limit it
}
