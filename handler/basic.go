package handler

import (
	"net/http"
	. "github.com/cool2645/kotori-ng/httputils"
)

func Pong(w http.ResponseWriter, req *http.Request) {
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"msg":    "OK",
	}
	ResponseJson(w, res, http.StatusOK)
	return
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	res := map[string]interface{}{
		"code":   http.StatusNotFound,
		"result": false,
		"msg":    "Not found",
	}
	ResponseJson(w, res, http.StatusNotFound)
	return
}