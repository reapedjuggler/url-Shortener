package main

import (
	"fmt"
	"reapedjuggler/url-shortener/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Testing Golang Redis")

	// client := redis.NewClient(&redis.Options{
	// 	Addr:     "127.0.0.1:6379",
	// 	Password: "",
	// 	DB:       0,
	// })
	// pong, err := client.Ping().Result()

	// fmt.Println(pong, "pong", err)
	router := gin.Default()

	router.GET("/resolve", controllers.Resolve)
	router.POST("/shorten", controllers.Shorten)

	router.Run("localhost:3000")

}
