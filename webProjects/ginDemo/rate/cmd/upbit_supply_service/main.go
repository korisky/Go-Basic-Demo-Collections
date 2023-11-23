package main

import "github.com/gin-gonic/gin"

func main() {

	router := gin.Default()

	router.GET("/api/functionx/info", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "info",
		})
	})

	err := router.Run(":20950")
	if err != nil {
		return
	}
}
