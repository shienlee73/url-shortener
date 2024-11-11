package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shienlee73/url-shortener/shortener"
	"github.com/shienlee73/url-shortener/store"
	"github.com/shienlee73/url-shortener/token"
)

func (server *Server) Index(c *gin.Context) {
	c.Redirect(http.StatusFound, "/url-shortener/")
}

type CreateShortUrlRequest struct {
	OriginalUrl string `json:"originalUrl"`
}

type CustomizeShortUrlRequest struct {
	OriginalUrl    string `json:"originalUrl"`
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

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	shortUrl, err := shortener.GenerateShortUrl(createShortUrlRequest.OriginalUrl, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// save to redis
	err = server.store.SaveUrlMapping(shortUrl, createShortUrlRequest.OriginalUrl, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// save to bolt
	err = server.store.CreateURLMapping(store.URLMapping{
		ID:          uuid.NewString(),
		ShortUrl:    shortUrl,
		OriginalUrl: createShortUrlRequest.OriginalUrl,
		UserId:      authPayload.UserID,
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

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	// save to redis
	err := server.store.SaveUrlMapping(shortUrl, customizeShortUrlRequest.OriginalUrl, authPayload.UserID)
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
		UserId:      authPayload.UserID,
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

func (server *Server) GetURLMappings(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	urlMappings, err := server.store.RetrieveURLMappings(authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"urlMappings": urlMappings,
	})
}
