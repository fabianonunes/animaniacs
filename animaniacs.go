package main

import (
	"bufio"
	"fmt"
	"github.com/DataDog/zstd"
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
	"sync"
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

var netClient = &http.Client{
	Timeout: time.Second * 600,
}

var f = fmt.Sprintf

const FilePath = "/tmp/dat.zec"

func compress(input io.Reader, output io.Writer, params *ZstdParams) *exec.Cmd {
	cmd := exec.Command(
		"zstd",
		"-qc",
		f("-T%d", params.Threads),
		f("-%d", params.Level),
	)
	cmd.Stdin = input
	cmd.Stdout = output
	return cmd
}

func main() {
	char := os.Getenv("CHAR")
	r := gin.Default()
	instanceName := randomdata.SillyName()
	store := memstore.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("sessionid", store))

	file, _ := os.Create(FilePath)

	defaultParams := &ZstdParams{
		Level:   3,
		Threads: 0,
	}

	var isFileValid bool
	var mutex = &sync.Mutex{}

	//contentLength := get.Header.Get("Content-Length")
	//context.Header("sf.com.zec-length", contentLength)

	r.GET("/zec-pipe", func(context *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()

		stat, _ := os.Stat(FilePath)
		now := time.Now()
		expiresAt := stat.ModTime().Add(1 * time.Minute)

		if isFileValid && now.Before(expiresAt) {
			context.File(FilePath)
			return
		}

		_ = os.Truncate("/tmp/dat.zec", 0)
		_ = context.ShouldBindQuery(&defaultParams)

		get, _ := netClient.Get("http://localhost:8080/apps/iterator")
		responseWriter := context.Writer
		fileWriter := bufio.NewWriter(file)

		// (getZec.Body) zecPipeW»|«zecPipeR (compress) zstdPipeW»|«zstdPipeR (tee >fileWriter) > responseWriter

		zecPipeR, zecPipeW := io.Pipe()
		zstdPipeR, zstdPipeW := io.Pipe()

		command := compress(zecPipeR, zstdPipeW, defaultParams)
		tee := io.TeeReader(zstdPipeR, fileWriter)

		defer zecPipeW.Close()
		defer zstdPipeW.Close()

		go func() {
			defer zecPipeW.Close()
			_, _ = io.Copy(zecPipeW, get.Body)
		}()

		go func() {
			defer zstdPipeW.Close()
			if _, err := io.Copy(responseWriter, tee); err != nil {
				_ = os.Truncate(FilePath, 0)
			} else {
				_ = fileWriter.Flush()
				isFileValid = true
			}
		}()

		_ = command.Run()
	})

	r.GET("/zec-compress", func(context *gin.Context) {
		writer := context.Writer
		fileName := os.Getenv("ZEC_PATH")

		params := &ZstdParams{
			Level:   2,
			Threads: 1,
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
				// TODO: testar se um defer pipeWriter.Close() nesse closure tem o mesmo efeito
				_ = pipeWriter.Close()
			}
		}()

		_ = cmd.Run()
	})

	r.GET("/zec-compress-go", func(context *gin.Context) {
		params := &ZstdParams{
			Level: 2,
		}
		_ = context.ShouldBindQuery(&params)

		writer := context.Writer
		filePath := os.Getenv("ZEC_PATH")

		file, _ := os.Open(filePath)

		zstdWriter := zstd.NewWriterLevel(writer, params.Level)
		defer zstdWriter.Close()

		if _, err := io.Copy(zstdWriter, file); err != nil {
			zstdWriter.Close()
		}
	})

	r.GET("/zec", func(context *gin.Context) {
		filePath := os.Getenv("ZEC_PATH")
		context.File(filePath)
	})

	r.GET("/", func(c *gin.Context) {
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
