package main

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
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
	instanceName := randomdata.SillyName()
	store := memstore.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("sessionid", store))

	r.GET("/healthcheck", func(c *gin.Context) {
		c.String(http.StatusOK, "Success")
	})

	r.GET("/count", func(c *gin.Context) {
		session := sessions.Default(c)
		value := session.Get("counter")

		var count int
		if value == nil {
			count = 0
		} else {
			count = value.(int)
			count++
		}
		session.Set("counter", count)
		_ = session.Save()
		c.JSON(200, gin.H{instanceName: count})
	})

	r.GET("/fail", func(c *gin.Context) {
		c.String(http.StatusInternalServerError, "Error")
	})

	r.GET("/", func(c *gin.Context) {
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
