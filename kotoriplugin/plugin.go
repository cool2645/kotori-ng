package kotoriplugin

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Plugin interface {
	GetName() string
	GetVersion() string
	LoadConfig() error
	RegRouter(*mux.Router) error
	InitDB(*gorm.DB) error
}
