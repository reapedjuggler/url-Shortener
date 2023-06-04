package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reapedjuggler/url-shortener/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var redisCtx = context.Background()

type Url struct {
	Urls string `json:"urls"`
}

func (i *Url) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}
func Shorten(ctx *gin.Context) {
	// recieves a url from
	urls := &Url{}
	if err := ctx.ShouldBindJSON(&urls); err != nil {
		fmt.Print(err)
		return
	}

	var client *redis.Client = utils.GetClient()
	nextId, err := client.Get("nextId").Result()
	if err == redis.Nil {
		fmt.Println("Inside redis empty error")
		client.Set("nextId", 1, 0)
		nextId = "1"
	}
	nextIdInt, err := strconv.ParseInt(nextId, 10, 64)

	fmt.Println(nextId, " ->sad ", nextIdInt)
	base64Encoded, err := utils.ConvertToBase64(nextIdInt)

	val, err := client.Get(base64Encoded).Result()
	fmt.Println(val, " ", "val0")
	if err != nil {
		panic(err)
	}
	status := client.Set(base64Encoded, urls, 0).Err()
	client.Set("nextId", nextIdInt+1, 0)
	fmt.Print("status \n", status)

	ctx.JSON(http.StatusAccepted, urls)
}
