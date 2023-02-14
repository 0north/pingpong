package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	PING_URL = os.Getenv("PING_URL")
)

func ping() {
	for {
		resp, err := http.Get(PING_URL)
		if err != nil {
			fmt.Printf("cannot call ping: %s\n", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Printf("response from %v: %s", PING_URL, string(body))

		time.Sleep(5 * time.Second)
	}
}

func pong() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func main() {
	fmt.Println("PING_URL: " + PING_URL)
	if PING_URL != "" {
		ping()
	} else {
		pong()
	}
}
