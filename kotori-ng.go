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
	v1Api.Methods("GET").Path("").HandlerFunc(handler.Pong)
	v1Api.Methods("GET").Path("/session").HandlerFunc(handler.GetMe)
	v1Api.Methods("POST").Path("/session").HandlerFunc(handler.Login)
	v1Api.Methods("POST").Path("/users").HandlerFunc(handler.Register)
	v1Api.Methods("GET").Path("/users").HandlerFunc(handler.ListUsers)
	v1Api.Methods("GET").Path("/users/username={username}").HandlerFunc(handler.GetUserByUsername)
	v1Api.Methods("GET").Path("/users/{uuid}").HandlerFunc(handler.GetUserByUUID)
	v1Api.Methods("PATCH").Path("/users/{uuid}").HandlerFunc(handler.UpdateUser)
	v1Api.Methods("PUT").Path("/users/{uuid}/username").HandlerFunc(handler.UpdateUserSetUsername)
	v1Api.Methods("PUT").Path("/users/{uuid}/password").HandlerFunc(handler.UpdateUserSetPassword)
	v1Api.Methods("PUT").Path("/users/{uuid}/admin").HandlerFunc(handler.PromoteAdmin)
	v1Api.Methods("DELETE").Path("/users/{uuid}/admin").HandlerFunc(handler.DismissAdmin)
}

func InitDB(db *gorm.DB) {
	db.AutoMigrate(model.User{})
}
