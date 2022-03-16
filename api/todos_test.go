package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/maslow123/todoapp-services/db/mock"
	db "github.com/maslow123/todoapp-services/db/sqlc"
	"github.com/maslow123/todoapp-services/token"
	"github.com/maslow123/todoapp-services/util"
	"github.com/stretchr/testify/require"
)

func TestCreateTodoAPI(t *testing.T) {
	todo := randomTodo(t)
	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"category_id": todo.CategoryID,
				"title":       todo.Title,
				"content":     todo.Content,
				"date":        "2020-01-01",
				"color":       todo.Color,
				"is_priority": todo.IsPriority,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTodoParams{
					CategoryID: todo.CategoryID,
					Title:      todo.Title,
					Content:    todo.Content,
					Date:       todo.Date,
					Color:      todo.Color,
					IsPriority: todo.IsPriority,
					UserEmail:  todo.UserEmail,
				}
				store.EXPECT().
					GetCategory(gomock.Any(), gomock.Eq(arg.CategoryID)).
					Times(1).
					Return(db.Category{}, nil)

				store.EXPECT().
					CreateTodo(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(todo, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTodo(t, recorder.Body, todo)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"category_id": todo.CategoryID,
				"title":       todo.Title,
				"content":     todo.Content,
				"date":        "2020-01-01",
				"color":       todo.Color,
				"is_priority": todo.IsPriority,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetCategory(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateTodo(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidCategoryID",
			body: gin.H{
				"category_id": 0,
				"title":       todo.Title,
				"content":     todo.Content,
				"date":        "2020-01-01",
				"color":       todo.Color,
				"is_priority": todo.IsPriority,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCategory(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateTodo(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidTitle",
			body: gin.H{
				"category_id": todo.CategoryID,
				"title":       "",
				"content":     todo.Content,
				"date":        "2020-01-01",
				"color":       todo.Color,
				"is_priority": todo.IsPriority,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCategory(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateTodo(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidContent",
			body: gin.H{
				"category_id": todo.CategoryID,
				"title":       todo.Title,
				"content":     "",
				"date":        "2020-01-01",
				"color":       todo.Color,
				"is_priority": todo.IsPriority,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCategory(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateTodo(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidDate",
			body: gin.H{
				"category_id": todo.CategoryID,
				"title":       todo.Title,
				"content":     todo.Content,
				"date":        "invalid-date",
				"color":       todo.Color,
				"is_priority": todo.IsPriority,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCategory(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateTodo(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidColor",
			body: gin.H{
				"category_id": todo.CategoryID,
				"title":       todo.Title,
				"content":     todo.Content,
				"date":        "2020-01-01",
				"color":       "",
				"is_priority": todo.IsPriority,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCategory(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateTodo(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidIsPriority",
			body: gin.H{
				"category_id": todo.CategoryID,
				"title":       todo.Title,
				"content":     todo.Content,
				"date":        "2020-01-01",
				"color":       todo.Color,
				"is_priority": "",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCategory(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateTodo(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidIsCategoryID",
			body: gin.H{
				"category_id": 9999,
				"title":       todo.Title,
				"content":     todo.Content,
				"date":        "2020-01-01",
				"color":       todo.Color,
				"is_priority": true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCategory(gomock.Any(), gomock.Any()).
					Return(db.Category{}, sql.ErrNoRows)

				store.EXPECT().
					CreateTodo(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/todo"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListTodosAPI(t *testing.T) {
	user, _ := randomUser(t)
	n := 5
	todos := make([]db.ListTodoByUserRow, n)

	for i := 0; i < n; i++ {
		todos[i] = randomTodoWithExistingUser(t, user.Email)
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListTodoByUserParams{
					UserEmail: user.Email, // no required
					Limit:     int32(n),
					Offset:    0,
				}

				store.EXPECT().
					ListTodoByUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(todos, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTodos(t, recorder.Body, todos)
			},
		},
		{
			name: "NoAuthorization",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListTodoByUserParams{
					UserEmail: user.Email,
					Limit:     int32(n),
					Offset:    0,
				}

				store.EXPECT().
					ListTodoByUser(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Invalid PageID",
			query: Query{
				pageID:   0,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					ListTodoByUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid PageSize",
			query: Query{
				pageID:   1,
				pageSize: 99,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					ListTodoByUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/todo"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query params
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetTodo(t *testing.T) {
	todo := randomTodo(t)
	resp := db.GetTodoRow{
		CategoryID: todo.CategoryID,
		UserEmail:  todo.UserEmail,
		Title:      todo.Title,
		Content:    todo.Content,
		CreatedAt:  todo.CreatedAt,
		UpdatedAt:  todo.UpdatedAt,
		Date:       todo.Date,
		Color:      todo.Color,
		IsPriority: todo.IsPriority,
	}

	testCases := []struct {
		name          string
		todoID        int32
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			todoID: todo.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(1).
					Return(resp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				requireBodyMatchTodoRow(t, recorder.Body, resp)

			},
		},
		{
			name:   "NoAuthorization",
			todoID: todo.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)

			},
		},
		{
			name:   "NotFound",
			todoID: 999,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTodo(gomock.Any(), gomock.Any()).
					Return(db.GetTodoRow{}, sql.ErrNoRows)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/todo/%d", tc.todoID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	todo := randomTodo(t)

	testCases := []struct {
		name          string
		todoID        int32
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			todoID: todo.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(1)

				store.EXPECT().
					DeleteTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "Unauthorized",
			todoID: todo.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTodo(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					DeleteTodo(gomock.Any(), gomock.Any()).
					Times(0).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "NotFound",
			todoID: 999,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTodo(gomock.Any(), gomock.Any()).
					Return(db.GetTodoRow{}, sql.ErrNoRows)

				store.EXPECT().
					DeleteTodo(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/todo/%d", tc.todoID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateTodo(t *testing.T) {
	todo := randomTodo(t)
	todo2 := randomTodo(t)

	resp := db.Todo{
		ID:         todo.ID,
		CategoryID: todo2.CategoryID,
		Title:      todo2.Content,
		Content:    todo2.Content,
		Date:       todo2.Date,
		Color:      todo2.Color,
		IsPriority: todo2.IsPriority,
	}

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"todo_id":     todo.ID,
				"category_id": todo2.CategoryID,
				"title":       todo2.Title,
				"content":     todo2.Content,
				"date":        "2020-01-01",
				"color":       todo2.Color,
				"is_priority": todo2.IsPriority,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateTodoByUserParams{
					ID:         todo.ID,
					CategoryID: todo2.CategoryID,
					Title:      todo2.Title,
					Content:    todo2.Content,
					Date:       todo2.Date,
					Color:      todo2.Color,
					IsPriority: todo2.IsPriority,
				}

				store.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(1)

				store.EXPECT().
					UpdateTodoByUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(resp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTodo(t, recorder.Body, resp)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"todo_id":     todo.ID,
				"category_id": todo2.CategoryID,
				"title":       todo2.Title,
				"content":     todo2.Content,
				"date":        "2020-01-01",
				"color":       todo2.Color,
				"is_priority": todo2.IsPriority,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(0)

				store.EXPECT().
					UpdateTodoByUser(gomock.Any(), gomock.Any()).
					Times(0).
					Return(resp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Invalid TodoID",
			body: gin.H{
				"todo_id":     "",
				"category_id": todo2.CategoryID,
				"title":       todo2.Title,
				"content":     todo2.Content,
				"date":        "2020-01-01",
				"color":       todo2.Color,
				"is_priority": todo2.IsPriority,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(0)

				store.EXPECT().
					UpdateTodoByUser(gomock.Any(), gomock.Any()).
					Times(0).
					Return(resp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid CategoryID",
			body: gin.H{
				"todo_id":     todo.ID,
				"category_id": "",
				"title":       todo2.Title,
				"content":     todo2.Content,
				"date":        "2020-01-01",
				"color":       todo2.Color,
				"is_priority": todo2.IsPriority,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(0)

				store.EXPECT().
					UpdateTodoByUser(gomock.Any(), gomock.Any()).
					Times(0).
					Return(resp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Title",
			body: gin.H{
				"todo_id":     todo.ID,
				"category_id": todo2.CategoryID,
				"title":       "",
				"content":     todo2.Content,
				"date":        "2020-01-01",
				"color":       todo2.Color,
				"is_priority": todo2.IsPriority,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(0)

				store.EXPECT().
					UpdateTodoByUser(gomock.Any(), gomock.Any()).
					Times(0).
					Return(resp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Content",
			body: gin.H{
				"todo_id":     todo.ID,
				"category_id": todo2.CategoryID,
				"title":       todo2.Title,
				"content":     "",
				"date":        "2020-01-01",
				"color":       todo2.Color,
				"is_priority": todo2.IsPriority,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(0)

				store.EXPECT().
					UpdateTodoByUser(gomock.Any(), gomock.Any()).
					Times(0).
					Return(resp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Date Format",
			body: gin.H{
				"todo_id":     todo.ID,
				"category_id": todo2.CategoryID,
				"title":       todo2.Title,
				"content":     todo2.Content,
				"date":        "invalid-date",
				"color":       todo2.Color,
				"is_priority": todo2.IsPriority,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(0)

				store.EXPECT().
					UpdateTodoByUser(gomock.Any(), gomock.Any()).
					Times(0).
					Return(resp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Date",
			body: gin.H{
				"todo_id":     todo.ID,
				"category_id": todo2.CategoryID,
				"title":       todo2.Title,
				"content":     todo2.Content,
				"date":        "",
				"color":       todo2.Color,
				"is_priority": todo2.IsPriority,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(0)

				store.EXPECT().
					UpdateTodoByUser(gomock.Any(), gomock.Any()).
					Times(0).
					Return(resp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Color",
			body: gin.H{
				"todo_id":     todo.ID,
				"category_id": todo2.CategoryID,
				"title":       todo2.Title,
				"content":     todo2.Content,
				"date":        "2020-01-01",
				"color":       "",
				"is_priority": todo2.IsPriority,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(0)

				store.EXPECT().
					UpdateTodoByUser(gomock.Any(), gomock.Any()).
					Times(0).
					Return(resp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid IsPriority",
			body: gin.H{
				"todo_id":     todo.ID,
				"category_id": todo2.CategoryID,
				"title":       todo2.Title,
				"content":     todo2.Content,
				"date":        "2020-01-01",
				"color":       todo2.Color,
				"is_priority": "",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(0)

				store.EXPECT().
					UpdateTodoByUser(gomock.Any(), gomock.Any()).
					Times(0).
					Return(resp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Todo NotFound",
			body: gin.H{
				"todo_id":     999,
				"category_id": todo2.CategoryID,
				"title":       todo2.Title,
				"content":     todo2.Content,
				"date":        "2020-01-01",
				"color":       todo2.Color,
				"is_priority": todo2.IsPriority,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, todo.UserEmail, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTodo(gomock.Any(), gomock.Any()).
					Return(db.GetTodoRow{}, sql.ErrNoRows)

				store.EXPECT().
					UpdateTodoByUser(gomock.Any(), gomock.Any()).
					Times(0).
					Return(resp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/todo")
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchTodo(t *testing.T, body *bytes.Buffer, todo db.Todo) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotTodo db.Todo
	err = json.Unmarshal(data, &gotTodo)

	require.NoError(t, err)
	require.Equal(t, todo.CategoryID, gotTodo.CategoryID)
	require.Equal(t, todo.Title, gotTodo.Title)
	require.Equal(t, todo.Content, gotTodo.Content)
	require.Equal(t, todo.Date, gotTodo.Date)
	require.Equal(t, todo.Color, gotTodo.Color)
	require.Equal(t, todo.IsPriority, gotTodo.IsPriority)
}

func requireBodyMatchTodoRow(t *testing.T, body *bytes.Buffer, todo db.GetTodoRow) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotTodo db.GetTodoRow
	err = json.Unmarshal(data, &gotTodo)

	require.NoError(t, err)
	require.Equal(t, todo, gotTodo)
}

func requireBodyMatchTodos(t *testing.T, body *bytes.Buffer, todo []db.ListTodoByUserRow) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotTodo []db.ListTodoByUserRow
	err = json.Unmarshal(data, &gotTodo)

	require.NoError(t, err)
	require.Equal(t, todo, gotTodo)
}

func randomTodo(t *testing.T) db.Todo {
	user, _ := randomUser(t)
	category := randomCategory()
	date, err := time.Parse("2006-01-02", "2020-01-01")
	require.NoError(t, err)

	todo := db.Todo{
		ID:         int32(util.RandomInt(1, 100)),
		CategoryID: category.ID,
		Title:      fmt.Sprintf("This title from testing: %s", util.RandomString(10)),
		Content:    fmt.Sprintf("This content from testing: %s", util.RandomString(30)),
		Date:       date,
		Color:      util.RandomColor(),
		IsPriority: false,
		UserEmail:  user.Email,
	}

	return todo
}

func randomTodoWithExistingUser(t *testing.T, userEmail string) db.ListTodoByUserRow {
	category := randomCategory()
	date, err := time.Parse("2006-01-02", "2020-01-01")
	require.NoError(t, err)

	todo := db.ListTodoByUserRow{
		CategoryID:   category.ID,
		Title:        fmt.Sprintf("This title from testing: %s", util.RandomString(10)),
		Content:      fmt.Sprintf("This content from testing: %s", util.RandomString(30)),
		Date:         date,
		Color:        util.RandomColor(),
		IsPriority:   false,
		UserEmail:    userEmail,
		CategoryName: sql.NullString{String: category.Name},
	}

	return todo
}
