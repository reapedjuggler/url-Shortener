package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	urlParser "net/url"
	"reapedjuggler/url-shortener/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var redisctx = context.Background()

type url struct {
	Urls string `form:"urls"`
}

type ErrorMessage struct {
	Message    string
	StatusCode int
}

func (i *url) marshalbinary() ([]byte, error) {
	return json.Marshal(i)
}
func Shorten(ctx *gin.Context) {
	// recieves a url from
	urls := &url{}
	if err := ctx.ShouldBind(urls); err != nil {
		ctx.String(http.StatusBadRequest, "bad request: %v", err)
		return
	}
	_, err := urlParser.ParseRequestURI(urls.Urls)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorMessage{"Invalid URL", 400})
		return
	}

	log.Print(ctx.ContentType(), " Content-Type")
	log.Print(urls, " Inside the shorten controller")

	var client *redis.Client = utils.GetClient()
	nextid, err := client.Get("nextid").Result()
	if err == redis.Nil {
		client.Set("nextId", 1, 0)
		nextid = "1"
	}

	nextidint, err := strconv.ParseInt(nextid, 10, 64)
	base64encoded, err := utils.ConvertToBase64(nextidint)
	log.Print(base64encoded)
	_, err = client.Get(urls.Urls).Result()

	if err != redis.Nil {
		fmt.Println(err, " err")
		ctx.JSON(http.StatusBadRequest, ErrorMessage{"URL already exists", 400})
		panic("url already exists in the database")
	}

	shorturl := base64encoded
	shorturl = utils.CompleteShortUrl(shorturl)
	status := client.Set(shorturl, urls.Urls, 3600*1e9)
	log.Print(status)
	client.Set("nextid", nextidint+1, 0)
	shorturl = "http://localhost:3000/resolve?shorturl=" + shorturl
	ctx.JSON(http.StatusAccepted, "Here is your shoterened URL: "+shorturl)
}
