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
	viper.SetDefault("database.server", "mariadb")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.user", "drlm3")
	viper.SetDefault("database.password", "drlm3db")
	viper.SetDefault("database.database", "drlm3")
}

func InitDatabase() {

	if DBConn != nil {
		return
	}

	connectionString := (Config.Database.User + ":" + Config.Database.Password + "@tcp(" + Config.Database.Server + ":" + Config.Database.Port + ")/" + Config.Database.Database + "?parseTime=true")

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
