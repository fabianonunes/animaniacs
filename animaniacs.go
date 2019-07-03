package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func randomChar() string {
	chars := []string{
		"DOT",
		"WAKKO",
		"YAKKO",
	}
	rand.Seed(time.Now().UnixNano())
	return chars[rand.Intn(len(chars))]
}

func main() {
	char := os.Getenv("CHAR")
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Success")
	})

	r.GET("/fail", func(c *gin.Context) {
		c.String(http.StatusInternalServerError, "Error")
	})

	r.GET("/v1/:name", func(c *gin.Context) {
		sleep := c.Query("sleep")
		if len(sleep) > 0 {
			duration, err := strconv.ParseInt(sleep, 10, 64)
			if err == nil {
				time.Sleep(time.Duration(duration) * time.Millisecond)
			}
		}

		if len(char) > 0 {
			c.File("./art/" + char)
		} else {
			c.File("./art/" + randomChar())
		}
	})
	_ = r.Run("0.0.0.0:3000")
}
