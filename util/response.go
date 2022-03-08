package util

import (
	"github.com/gin-gonic/gin"
)

type MediaDto struct {
	StatusCode int                    `json:"status_code"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"data"`
}

func GenericUploadResponse(
	c *gin.Context,
	code int,
	message string,
	data map[string]interface{},

) {
	c.JSON(
		code,
		MediaDto{
			StatusCode: code,
			Message:    message,
			Data:       data,
		},
	)
}
