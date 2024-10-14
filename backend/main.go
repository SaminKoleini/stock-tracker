package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		testParam := c.Query("test")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"test":    testParam,
		})
	})

	fmt.Println("Everything is working fine")

	r.Run("localhost:8000")
}
