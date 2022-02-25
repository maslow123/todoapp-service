package api

import db "github.com/maslow123/todoapp-services/db/sqlc"

// Category
type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color" binding:"required"`
}

// User
type CreateUserRequest struct {
	Name     string `json:"name" binding="required"`
	Address  string `json:"address" binding="required"`
	Pic      string `json:"pic" binding="required"`
	Password string `json:"password" binding="required"`
	Email    string `json:"email" binding="required"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding="required"`
	Password string `json:"password" binding="required"`
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
