package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/maslow123/todoapp-services/db/sqlc"
	"github.com/maslow123/todoapp-services/util"
)

func newUserResponse(user db.User) GenericUserResponse {
	return GenericUserResponse{
		Name:    user.Name,
		Address: user.Address,
		Pic:     user.Pic,
		Email:   user.Email,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Name:           req.Name,
		Address:        req.Address,
		Pic:            req.Pic,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := newUserResponse(user)
	ctx.JSON(http.StatusOK, &resp)
}
