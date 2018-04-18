package handler

import (
	"net/http"
	"github.com/cool2645/kotori-ng/status"
	"github.com/cool2645/kotori-ng/httputils"
)

func GetStatus(w http.ResponseWriter, req *http.Request) {
	hss := status.Stat.Data()
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   hss,
	}
	httputils.Respond(w, res, http.StatusOK, req)
	return
}
