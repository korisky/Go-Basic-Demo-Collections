package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	r := gin.New()

	r.GET("/", helloWorldGinHandler)

	r.Run()
}

func helloWorldGinHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello From Gin",
	})
}
