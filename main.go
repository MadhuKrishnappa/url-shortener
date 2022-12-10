package main

import (
	"fmt"

	"url-shortener/handler"

	"url-shortener/db"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.POST("/create-short-url", func(c *gin.Context) {
		handler.CreateShortUrl(c)
	})

	r.GET("/:shortUrl", func(c *gin.Context) {
		handler.HandleShortUrlRedirect(c)
	})

	db.Connect()

	err := r.Run(":9808")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}
}
