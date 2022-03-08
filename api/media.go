package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/todoapp-services/models"
	"github.com/maslow123/todoapp-services/util"
)

func FileUpload() gin.HandlerFunc {
	return func(c *gin.Context) {
		// upload
		formFile, _, err := c.Request.FormFile("file")
		if err != nil {
			util.GenericUploadResponse(c, http.StatusInternalServerError, "error", map[string]interface{}{"data": "Select a file to upload"})
			return
		}

		uploadUrl, err := models.NewMediaUpload().FileUpload(models.File{File: formFile})
		if err != nil {
			util.GenericUploadResponse(c, http.StatusInternalServerError, "error", map[string]interface{}{"data": "Error uploading file"})
			return
		}

		util.GenericUploadResponse(c, http.StatusOK, "success", map[string]interface{}{"data": uploadUrl})
	}
}

func RemoteUpload() gin.HandlerFunc {
	return func(c *gin.Context) {
		var url models.Url

		// validate the request body
		if err := c.BindJSON(&url); err != nil {
			util.GenericUploadResponse(c, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		uploadUrl, err := models.NewMediaUpload().RemoteUpload(url)
		if err != nil {

			util.GenericUploadResponse(c, http.StatusInternalServerError, "error", map[string]interface{}{"data": "Error uploading file"})
			return
		}

		util.GenericUploadResponse(c, http.StatusOK, "success", map[string]interface{}{"data": uploadUrl})
	}
}
