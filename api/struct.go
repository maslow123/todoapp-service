package api

import db "github.com/maslow123/todoapp-services/db/sqlc"

// Category
type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color" binding:"required"`
}

type ListCategoryRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type UpdateCategoryRequest struct {
	CategoryID int32  `json:"category_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Color      string `json:"color" binding:"required"`
}

type DeleteCategoryRequest struct {
	CategoryID int32 `uri:"category_id" binding:"required,min=1"`
}

// User
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Pic      string `json:"pic" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
	AccessToken string `json:"access_token"`
	User        db.User
}

type GenericUserResponse struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Pic     string `json:"pic"`
	Email   string `json:"email"`
}

// Todo
type CreateTodoRequest struct {
	CategoryID int32  `json:"category_id" binding:"required,min=1"`
	Title      string `json:"title"`
	Content    string `json:"content"`
}

type GetTodoRequest struct {
	TodoID int32 `uri:"todo_id" binding:"required,min=1"`
}

type ListTodoRequest struct {
	PageID    int32  `form:"page_id" binding:"required,min=1"`
	PageSize  int32  `form:"page_size" binding:"required,min=5,max=10"`
	UserEmail string `form:"user_email" binding:"required"`
}
