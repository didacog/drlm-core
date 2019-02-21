package lib

import (
	"fmt"
	"os"
	"strings"

	"github.com/brainupdaters/drlm-comm/logger"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type DrlmcoreConfig struct {
	Database DatabaseConfig       `mapstructure:"database"`
	Minio    MinioConfig          `mapstructure:"minio"`
	Drlmapi  DrlmapiConfig        `mapstructure:"drlmapi"`
	Logging  logger.LoggingConfig `mapstructure:"logging"`
}

var Config *DrlmcoreConfig

func InitConfig(c string) {
	if c != "" {
		// Use config file from the flag.
		viper.SetConfigFile(c)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".drlm-core" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".drlm-core")
	}

	// Enable environment variables
	// ex.: DRLMCORE_DRLMAPI_PORT=8000
	viper.SetEnvPrefix("DRLMCORE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// var config Config

	err := viper.Unmarshal(&Config)
	if err != nil {
		panic("Unable to unmarshal config")
	}
}
