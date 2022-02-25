package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/maslow123/todoapp-services/util"
	"github.com/stretchr/testify/require"
)

func createRandomCategory(t *testing.T) Category {
	color := fmt.Sprintf("#%s", util.RandomString(5))

	arg := CreateCategoryParams{
		Name:  util.RandomString(6),
		Color: color,
	}

	category, err := testQueries.CreateCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, category)

	require.Equal(t, arg.Name, category.Name)
	require.Equal(t, arg.Color, category.Color)

	require.NotZero(t, category.CreatedAt)

	return category
}

func TestCreateCategory(t *testing.T) {
	createRandomCategory(t)
}

func TestListCategories(t *testing.T) {
	// Create 10 data
	for i := 0; i < 10; i++ {
		createRandomCategory(t)
	}

	arg := ListCategoriesParams{
		Limit:  5,
		Offset: 0,
	}

	categories, err := testQueries.ListCategories(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, categories)
	require.Equal(t, 5, len(categories))

	for _, category := range categories {
		require.NotEmpty(t, category)
	}
}

func TestUpdateCategory(t *testing.T) {
	category1 := createRandomCategory(t)

	arg := UpdateCategoryParams{
		ID:    category1.ID,
		Color: fmt.Sprintf("#%s", util.RandomString(6)),
		Name:  util.RandomString(6),
	}

	category2, err := testQueries.UpdateCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, category2)

	require.Equal(t, arg.ID, category2.ID)
	require.Equal(t, arg.Color, category2.Color)
	require.Equal(t, arg.Name, category2.Name)
}

func TestDeleteCategory(t *testing.T) {
	category1 := createRandomCategory(t)
	err := testQueries.DeleteCategory(context.Background(), category1.ID)
	require.NoError(t, err)
}
