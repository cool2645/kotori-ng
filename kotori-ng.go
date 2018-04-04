package main

import (
	"github.com/cool2645/kotori-ng/handler"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/cool2645/kotori-ng/model"
)

const (
	BaseApiVer = "/v1"
	Base       = BaseApi + BaseApiVer
)

var (
	v1Api = api.PathPrefix(BaseApiVer).Subrouter()
)

func RegRouter() {
	// Ping
	v1Api.Methods("GET").Path("").HandlerFunc(handler.Pong)
	v1Api.Methods("GET").Path("/session").HandlerFunc(handler.GetMe)
	v1Api.Methods("POST").Path("/session").HandlerFunc(handler.Login)
	v1Api.Methods("POST").Path("/users").HandlerFunc(handler.Register)
}

func InitDB(db *gorm.DB) {
	db.AutoMigrate(model.User{})
}
