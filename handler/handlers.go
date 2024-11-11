package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shienlee73/url-shortener/shortener"
	"github.com/shienlee73/url-shortener/store"
)

func (server *Server) Index(c *gin.Context) {
	c.Redirect(http.StatusFound, "/url-shortener/")
}

type CreateShortUrlRequest struct {
	OriginalUrl string `json:"originalUrl"`
	UserId      string `json:"userId"`
}

type CustomizeShortUrlRequest struct {
	OriginalUrl    string `json:"originalUrl"`
	UserId         string `json:"userId"`
	CustomShortUrl string `json:"customShortUrl"`
}

func (server *Server) CreateShortUrl(c *gin.Context) {
	var createShortUrlRequest CreateShortUrlRequest
	if err := c.ShouldBindJSON(&createShortUrlRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	shortUrl, err := shortener.GenerateShortUrl(createShortUrlRequest.OriginalUrl, createShortUrlRequest.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// save to redis
	err = server.store.SaveUrlMapping(shortUrl, createShortUrlRequest.OriginalUrl, createShortUrlRequest.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// save to bolt
	err = server.store.CreateURLMapping(store.URLMapping{
		ShortUrl:    shortUrl,
		OriginalUrl: createShortUrlRequest.OriginalUrl,
		UserId:      createShortUrlRequest.UserId,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"shortUrl": shortUrl,
	})
}

func (server *Server) HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")

	// HIT: retrieve from redis
	originalUrl, err := server.store.RetrieveOriginalUrl(shortUrl)
	if err == nil {
		fmt.Printf("HIT: %s\n", shortUrl)
		c.Redirect(http.StatusFound, originalUrl)
		return
	}

	// NOT HIT: retrieve from bolt and add it to redis
	fmt.Printf("NOT HIT: %s\n", shortUrl)
	urlMapping, err := server.store.RetrieveURLMapping(shortUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = server.store.SaveUrlMapping(shortUrl, urlMapping.OriginalUrl, urlMapping.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, urlMapping.OriginalUrl)
}

func (server *Server) CustomizeShortUrl(c *gin.Context) {
	var customizeShortUrlRequest CustomizeShortUrlRequest
	if err := c.ShouldBindJSON(&customizeShortUrlRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	shortUrl := customizeShortUrlRequest.CustomShortUrl

	// save to redis
	err := server.store.SaveUrlMapping(shortUrl, customizeShortUrlRequest.OriginalUrl, customizeShortUrlRequest.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// save to bolt
	err = server.store.CreateURLMapping(store.URLMapping{
		ShortUrl:    shortUrl,
		OriginalUrl: customizeShortUrlRequest.OriginalUrl,
		UserId:      customizeShortUrlRequest.UserId,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"shortUrl": shortUrl,
	})
}
