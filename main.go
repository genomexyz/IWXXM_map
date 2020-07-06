package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.Static("/static", "./static")

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	
	//show map
//	r.LoadHTMLGlob("template/*")
	r.GET("/index", func(c *gin.Context) {
		r.LoadHTMLGlob("template/*")
		c.HTML(http.StatusOK, "index.html", nil)
	})
	
	//data
	r.GET("/data", func(c *gin.Context) {
		r.LoadHTMLGlob("static/data/*")
		c.HTML(http.StatusOK, "recent_weather.txt", nil)
	})

	return r
}

func main() {
	r := setupRouter()
	
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
