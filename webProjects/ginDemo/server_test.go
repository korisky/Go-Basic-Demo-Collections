package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

func Test_start(t *testing.T) {
	r := gin.New()
	r.GET("/", helloWorldGinHandler)
	r.Run()
}

func helloWorldGinHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello From Gin",
	})
}
