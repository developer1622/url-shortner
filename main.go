package main

import (
	"hash/fnv"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LongURL struct {
	Url string `json:"longURL"`
}

func main() {
	// add it in the etc/hosts file
	host := "http://localhost:8087"
	storage := make(map[uint64]string)
	r := gin.Default()

	// create short URL
	r.POST("/v1/url", func(c *gin.Context) {
		var reqBody LongURL
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, "invalid request")
			return
		}

		// FNV-64 for optimization
		fnv64 := fnv.New64()
		fnv64.Write([]byte(reqBody.Url))
		urlID := fnv64.Sum64()

		// check if already exists or not?
		_, isIDexist := storage[urlID]
		if isIDexist {
			log.Println("short URL already existing for this URL")
			c.JSON(http.StatusOK, host+"/v1/url/"+strconv.FormatUint(urlID, 10))
			return
		}

		storage[urlID] = reqBody.Url
		c.JSON(http.StatusOK, host+"/v1/url/"+strconv.FormatUint(urlID, 10))
	})

	// access actual URL with this short URL
	r.GET("/v1/url/:id", func(c *gin.Context) {
		longURLID := c.Param("id")
		if longURLID == "" {
			c.JSON(http.StatusBadRequest, "invalid request")
			return
		}

		urlID, err := strconv.ParseUint(longURLID, 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "invalid request")
			return
		}

		longURL, isIDexist := storage[urlID]
		if !isIDexist {
			c.JSON(http.StatusBadRequest, "invalid request")
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, longURL)
	})

	// just added comment, which will be reverted
	// To fetch all the URLs
	r.GET("/v1/urls", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, storage)
	})
	r.Run(":8087")
}
