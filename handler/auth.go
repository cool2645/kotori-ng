package handler

import (
	"net/http"
	"encoding/json"
	"github.com/cool2645/kotori-ng/model"
	"github.com/yanzay/log"
	. "github.com/cool2645/kotori-ng/httputils"
	"github.com/cool2645/kotori-ng/auth"
	"strconv"
	"github.com/cool2645/kotori-ng/database"
)

func Login(w http.ResponseWriter, req *http.Request) {
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

func Register(w http.ResponseWriter, req *http.Request) {
	var user model.User
	switch req.Header.Get("Content-Type") {
	case "application/json":
		err := json.NewDecoder(req.Body).Decode(&user)
		if err != nil || user.Username == "" || user.Password == "" {
			log.Error(err)
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Error occurred parsing json request",
			}
			Respond(w, res, http.StatusBadRequest, req)
			return
		}
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
		user.Username = req.Form["username"][0]
		if len(req.Form["password"]) != 1 {
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Invalid password",
			}
			Respond(w, res, http.StatusBadRequest, req)
			return
		}
		user.Password = req.Form["password"][0]
		if len(req.Form["name"]) == 1 {
			user.Name = req.Form["name"][0]
		}
	}
	err := auth.MakeUser(&user)
	user.Password = ""
	if err != nil {
		switch err.Error() {
		case "MakeUser: StoreUser: UNIQUE constraint failed: users.username":
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "This username has already been used",
			}
			Respond(w, res, http.StatusBadRequest, req)
		default:
			res := map[string]interface{}{
				"code":   http.StatusInternalServerError,
				"result": false,
				"msg":    err.Error(),
			}
			Respond(w, res, http.StatusInternalServerError, req)
		}
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   user,
	}
	Respond(w, res, http.StatusOK, req)
}

func ChangePrivilege(w http.ResponseWriter, req *http.Request) {
	uuid, err := GetUrlParameter(req, "uuid")
	if err != nil {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    err.Error(),
		}
		Respond(w, res, http.StatusBadRequest, req)
		return
	}
	ok, user, msg := auth.CheckAuthorization(req)
	if !ok {
		res := map[string]interface{}{
			"code":   http.StatusUnauthorized,
			"result": false,
			"msg":    msg,
		}
		Respond(w, res, http.StatusUnauthorized, req)
	}
	if !user.IsAdmin {
		res := map[string]interface{}{
			"code":   http.StatusUnauthorized,
			"result": false,
			"msg":    "You have no privilege to do so",
		}
		Respond(w, res, http.StatusUnauthorized, req)
		return
	}
	var isAdmin bool
	switch req.Header.Get("Content-Type") {
	case "application/json":
		var m map[string]interface{}
		err := json.NewDecoder(req.Body).Decode(&m)
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
		isAdmin, err = strconv.ParseBool(m["is_admin"].(string))
		if err != nil {
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Error occurred parsing privilege",
			}
			Respond(w, res, http.StatusBadRequest, req)
			return
		}
	default:
		req.ParseForm()
		if len(req.Form["is_admin"]) != 1 {
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Invalid privilege",
			}
			Respond(w, res, http.StatusBadRequest, req)
			return
		}
		isAdmin, err = strconv.ParseBool(req.Form["is_admin"][0])
		if err != nil {
			res := map[string]interface{}{
				"code":   http.StatusBadRequest,
				"result": false,
				"msg":    "Error occurred parsing privilege",
			}
			Respond(w, res, http.StatusBadRequest, req)
			return
		}
	}
	user, err = model.GetUserByUUID(database.DB, uuid)
	if err != nil {
		if err.Error() == "GetUserByUUID: record not found" {
			res := map[string]interface{}{
				"code":   http.StatusNotFound,
				"result": false,
				"msg":    "User not found",
			}
			Respond(w, res, http.StatusNotFound, req)
			return
		} else {
			res := map[string]interface{}{
				"code":   http.StatusInternalServerError,
				"result": false,
				"msg":    "Error occurred querying user",
			}
			Respond(w, res, http.StatusInternalServerError, req)
			return
		}
	}
	if isAdmin {
		err = model.PromoteUserToAdmin(database.DB, &user)
	} else {
		err = model.DismissUserFromAdmin(database.DB, &user)
	}
	if err != nil {
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred updating user privilege",
		}
		Respond(w, res, http.StatusInternalServerError, req)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   user,
	}
	Respond(w, res, http.StatusOK, req)
	return
}
