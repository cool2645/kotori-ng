package status

import (
	. "github.com/cool2645/kotori-ng/kotoriplugin"
	"github.com/cool2645/kotori-ng/pluginmanager"
)

func GetPluginInfo() (pis []PluginInfo) {
	for _, p := range pluginmanager.PM.Plugins {
		pis = append(pis, p.Info)
	}
	return
}
