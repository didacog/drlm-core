package lib

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	User     string
	Password string
}

func InitUser() {
	DBConn.AutoMigrate(&User{})

}

func (u *User) AddUser() {
	DBConn.Create(u)
}
