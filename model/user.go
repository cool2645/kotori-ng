package model

import (
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
)

type User struct {
	ID        uint       `gorm:"AUTO_INCREMENT" json:"id"`
	UUID      string     `gorm:"not null;unique" json:"uuid"`
	Name      string     `json:"name"`
	Username  string     `gorm:"not null;unique" json:"username"`
	Password  string     `json:"password"`
	IsAdmin   bool       `json:"is_admin"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func GetUserByUsername(db *gorm.DB, username string) (user User, err error) {
	err = db.Where("username = ?", username).First(&user).Error
	if err != nil {
		err = errors.Wrap(err, "GetUserByUsername")
		return
	}
	return
}

func GetUserByUUID(db *gorm.DB, UUID string) (user User, err error) {
	err = db.Where("uuid = ?", UUID).First(&user).Error
	if err != nil {
		err = errors.Wrap(err, "GetUserByUUID")
		return
	}
	return
}

func StoreUser(db *gorm.DB, user User) (newUser User, err error) {
	err = db.Create(&user).Error
	if err != nil {
		err = errors.Wrap(err, "StoreUser")
		return
	}
	newUser = user
	return
}
