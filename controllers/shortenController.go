package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reapedjuggler/url-shortener/utils"
	"strconv"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var redisctx = context.Background()

type url struct {
	Urls string `json:"urls"`
}

func (i *url) marshalbinary() ([]byte, error) {
	return json.Marshal(i)
}
func Shorten(ctx *gin.Context) {
	// recieves a url from
	urls := &url{}
	if err := ctx.ShouldBindJSON(&urls); err != nil {
		fmt.Print(err)
		return
	}
	log.Print(urls)
	var client *redis.Client = utils.GetClient()
	nextid, err := client.Get("nextid").Result()
	if err == redis.Nil {
		fmt.Println("inside redis empty")
		client.Set("nextid", 1, 0)
		nextid = "1"
	}
	nextidint, err := strconv.ParseInt(nextid, 10, 64)
	log.Print(nextid, " -> ", nextidint)
	base64encoded, err := utils.ConvertToBase64(nextidint)
	log.Print(base64encoded)
	val, err := client.Get(urls.Urls).Result()
	fmt.Println(val, " ", "val0")
	if err != redis.Nil {
		fmt.Println(err, " err")
		ctx.JSON(http.StatusBadRequest, urls)
		panic("url already exists in the database")
	}
	shorturl := base64encoded
	shorturl = utils.CompleteShortUrl(shorturl)
	status := client.Set(shorturl, urls.Urls, 0)
	log.Print(status)
	client.Set("nextid", nextidint+1, 0)
	ctx.JSON(http.StatusAccepted, urls)
}
