package main

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
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

type ZstdParams struct {
	Level   int `form:"l"`
	Threads int `form:"T"`
}

var f = fmt.Sprintf

func main() {
	char := os.Getenv("CHAR")
	r := gin.Default()
	instanceName := randomdata.SillyName()
	store := memstore.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("sessionid", store))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Success")
	})

	r.GET("/zec-compress", func(context *gin.Context) {
		writer := context.Writer
		fileName := os.Getenv("ZEC_PATH")

		params := &ZstdParams{
			Level:   1,
			Threads: 2,
		}
		_ = context.ShouldBindQuery(&params)

		pipeReader, pipeWriter := io.Pipe()
		defer pipeWriter.Close()

		cmd := exec.Command(
			"zstd",
			"-qc",
			f("-T%d", params.Threads),
			f("-%d", params.Level),
			fileName,
		)
		cmd.Stdout = pipeWriter

		go func() {
			if _, err := io.Copy(writer, pipeReader); err != nil {
				_ = pipeWriter.Close()
			}
		}()

		_ = cmd.Run()
	})

	r.GET("/zec", func(context *gin.Context) {
		filePath := os.Getenv("ZEC_PATH")
		context.File(filePath)
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
