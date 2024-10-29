package handler

import (
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shienlee73/url-shortener/frontend"
	"github.com/shienlee73/url-shortener/rate_limiter"
	"github.com/shienlee73/url-shortener/store"
)

type Server struct {
	address     string
	router      *gin.Engine
	store       store.StorageService
	rateLimiter rate_limiter.RateLimiter
}

func NewServer(store *store.StorageService, rateLimiter *rate_limiter.RateLimiter) *Server {
	server := &Server{
		store:       *store,
		rateLimiter: *rateLimiter,
	}
	server.setupRoutes()
	return server
}

func (server *Server) setupRoutes() {
	r := gin.Default()

	contentStatic, err := fs.Sub(frontend.Assets(), "dist")
	if err != nil {
		panic(err)
	}

	r.StaticFS("/url-shortener", http.FS(contentStatic))

	r.GET("/", server.Index)
	r.POST("/shorten", server.rateLimiter.Limit("shorten", time.Minute, 5), server.CreateShortUrl)
	r.POST("/customize", server.rateLimiter.Limit("customize", time.Minute, 5), server.CustomizeShortUrl)
	r.GET("/:shortUrl", server.rateLimiter.Limit("/:shortUrl", time.Minute, 10), server.HandleShortUrlRedirect)

	server.router = r
}

func (server *Server) Start(address string) error {
	fmt.Println("Hello URL Shortener! 🚀")
	server.address = address
	return server.router.Run(address)
}
