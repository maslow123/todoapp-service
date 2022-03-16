package db

import (
	"context"
	"testing"
	"time"

	"github.com/maslow123/todoapp-services/util"
	"github.com/stretchr/testify/require"
)

func createRandomTodo(t *testing.T, userEmail string, categoryID int32) Todo {
	title := "Todo title 1"
	content := "Todo content 1"
	date, err := time.Parse("2006-01-02", "2021-10-29")
	require.NoError(t, err)

	color := util.RandomColor()

	arg := CreateTodoParams{
		CategoryID: categoryID,
		UserEmail:  userEmail,
		Title:      title,
		Content:    content,
		Date:       date,
		Color:      color,
		IsPriority: false,
	}

	todo, err := testQueries.CreateTodo(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, todo)

	require.Equal(t, todo.CategoryID, arg.CategoryID)
	require.Equal(t, todo.Title, arg.Title)
	require.Equal(t, todo.Content, arg.Content)
	require.Equal(t, todo.Date.Unix(), arg.Date.Unix())
	require.Equal(t, todo.Color, arg.Color)
	require.Equal(t, todo.IsPriority, arg.IsPriority)

	return todo
}

func TestCreateTodo(t *testing.T) {
	user := createRandomUser(t)
	category := createRandomCategory(t)

	createRandomTodo(t, user.Email, category.ID)
}

func TestListTodoByUser(t *testing.T) {
	user := createRandomUser(t)
	category := createRandomCategory(t)
	// Create 10 todo per user
	for i := 0; i < 10; i++ {
		createRandomTodo(t, user.Email, category.ID)
	}

	arg := ListTodoByUserParams{
		Limit:     10,
		Offset:    0,
		UserEmail: user.Email,
	}
	todos, err := testQueries.ListTodoByUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, todos)

	require.Equal(t, category.ID, todos[0].CategoryID)
	require.Equal(t, category.Name, todos[0].CategoryName.String)
	require.Equal(t, 10, len(todos))
}

func TestUpdateTodoByUser(t *testing.T) {
	user := createRandomUser(t)
	category := createRandomCategory(t)
	date, err := time.Parse("2006-01-02", "2021-10-28")
	require.NoError(t, err)

	color := util.RandomColor()

	todo1 := createRandomTodo(t, user.Email, category.ID)

	arg := UpdateTodoByUserParams{
		ID:         todo1.ID,
		Title:      "Update title 1",
		Content:    "Update content 1",
		Date:       date,
		Color:      color,
		IsPriority: true,
	}

	todo2, err := testQueries.UpdateTodoByUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, todo2)

	require.Equal(t, arg.ID, todo2.ID)
	require.Equal(t, arg.Title, todo2.Title)
	require.Equal(t, arg.Content, todo2.Content)
	require.Equal(t, arg.Date.Unix(), todo2.Date.Unix())
	require.Equal(t, arg.Color, todo2.Color)
	require.Equal(t, arg.IsPriority, todo2.IsPriority)

	require.NotZero(t, todo2.UpdatedAt)
}

func TestDeleteTodo(t *testing.T) {
	user := createRandomUser(t)
	category := createRandomCategory(t)
	todo := createRandomTodo(t, user.Email, category.ID)

	err := testQueries.DeleteTodo(context.Background(), todo.ID)
	require.NoError(t, err)
}
