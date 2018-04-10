package handler

import (
	"net/http"
	"encoding/json"
	"github.com/cool2645/kotori-ng/model"
	"github.com/yanzay/log"
	. "github.com/cool2645/kotori-ng/httputils"
	"github.com/cool2645/kotori-ng/auth"
)

func Login(w http.ResponseWriter, req *http.Request) {
	// Parse Request
	var username, password string
	switch req.Header.Get("Content-Type") {
	case "application/json":
		var user model.User
		err := json.NewDecoder(req.Body).Decode(&user)
		if err != nil {
			log.Error(err)
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Error occurred parsing json request",
			}
			Respond(w, res, http.StatusBadRequest, req)
			return
		}
		username = user.Username
		password = user.Password
	default:
		req.ParseForm()
		if len(req.Form["username"]) != 1 {
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Invalid username",
			}
			Respond(w, res, http.StatusBadRequest, req)
			return
		}
		username = req.Form["username"][0]
		if len(req.Form["password"]) != 1 {
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Invalid password",
			}
			Respond(w, res, http.StatusBadRequest, req)
			return
		}
		password = req.Form["password"][0]
	}
	// Process
	if ok, token, msg := auth.GenerateToken(username, password); ok {
		res := map[string]interface{}{
			"code":   http.StatusOK,
			"result": true,
			"data":   token,
		}
		Respond(w, res, http.StatusOK, req)
	} else {
		res := map[string]interface{}{
			"code":   http.StatusOK,
			"result": false,
			"msg":    msg,
		}
		Respond(w, res, http.StatusOK, req)
	}
}

func GetMe(w http.ResponseWriter, req *http.Request) {
	if ok, user, msg := auth.CheckAuthorization(req); ok {
		user.Password = ""
		res := map[string]interface{}{
			"code":   http.StatusOK,
			"result": true,
			"data":   user,
		}
		Respond(w, res, http.StatusOK, req)
	} else {
		res := map[string]interface{}{
			"code":   http.StatusUnauthorized,
			"result": false,
			"msg":    msg,
		}
		Respond(w, res, http.StatusUnauthorized, req)
	}
}
