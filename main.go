package main

import (
	"fmt"
	"reapedjuggler/url-shortener/controllers"
	"reapedjuggler/url-shortener/utils"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Testing Golang Redis")
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./views", true)))
	utils.GoDotEnvVariable()

	router.GET("/resolve", controllers.Resolve)
	router.POST("/shorten", controllers.Shorten)

	router.Run("localhost:3000")
}
