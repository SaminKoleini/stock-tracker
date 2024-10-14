package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"os"
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

	r.GET("/stocks", func(c *gin.Context) {
		err := godotenv.Load()
		if err != nil {
			return
		}
		symbol := c.Query("symbol")
		api_key := os.Getenv("API_KEY")
		api_url := os.Getenv("API_URL")
		interval := "1min"
		outputsize := "1"
		if api_key == "" || api_url == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "API_KEY or API_URL not set"})
			return
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("%s?symbol=%s&interval=%s&outputsize=%s", api_url, symbol, interval, outputsize), nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		req.Header.Set("Authorization", fmt.Sprintf("apikey %s", api_key))

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{"error": "Failed to fetch data from API"})
			return
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
			return
		}

		values, ok := result["values"].([]interface{})
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid data format for values"})
			return
		}

		if len(values) > 0 {
			firstValue, ok := values[0].(map[string]interface{})
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid data format for values item"})
				return
			}

			high, ok := firstValue["high"].(string)
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid data format for high"})
				return
			}

			low, ok := firstValue["low"].(string)
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid data format for low"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"symbol": symbol, "high": high, "low": low})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No data found in values"})
		}

	})

	fmt.Println("Everything is working fine")

	r.Run() // Listen on localhost:8080
}
