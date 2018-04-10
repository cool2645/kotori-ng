package model

import (
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/cool2645/kotori-ng/database"
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

func ListUsers(db *gorm.DB, page uint, perPage uint, orderBy string, order string) (users []User, total uint, err error) {
	if perPage == 0 {
		perPage = 15
	}
	err = db.Order(orderBy + " " + order).Limit(perPage).Offset((page - 1) * perPage).Find(&users).Error
	if err != nil {
		err = errors.Wrap(err, "ListUsers")
		return
	}
	err = db.Model(&User{}).Count(&total).Error
	if err != nil {
		err = errors.Wrap(err, "ListUsers")
		return
	}
	return
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

func StoreUser(db *gorm.DB, user *User) (err error) {
	err = db.Create(user).Error
	if err != nil {
		err = errors.Wrap(err, "StoreUser")
		return
	}
	return
}

func MakeUser(user *User) (err error) {
	u2, err := uuid.NewV4()
	if err != nil {
		err = errors.Wrap(err, "MakeUser")
		return
	}
	user.UUID = u2.String()
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		err = errors.Wrap(err, "MakeUser")
		return
	}
	user.Password = string(hash)
	user.IsAdmin = false
	err = StoreUser(database.DB, user)
	if err != nil {
		err = errors.Wrap(err, "MakeUser")
		return
	}
	return
}

func MakeAdmin(user *User) (err error) {
	u2, err := uuid.NewV4()
	if err != nil {
		err = errors.Wrap(err, "MakeAdmin")
		return
	}
	user.UUID = u2.String()
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		err = errors.Wrap(err, "MakeAdmin")
		return
	}
	user.Password = string(hash)
	user.IsAdmin = false
	err = StoreUser(database.DB, user)
	if err != nil {
		err = errors.Wrap(err, "MakeAdmin")
		return
	}
	return
}

func UpdateUser(db *gorm.DB, user *User, patch map[string]interface{}) (err error) {
	err = db.Model(user).Select("name").Updates(patch).Error
	if err != nil {
		err = errors.Wrap(err, "UpdateUser")
		return
	}
	return
}

func UpdateUserSetUsername(db *gorm.DB, user *User, username string) (err error) {
	err = db.Model(user).UpdateColumn("username", username).Error
	if err != nil {
		err = errors.Wrap(err, "UpdateUserSetUsername")
		return
	}
	return
}

func UpdateUserSetPassword(db *gorm.DB, user *User, password string) (err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		err = errors.Wrap(err, "UpdateUserSetPassword")
		return
	}
	password = string(hash)
	err = db.Model(user).UpdateColumn("password", password).Error
	if err != nil {
		err = errors.Wrap(err, "UpdateUserSetPassword")
		return
	}
	return
}

func PromoteUserToAdmin(db *gorm.DB, user *User) (err error) {
	user.IsAdmin = true
	db.Save(user)
	if err != nil {
		err = errors.Wrap(err, "PromoteUserToAdmin")
		return
	}
	return
}

func DismissUserFromAdmin(db *gorm.DB, user *User) (err error) {
	user.IsAdmin = false
	db.Save(user)
	if err != nil {
		err = errors.Wrap(err, "DismissUserFromAdmin")
		return
	}
	return
}
