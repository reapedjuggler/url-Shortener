package controllers

import (
	"fmt"
	"net/http"
	"reapedjuggler/url-shortener/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func Resolve(ctx *gin.Context) { // resolving a url
	// ctx.JSON(http.StatusAccepted, "fuck world")

	// code, _ := ctx.Params.Get("shorturl")

	code := ctx.Request.URL.Query().Get("shorturl")

	// write the logic to find the corresponding original
	// url and reirect it

	var client *redis.Client = utils.GetClient()

	// fmt.Println("code\n", code, " -> \n sad", code1)

	val, err := client.Get(code).Result()

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusNotFound, "The given short url is invalid")
	}

	fmt.Println(val, "\nval")

	loginUrl := val // "httsp://www.google.com/"

	// rdirect it here
	ctx.Redirect(http.StatusOK, loginUrl)

}
