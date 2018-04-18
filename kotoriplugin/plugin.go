package kotoriplugin

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Plugin interface {
	GetPluginInfo() PluginInfo
	LoadConfig() error
	RegRouter(*mux.Router) error
	InitDB(*gorm.DB) error
}

type PluginInfo struct {
	Name    string `json:"name"`
	Author  string `json:"author"`
	Version string `json:"version"`
	License string `json:"license"`
	URL     string `json:"url"`
}
