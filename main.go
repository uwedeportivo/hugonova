package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)

	router := gin.New()

	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasSuffix(path, "index.md") {
			path = strings.TrimSuffix(path, "index.md")
		} else if strings.HasSuffix(path, ".md") {
			path = strings.TrimSuffix(path, ".md")
		}
		backendUrl := "http://localhost:1313" + path

		response, err := http.Get(backendUrl)
		if err != nil || response.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		reader := response.Body
		defer func(reader io.ReadCloser) {
			err := reader.Close()
			if err != nil {
				log.Println("failed to close response body")
			}
		}(reader)
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")

		extraHeaders := map[string]string{}

		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
