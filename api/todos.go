package api

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/maslow123/todoapp-services/db/sqlc"
	"github.com/maslow123/todoapp-services/token"
)

func (server *Server) createTodo(ctx *gin.Context) {
	var req CreateTodoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid-date")))
		return
	}

	// check category is exists or no
	_, err = server.store.GetCategory(context.Background(), req.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("invalid-category")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateTodoParams{
		UserEmail:  authPayload.Username,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Content:    req.Content,
		Date:       date,
		Color:      req.Color,
		IsPriority: *req.IsPriority,
	}

	todo, err := server.store.CreateTodo(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

func (server *Server) getTodo(ctx *gin.Context) {
	var req GetTodoRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	todo, err := server.store.GetTodo(ctx, req.TodoID)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("not-found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if todo.UserEmail != authPayload.Username {
		err := errors.New("wrong-user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

func (server *Server) listTodo(ctx *gin.Context) {
	var req ListTodoRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	var resp ListTodoResponse

	// Get Today List
	argTodayList := db.ListTodayTodoParams{
		UserEmail: authPayload.Username,
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	}
	todayTodo, err := server.store.ListTodayTodo(ctx, argTodayList)
	if err != nil {
		if err != sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	resp.Today = todayTodo

	// Get Upcoming List
	argUpcomingList := db.ListUpcomingTodoParams{
		UserEmail: authPayload.Username,
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	}
	upcomingTodo, err := server.store.ListUpcomingTodo(ctx, argUpcomingList)
	if err != nil {
		if err != sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	resp.Upcoming = upcomingTodo

	// Get Done List
	argDoneList := db.ListDoneTodoParams{
		UserEmail: authPayload.Username,
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	}
	doneTodo, err := server.store.ListDoneTodo(ctx, argDoneList)
	if err != nil {
		if err != sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	resp.Done = doneTodo

	ctx.JSON(http.StatusOK, resp)
}

func (server *Server) deleteTodo(ctx *gin.Context) {
	var req GetTodoRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err := server.store.GetTodo(ctx, req.TodoID)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("not-found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteTodo(ctx, req.TodoID)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "OK")
}

func (server *Server) updateTodo(ctx *gin.Context) {
	var req UpdateTodoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid-date")))
		return
	}

	// check todo is exists or no
	_, err = server.store.GetTodo(context.Background(), req.TodoID)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("todo-not-found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// update todo
	arg := db.UpdateTodoByUserParams{
		ID:         req.TodoID,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Content:    req.Content,
		Date:       date,
		Color:      req.Color,
		IsPriority: *req.IsPriority,
	}

	todo, err := server.store.UpdateTodoByUser(context.Background(), arg)
	if err != nil {
		log.Println(err)

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

func (server *Server) markCompleteTodo(ctx *gin.Context) {
	var req MarkCompleteTodoRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check todo is exists or no
	_, err := server.store.GetTodo(context.Background(), req.TodoID)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("todo-not-found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	todo, err := server.store.MarkAsCompleteTodo(context.Background(), req.TodoID)
	if err != nil {
		log.Println(err)

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, todo)
}
