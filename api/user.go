package api

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	db "github.com/maslow123/todoapp-services/db/sqlc"
	"github.com/maslow123/todoapp-services/token"
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
		if strings.Contains(err.Error(), "pq: duplicate key") {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("email-already-exists")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := newUserResponse(user)
	ctx.JSON(http.StatusOK, &resp)
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req LoginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("invalid-password")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid-password")))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(
		user.Email,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	response := LoginUserResponse{
		AccessToken: accessToken,
		User:        user,
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) me(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload == nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Unauthorized")))
		return
	}

	_, err := server.store.GetUser(context.Background(), authPayload.Username)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, errors.New("invalid-user"))
			return
		}
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, authPayload)
}
