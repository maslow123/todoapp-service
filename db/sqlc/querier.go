// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
)

type Querier interface {
	CreateCategory(ctx context.Context, name string) (Category, error)
	CreateTodo(ctx context.Context, arg CreateTodoParams) (Todo, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteCategory(ctx context.Context, id int32) error
	DeleteTodo(ctx context.Context, id int32) error
	DeleteUser(ctx context.Context, id int32) error
	GetCategory(ctx context.Context, id int32) (Category, error)
	GetTodo(ctx context.Context, id int32) (GetTodoRow, error)
	GetUser(ctx context.Context, email string) (User, error)
	ListCategories(ctx context.Context, arg ListCategoriesParams) ([]Category, error)
	ListDoneTodo(ctx context.Context, arg ListDoneTodoParams) ([]ListDoneTodoRow, error)
	ListTodayTodo(ctx context.Context, arg ListTodayTodoParams) ([]ListTodayTodoRow, error)
	ListTodoByUser(ctx context.Context, arg ListTodoByUserParams) ([]ListTodoByUserRow, error)
	ListUpcomingTodo(ctx context.Context, arg ListUpcomingTodoParams) ([]ListUpcomingTodoRow, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateCategory(ctx context.Context, arg UpdateCategoryParams) (Category, error)
	UpdateTodoByUser(ctx context.Context, arg UpdateTodoByUserParams) (Todo, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateUserPhoto(ctx context.Context, arg UpdateUserPhotoParams) (User, error)
}

var _ Querier = (*Queries)(nil)
