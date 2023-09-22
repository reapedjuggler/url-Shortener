package main

import (
	"fmt"
	"reapedjuggler/url-shortener/controllers"
	"reapedjuggler/url-shortener/middlewares"
	"reapedjuggler/url-shortener/utils"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	// Jai Shree Ram, Jai Hanuman ji
	fmt.Println("Testing Golang Redis")
	router := gin.Default()

	router.Use(middlewares.RateLimitMiddleware())
	router.Use()
	router.Use(static.Serve("/", static.LocalFile("./views", true)))
	utils.GoDotEnvVariable()
	router.GET("/resolve", controllers.Resolve)
	router.POST("/shorten", controllers.ShortenController)

	router.Run("localhost:3000")
}
