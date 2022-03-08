package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/maslow123/todoapp-services/db/sqlc"
	"github.com/maslow123/todoapp-services/models"
	"github.com/maslow123/todoapp-services/token"
	"github.com/maslow123/todoapp-services/util"
)

func (server *Server) UpdateUserPhoto(ctx *gin.Context) {
	// upload
	formFile, _, err := ctx.Request.FormFile("file")
	if err != nil {
		// util.GenericUploadResponse(ctx, http.StatusInternalServerError, "error", map[string]interface{}{"data": "Select a file to upload"})
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	uploadUrl, err := models.NewMediaUpload().FileUpload(models.File{File: formFile})
	if err != nil {
		// util.GenericUploadResponse(ctx, http.StatusInternalServerError, "error", map[string]interface{}{"data": "Error uploading file"})
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.UpdateUserPhotoParams{
		Email: authPayload.Username,
		Pic:   uploadUrl,
	}

	user, err := server.store.UpdateUserPhoto(ctx, arg)
	if err != nil {
		// log.Println("error: ", err)
		// util.GenericUploadResponse(ctx, http.StatusInternalServerError, "error", map[string]interface{}{"data": "Error update user"})
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// util.GenericUploadResponse(ctx, http.StatusOK, "success", map[string]interface{}{"data": uploadUrl})
	ctx.JSON(http.StatusOK, newUserResponse(user))
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
