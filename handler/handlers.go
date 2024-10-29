package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shienlee73/url-shortener/shortener"
)

func (server *Server) Index(c *gin.Context) {
	c.Redirect(http.StatusFound, "/url-shortener/")
}

type CreateShortUrlRequest struct {
	OriginalUrl string `json:"originalUrl"`
	UserId      string `json:"userId"`
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

	err = server.store.SaveUrlMapping(shortUrl, createShortUrlRequest.OriginalUrl, createShortUrlRequest.UserId)
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

	originalUrl, err := server.store.RetrieveOriginalUrl(shortUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, originalUrl)
}
