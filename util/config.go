package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver               string        `mapstructure:"DB_DRIVER"`
	DBSource               string        `mapstructure:"DB_SOURCE"`
	ServerAddress          string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey      string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration    time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	CloudinaryCloudName    string        `mapstructure:"CLOUDINARY_CLOUD_NAME"`
	CloudinaryApiKey       string        `mapstructure:"CLOUDINARY_API_KEY"`
	CloudinaryApiSecret    string        `mapstructure:"CLOUDINARY_API_SECRET"`
	CloudinaryUploadFolder string        `mapstructure:"CLOUDINARY_UPLOAD_FOLDER"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
