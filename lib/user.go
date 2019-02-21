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

func (u *User) SetPassword(pass string) {
	u.Password = pass
}

func (u *User) SetUser(user string) {
	u.User = user
}

func (u *User) LoadUser(user string) {
	DBConn.Where("user = ?", user).First(u)
}

func (u *User) LoadId(id string) {
	DBConn.First(u, id)
}

func (u *User) Save() {
	DBConn.Model(u).Update(u.User, u.Password)
}

func (u *User) Delete() {
	DBConn.Delete(u)
}
