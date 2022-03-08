package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/todoapp-services/token"
	"github.com/maslow123/todoapp-services/util"
	"github.com/stretchr/testify/require"
)

func TestUploadImageWithFormData(t *testing.T) {
	testCases := []struct {
		name          string
		filePath      string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "StatusUnauthorized",
			filePath: "",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:     "OK",
			filePath: "./../assets/images/simpleimage.jpg",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatch(t, recorder.Body, "success")
			},
		},
		{
			name:     "InternalServerError",
			filePath: "",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				requireBodyMatch(t, recorder.Body, "error")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			fieldName := "file"
			body := new(bytes.Buffer)

			mw := multipart.NewWriter(body)

			if tc.filePath != "" {
				file, err := os.Open(tc.filePath)
				if err != nil {
					t.Fatalf(err.Error())
				}

				w, err := mw.CreateFormFile(fieldName, tc.filePath)
				if err != nil {
					t.Fatalf(err.Error())
				}

				if _, err := io.Copy(w, file); err != nil {
					t.Fatalf(err.Error())
				}

				mw.Close()
			}

			req := httptest.NewRequest(http.MethodPost, "/file", body)
			req.Header.Add("Content-Type", mw.FormDataContentType())

			server := newTestServer(t, nil)
			recorder := httptest.NewRecorder()

			server.router.ServeHTTP(recorder, req)
		})
	}
}

func TestUploadImageWithRemoteURL(t *testing.T) {
	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "StatusUnauthorized",
			body: gin.H{
				"url": "",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "OK",
			body: gin.H{
				"url": "https://picsum.photos/200/300",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatch(t, recorder.Body, "success")
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"url": "",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				requireBodyMatch(t, recorder.Body, "error")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/remote"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatch(t *testing.T, body *bytes.Buffer, errorMessage string) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var uploadResp util.MediaDto
	err = json.Unmarshal(data, &uploadResp)
	require.NoError(t, err)
	require.Equal(t, uploadResp.Message, errorMessage)
}
