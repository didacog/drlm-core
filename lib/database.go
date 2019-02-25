package lib

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Server   string
	Port     string
	User     string
	Password string
	Database string
}

var DBConn *gorm.DB

func SetDatabaseConfigDefaults() {
	viper.SetDefault("database.server", "localhost")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.user", "drlm3")
	viper.SetDefault("database.password", "drlm3db")
	viper.SetDefault("database.database", "drlm3")
}

func InitDatabase(cfg DatabaseConfig) {

	if DBConn != nil {
		return
	}

	connectionString := (cfg.User + ":" + cfg.Password + "@tcp(" + cfg.Server + ":" + cfg.Port + ")/" + cfg.Database + "?parseTime=true")

	db, err := gorm.Open("mysql", connectionString)
	if err != nil {
		log.Panic("failed to connect database")
	}
	log.Info("Connected to database!")

	DBConn = db

	InitUser()
}

func closeDatabase() {
	DBConn.Close()
}
