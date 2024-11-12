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

func (server *Server) LoginPage(c *gin.Context) {
	c.Redirect(http.StatusFound, "/url-shortener/login.html")
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
	shortUrlId := uuid.NewString()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// save to redis
	err = server.store.SaveUrlMappingToRedis(shortUrl, createShortUrlRequest.OriginalUrl, shortUrlId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// save to bolt
	err = server.store.CreateURLMapping(store.URLMapping{
		ID:          shortUrlId,
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

	var urlMapping store.URLMapping
	var err error

	// HIT: retrieve from redis
	urlMappingCache, err := server.store.RetrieveUrlMappingFromRedis(shortUrl)
	if err == nil {
		fmt.Println(urlMappingCache)
		fmt.Println("HIT:", shortUrl)
		urlMapping.ID = urlMappingCache["urlMappingId"]
		urlMapping.OriginalUrl = urlMappingCache["originalUrl"]
	} else {
		// NOT HIT: retrieve from bolt and add it to redis
		fmt.Println("NOT HIT:", shortUrl)
		urlMapping, err = server.store.RetrieveURLMapping(shortUrl)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		err = server.store.SaveUrlMappingToRedis(shortUrl, urlMapping.OriginalUrl, urlMapping.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	// create click stat
	err = server.store.CreateClickStat(store.ClickStat{
		ID:           uuid.NewString(),
		UrlMappingId: urlMapping.ID,
		ClickTime:    time.Now(),
		IpAddress:    c.ClientIP(),
		Referer:      c.Request.Referer(),
		UserAgent:    c.Request.UserAgent(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
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
	err := server.store.SaveUrlMappingToRedis(shortUrl, customizeShortUrlRequest.OriginalUrl, authPayload.UserID)
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
