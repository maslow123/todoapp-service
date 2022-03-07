// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
)

type Querier interface {
	CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error)
	CreateTodo(ctx context.Context, arg CreateTodoParams) (Todo, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteCategory(ctx context.Context, id int32) error
	DeleteTodo(ctx context.Context, id int32) error
	DeleteUser(ctx context.Context, id int32) error
	GetTodo(ctx context.Context, id int32) (GetTodoRow, error)
	GetUser(ctx context.Context, email string) (User, error)
	ListCategories(ctx context.Context, arg ListCategoriesParams) ([]Category, error)
	ListTodoByUser(ctx context.Context, arg ListTodoByUserParams) ([]ListTodoByUserRow, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateCategory(ctx context.Context, arg UpdateCategoryParams) (Category, error)
	UpdateTodoByUser(ctx context.Context, arg UpdateTodoByUserParams) (Todo, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
