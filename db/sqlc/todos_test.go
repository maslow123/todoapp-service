package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomTodo(t *testing.T, userID int32, categoryID int32) Todo {
	title := "Todo title 1"
	content := "Todo content 1"

	arg := CreateTodoParams{
		CategoryID: categoryID,
		UserID:     userID,
		Title:      title,
		Content:    content,
	}

	todo, err := testQueries.CreateTodo(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, todo)

	require.Equal(t, todo.CategoryID, arg.CategoryID)
	require.Equal(t, todo.UserID, arg.UserID)
	require.Equal(t, todo.Title, arg.Title)
	require.Equal(t, todo.Content, arg.Content)

	return todo
}

func TestCreateTodo(t *testing.T) {
	user := createRandomUser(t)
	category := createRandomCategory(t)

	createRandomTodo(t, user.ID, category.ID)
}

func TestListTodoByUser(t *testing.T) {
	user := createRandomUser(t)
	category := createRandomCategory(t)
	// Create 10 todo per user
	for i := 0; i < 10; i++ {
		createRandomTodo(t, user.ID, category.ID)
	}

	arg := ListTodoByUserParams{
		Limit:  10,
		Offset: 0,
		UserID: user.ID,
	}
	todos, err := testQueries.ListTodoByUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, todos)

	require.Equal(t, user.ID, todos[0].UserID)
	require.Equal(t, category.ID, todos[0].CategoryID)
	require.Equal(t, category.Name, todos[0].CategoryName.String)
	require.Equal(t, category.Color, todos[0].CategoryColor.String)
	require.Equal(t, 10, len(todos))
}

func TestUpdateTodoByUser(t *testing.T) {
	user := createRandomUser(t)
	category := createRandomCategory(t)

	todo1 := createRandomTodo(t, user.ID, category.ID)

	arg := UpdateTodoByUserParams{
		ID:      todo1.ID,
		Title:   "Update title 1",
		Content: "Update content 1",
		UserID:  user.ID,
	}

	todo2, err := testQueries.UpdateTodoByUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, todo2)

	require.Equal(t, arg.ID, todo2.ID)
	require.Equal(t, arg.Title, todo2.Title)
	require.Equal(t, arg.Content, todo2.Content)
	require.Equal(t, arg.UserID, todo2.UserID)

	require.NotZero(t, todo2.UpdatedAt)
}

func TestDeleteTodo(t *testing.T) {
	user := createRandomUser(t)
	category := createRandomCategory(t)
	todo := createRandomTodo(t, user.ID, category.ID)

	err := testQueries.DeleteTodo(context.Background(), todo.ID)
	require.NoError(t, err)
}
