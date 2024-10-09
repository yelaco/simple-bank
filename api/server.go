package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/yelaco/simple-bank/db/sqlc"
	"github.com/yelaco/simple-bank/token"
	"github.com/yelaco/simple-bank/util"
)

// Server serves HTTP requests for our banking service
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("api.NewServer: cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
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

	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/", server.createUser)
		userRoutes.POST("/login", server.loginUser)
	}

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
