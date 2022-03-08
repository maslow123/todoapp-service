package util

import (
	"context"
	"log"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func ImageUploadHelper(input interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := LoadConfig("./..")
	if err != nil {
		log.Fatal("cannot-load-config", err)
	}
	// create cloudinary instance
	cld, err := cloudinary.NewFromParams(
		config.CloudinaryCloudName,
		config.CloudinaryApiKey,
		config.CloudinaryApiSecret,
	)
	if err != nil {
		return "", err
	}

	// upload file
	uploadParam, err := cld.Upload.Upload(ctx, input, uploader.UploadParams{
		Folder: config.CloudinaryUploadFolder,
	})
	if err != nil {
		return "", err
	}

	return uploadParam.SecureURL, nil
}
