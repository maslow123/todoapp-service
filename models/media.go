package models

import (
	"mime/multipart"

	"github.com/go-playground/validator/v10"
	"github.com/maslow123/todoapp-services/util"
)

var (
	validate = validator.New()
)

type File struct {
	File multipart.File `json:"file,omitempty" validate:"required"`
}

type Url struct {
	Url string `json:"url,omitempty" validate:"required"`
}

type mediaUpload interface {
	FileUpload(file File) (string, error)
	RemoteUpload(url Url) (string, error)
}

type media struct{}

func NewMediaUpload() mediaUpload {
	return &media{}
}

func (*media) FileUpload(file File) (string, error) {
	// validate
	err := validate.Struct(file)
	if err != nil {
		return "", err
	}

	// upload
	uploadUrl, err := util.ImageUploadHelper(file.File)
	if err != nil {
		return "", err
	}

	return uploadUrl, nil
}

func (*media) RemoteUpload(url Url) (string, error) {
	// validate
	err := validate.Struct(url)
	if err != nil {
		return "", err
	}

	// upload
	uploadUrl, errUrl := util.ImageUploadHelper(url.Url)
	if errUrl != nil {
		return "", err
	}

	return uploadUrl, nil
}
