package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	urlParser "net/url"
	"reapedjuggler/url-shortener/services"
)

type url struct {
	Urls string `form:"urls"`
}

type ErrorMessage struct {
	Message    string
	StatusCode int
}

func ShortenController(ctx *gin.Context) {
	// receives a url from
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

	// Service call
	serviceUrl := &services.ServiceUrl{Urls: urls.Urls, LongUrl: urls.Urls}
	shorturl := services.ShortenService(ctx, serviceUrl)
	ctx.JSON(http.StatusAccepted, "Here is your shoterened URL: "+shorturl)
}
