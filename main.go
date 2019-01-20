package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	r := createRouter()

	setRoute(r)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	r.Run(":" + port)
}

func createRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/css", "assets/css")
	r.Static("/image", "assets/image")
	return r
}

func setRoute(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		index(c)
	})
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Hello, world",
	})
}
