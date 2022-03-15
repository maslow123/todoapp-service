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

func randomTodo(t *testing.T) db.Todo {
	user, _ := randomUser(t)
	category := randomCategory()
	date, err := time.Parse("2006-01-02", "2020-01-01")
	require.NoError(t, err)

	todo := db.Todo{
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
