package handler

import (
	"net/http"
	"github.com/cool2645/kotori-ng/status"
	"github.com/cool2645/kotori-ng/httputils"
	"github.com/cool2645/kotori-ng/pluginmanager"
	. "github.com/cool2645/kotori-ng/kotoriplugin"
	"github.com/cool2645/kotori-ng/version"
)

func GetStatus(w http.ResponseWriter, req *http.Request) {
	hss := status.Stat.Data()
	pis := getPluginInfo()
	vis := version.GetVersionInfo()
	bis := version.GetBuildInfo()
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data": map[string]interface{}{
			"build":     bis,
			"kotori-ng": vis,
			"plugins":   pis,
			"http":      hss,
		},
	}
	httputils.Respond(w, res, http.StatusOK, req)
	return
}

func getPluginInfo() (pis []PluginInfo) {
	for _, p := range pluginmanager.PM.Plugins {
		pis = append(pis, p.Info)
	}
	return
}
