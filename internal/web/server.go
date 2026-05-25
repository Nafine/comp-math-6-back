package web

import (
	"comp-math-6/internal/config"
	"comp-math-6/internal/web/handler"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	address string
	router  *gin.Engine
}

func New(cfg *config.Config) *Server {
	router := gin.New()

	router.Use(gin.Recovery())

	router.GET("/equations", handler.Equations())
	router.POST("/solve", handler.Solve())

	return &Server{
		address: fmt.Sprintf("%v:%v", cfg.Host, cfg.Port),
		router:  router,
	}
}

func (s *Server) Start() error {
	return s.router.Run(s.address)
}
