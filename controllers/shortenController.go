package controllers

import (
	"io"
	"reapedjuggler/url-shortener/utils"

	"github.com/gin-gonic/gin"
)

func Shorten (ctx *gin.Context) {

	// recieves a url from 

	data, err := io.ReadAll(ctx.Request.Body)

	if err != nil {
		shortUrl, errFromUtil := utils.ConvertToBase64(string(data))
		
		if errFromUtil != nil {

			// save in redis

		} else {
			panic(errFromUtil)
		}

	} else {
		panic(err)
	}

}
