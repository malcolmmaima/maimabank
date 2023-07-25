package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/malcolmmaima/maimabank/db/sqlc"
	"github.com/malcolmmaima/maimabank/token"
	"github.com/malcolmmaima/maimabank/util"
)

// This server will serve all http requests for our banking service.

type Server struct {
	config util.Config
	store db.Store
	tokenMaker token.Maker
	router *gin.Engine
}

func NewServer(config util.Config,store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey) // to use JWT simply change from token.NewPasetoMaker to token.NewJWTMaker
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w ", err)
	}

	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}

	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	// Protected routes
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.GET("/accounts/statement", server.listTransfers)
	authRoutes.POST("/transfers", server.createTransfer)
	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}