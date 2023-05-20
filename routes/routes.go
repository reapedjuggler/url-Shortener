package routes

import "github.com/gin-gonic/gin"



func InitServer() {

	router := gin.Default()

	

	router.Run("localhost:8000")
}
