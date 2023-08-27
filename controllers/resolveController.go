package controllers

import (
	"context"
	"log"
	"net/http"
	"reapedjuggler/url-shortener/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResolvedRequest struct {
	ActualUrl string `json:"actualUrl"`
}

type ResultFromMongoDB struct {
	Urls    string `json:"urls" bson:"urls"`
	Longurl string `json:"longurl" bson:"longurl"`
	// _id     ObjectId `json:"bson_id" bson:"_id,omitempty"`
	Id primitive.ObjectID `bson:"_id" json:"id"`
}

func Resolve(ctx *gin.Context) {
	code := ctx.Request.URL.Query().Get("shorturl")
	var client *redis.Client = utils.GetClient()
	log.Print(code, " code")
	val, err := client.Get(code).Result()

	// Cache lookup
	if err == nil {
		log.Print(val, " Corresponding Resolved URL")
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
		panic(err)
	}

	log.Print(correspondingUrl, " correspondingUrl")
	ctx.Redirect(http.StatusMovedPermanently, correspondingUrl.Longurl)


	// Add in the cache as well, I think this should be done by a goroutine
	
}
