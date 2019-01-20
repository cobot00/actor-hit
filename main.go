package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	rand.Seed(time.Now().UnixNano())

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
	value := getCookie(c, "test")
	log.Printf("cookie value: %v", value)

	if value == "" {
		setCookie(c, "test", strconv.Itoa(rand.Intn(100)))
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Hello, world",
	})
}

func setCookie(c *gin.Context, key string, value string) {
	log.Println("setCookie")
	c.SetCookie(key, value, 3600, "/", "", false, true)
}

func getCookie(c *gin.Context, key string) string {
	value, err := c.Cookie(key)
	if err != nil {
		return ""
	}
	return value
}
