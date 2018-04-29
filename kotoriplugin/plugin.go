package kotoriplugin

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/cool2645/kotori-ng/version"
)

type Plugin interface {
	GetPluginInfo() PluginInfo
	LoadConfig() error
	RegRouter(*mux.Router) error
	InitDB(*gorm.DB) error
}

type PluginInfo struct {
	BasicInfo BasicInfo `json:"basic_info"`
	BuildInfo BuildInfo `json:"build_info"`
}

type BasicInfo version.VersionInfo

type BuildInfo version.BuildInfo
