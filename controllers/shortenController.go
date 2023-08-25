package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	urlParser "net/url"
	"reapedjuggler/url-shortener/services"

	"github.com/gin-gonic/gin"
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

	// Service call
	serviceUrl := &services.ServiceUrl{Urls: urls.Urls}
	shorturl := services.ShortenService(ctx, serviceUrl)

	shorturl = "http://localhost:3000/resolve?shorturl=" + shorturl
	ctx.JSON(http.StatusAccepted, "Here is your shoterened URL: "+shorturl)
}
