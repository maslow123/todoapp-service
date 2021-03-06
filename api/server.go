package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	db "github.com/maslow123/todoapp-services/db/sqlc"
	"github.com/maslow123/todoapp-services/token"
	"github.com/maslow123/todoapp-services/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("cannot-create-token")
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

	// Users
	router.Use(CORSMiddleware())
	router.POST("/users/register", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker), CORSMiddleware())
	authRoutes.GET("/users/me", server.me)
	// Category
	authRoutes.POST("/categories", server.createCategory)
	authRoutes.GET("/categories", server.listCategories)
	authRoutes.PATCH("/categories", server.updateCategory)
	authRoutes.DELETE("/categories/:category_id", server.deleteCategory)

	// Todo
	authRoutes.POST("/todo", server.createTodo)
	authRoutes.GET("/todo", server.listTodo)
	authRoutes.GET("/todo/:todo_id", server.getTodo)
	authRoutes.DELETE("/todo/:todo_id", server.deleteTodo)
	authRoutes.PUT("/todo", server.updateTodo)
	authRoutes.PUT("/todo/:todo_id", server.markCompleteTodo)

	// Upload
	authRoutes.POST("/file", server.UpdateUserPhoto)
	authRoutes.POST("/remote", RemoteUpload())
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
