package controllers

import (
	"context"
	"encoding/json"
	"fmt"
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
		// fmt.Println(err, " err")
		// ctx.JSON(http.StatusNotFound, "The given short url is invalid")
		log.Print(val, " Corresponding Resolved URL")
		fmt.Println(val, "\nval")
		ctx.Redirect(http.StatusMovedPermanently, val)
		return
	}

	// Cache miss, lookup in mongodb
	mongoClient := utils.GetMongoClient()
	db := mongoClient.Database(utils.GetKeyFromEnv(utils.DatabaseName))
	coll := db.Collection(utils.GetKeyFromEnv(utils.CollectionName))
	filter := bson.D{{Key: "urls", Value: code}}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	log.Print(code, " . ", cursor)
	results := []ResultFromMongoDB{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for _, result := range results {
		res, _ := json.Marshal(result)
		// log.Print(string(res), " res")
		fmt.Println(string(res))
	}
}
