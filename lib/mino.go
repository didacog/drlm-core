package lib

import "github.com/spf13/viper"

type MinioConfig struct {
	Server     string
	Port       string
	Ssl        bool
	Access_key string
	Secret_key string
}

func SetMinioConfigDefaults() {
	viper.SetDefault("minio.server", "minio")
	viper.SetDefault("minio.port", "9000")
	viper.SetDefault("minio.ssl", false)
	viper.SetDefault("minio.access_key", "drlm3minio")
	viper.SetDefault("minio.secret_key", "drlm3minio")
}
