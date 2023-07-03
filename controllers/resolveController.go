package controllers

import (
	"fmt"
	"net/http"
	"reapedjuggler/url-shortener/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type ResolvedRequest struct {
	ActualUrl string `json:"actualUrl"`
}

func Resolve(ctx *gin.Context) {
	code := ctx.Request.URL.Query().Get("shorturl")
	var client *redis.Client = utils.GetClient()

	val, err := client.Get(code).Result()

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusNotFound, "The given short url is invalid")
	}

	fmt.Println(val, "\nval")
	ctx.JSON(200, ResolvedRequest{ActualUrl: val})
}
