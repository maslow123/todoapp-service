package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/maslow123/todoapp-services/db/sqlc"
)

func (server *Server) createCategory(ctx *gin.Context) {
	var req CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateCategoryParams{
		Name:  req.Name,
		Color: req.Color,
	}

	category, err := server.store.CreateCategory(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}
