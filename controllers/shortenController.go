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
	urls string `json:"urls"`
}

func (i *url) marshalbinary() ([]byte, error) {
	return json.Marshal(i)
}
func shorten(ctx *gin.Context) {
	// recieves a url from
	urls := &url{}
	if err := ctx.shouldbindjson(&urls); err != nil {
		fmt.print(err)
		return
	}
	log.print(urls)
	var client *redis.client = utils.getclient()
	nextid, err := client.get("nextid").result()
	if err == redis.nil {
		fmt.println("inside redis empty")
		client.set("nextid", 1, 0)
		nextid = "1"
	}
	nextidint, err := strconv.parseint(nextid, 10, 64)
	log.print(nextid, " -> ", nextidint)
	base64encoded, err := utils.converttobase64(nextidint)
	log.print(base64encoded)
	val, err := client.get(urls.urls).result()
	fmt.println(val, " ", "val0")
	if err != redis.nil {
		fmt.println(err, " err")
		ctx.json(http.statusbadrequest, urls)
		panic("url already exists in the database")
	}
	shorturl := base64encoded
	shorturl = utils.completeshorturl(shorturl)
	status := client.set(urls.urls, shorturl, 0)
	log.print(status)
	client.set("nextid", nextidint+1, 0)
	ctx.json(http.statusaccepted, urls)
}

