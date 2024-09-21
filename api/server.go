package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/yelaco/simple-bank/db/sqlc"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	_ = server.bindValidators()

	accountRoutes := router.Group("/accounts")
	{
		accountRoutes.POST("/", server.createAccount)
		accountRoutes.GET("/", server.listAccounts)
		accountRoutes.GET("/:id", server.getAccount)
	}

	transferRoutes := router.Group("/transfers")
	{
		transferRoutes.POST("/", server.createTransfer)
	}

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
