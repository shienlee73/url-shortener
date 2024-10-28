package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shienlee73/url-shortener/store"
)

type Server struct {
	address string
	router *gin.Engine
	store  store.StorageService
}

func NewServer(store *store.StorageService) *Server {
	server := &Server{store: *store}
	server.setupRoutes()
	return server
}

func (server *Server) setupRoutes() {
	r := gin.Default()

	r.StaticFile("/", "./frontend/index.html")

	// r.GET("/", server.Index)
	r.POST("/shorten", server.CreateShortUrl)
	r.GET("/:shortUrl", server.HandleShortUrlRedirect)

	server.router = r
}

func (server *Server) Start(address string) error {
	fmt.Println("Hello URL Shortener! 🚀")
	server.address = address
	return server.router.Run(address)
}
