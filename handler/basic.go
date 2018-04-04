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
	Respond(w, res, http.StatusOK, req)
	return
}

func NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	res := map[string]interface{}{
		"code":   http.StatusNotFound,
		"result": false,
		"msg":    "Not found",
	}
	Respond(w, res, http.StatusNotFound, req)
	return
}