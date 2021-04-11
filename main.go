package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/", ping)
	v1 := r.Group("/v1")
	{
		v1.GET("/", ping)
		v1.GET("/ping", ping)
		v1.GET("/generator", generator)
		v1.GET("/validator", validator)
	}

	r.Run(":3000")
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
