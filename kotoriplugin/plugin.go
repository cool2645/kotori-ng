package kotoriplugin

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Plugin interface {
	GetName() string
	GetVersion() string
	RegRouter(*mux.Router)
	InitDB(*gorm.DB)
}
