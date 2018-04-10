package handler

import (
	"net/http"
	"encoding/json"
	"github.com/cool2645/kotori-ng/model"
	"github.com/yanzay/log"
	. "github.com/cool2645/kotori-ng/httputils"
	"github.com/cool2645/kotori-ng/auth"
	"github.com/cool2645/kotori-ng/database"
	"strconv"
	"os"
	"github.com/cool2645/kotori-ng/config"
)

func Register(w http.ResponseWriter, req *http.Request) {
	// Parse Request
	var user model.User
	switch req.Header.Get("Content-Type") {
	case "application/json":
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
	// Validate Request
	if user.Username == "" {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Username cannot be blank",
		}
		Respond(w, res, http.StatusBadRequest, req)
		return
	}
	if user.Password == "" {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Password cannot be blank",
		}
		Respond(w, res, http.StatusBadRequest, req)
		return
	}
	// Process
	var err error
	if _, e := os.Stat(config.GlobCfg.LOCK_FILE); e == nil {
		err = model.MakeUser(&user)
	} else {
		os.OpenFile(config.GlobCfg.LOCK_FILE, os.O_RDONLY|os.O_CREATE, 0666)
		err = model.MakeAdmin(&user)
	}
	user.Password = ""
	if err != nil {
		switch err.Error() {
		case "MakeUser: StoreUser: UNIQUE constraint failed: users.username":
			res := map[string]interface{}{
				"code":   http.StatusOK,
				"result": false,
				"msg":    "This username has already been used",
			}
			Respond(w, res, http.StatusOK, req)
		case "MakeAdmin: StoreUser: UNIQUE constraint failed: users.username":
			res := map[string]interface{}{
				"code":   http.StatusOK,
				"result": false,
				"msg":    "This username has already been used",
			}
			Respond(w, res, http.StatusOK, req)
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
		"code":   http.StatusCreated,
		"result": true,
		"data":   user,
	}
	Respond(w, res, http.StatusCreated, req)
}

func GetUserByUUID(w http.ResponseWriter, req *http.Request) {
	// Parse Request
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
	// Check Privilege
	ok, user, msg := auth.CheckAuthorization(req)
	if !ok {
		res := map[string]interface{}{
			"code":   http.StatusUnauthorized,
			"result": false,
			"msg":    msg,
		}
		Respond(w, res, http.StatusUnauthorized, req)
		return
	}
	if !user.IsAdmin && user.UUID != uuid {
		res := map[string]interface{}{
			"code":   http.StatusUnauthorized,
			"result": false,
			"msg":    "You have no privilege to do so",
		}
		Respond(w, res, http.StatusUnauthorized, req)
		return
	}
	if user.UUID != uuid {
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
	}
	// Process
	user.Password = ""
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   user,
	}
	Respond(w, res, http.StatusOK, req)
	return
}

func GetUserByUsername(w http.ResponseWriter, req *http.Request) {
	// Parse Request
	username, err := GetUrlParameter(req, "username")
	if err != nil {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    err.Error(),
		}
		Respond(w, res, http.StatusBadRequest, req)
		return
	}
	// Check Privilege
	ok, user, msg := auth.CheckAuthorization(req)
	if !ok {
		res := map[string]interface{}{
			"code":   http.StatusUnauthorized,
			"result": false,
			"msg":    msg,
		}
		Respond(w, res, http.StatusUnauthorized, req)
		return
	}
	if !user.IsAdmin && user.Username != username {
		res := map[string]interface{}{
			"code":   http.StatusUnauthorized,
			"result": false,
			"msg":    "You have no privilege to do so",
		}
		Respond(w, res, http.StatusUnauthorized, req)
		return
	}
	if user.Username != username {
		user, err = model.GetUserByUsername(database.DB, username)
		if err != nil {
			if err.Error() == "GetUserByUsername: record not found" {
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
	}
	// Process
	user.Password = ""
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   user,
	}
	Respond(w, res, http.StatusOK, req)
	return
}

func PromoteAdmin(w http.ResponseWriter, req *http.Request) {
	// Parse Request
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
	// Check Privilege
	ok, user, msg := auth.CheckAuthorization(req)
	if !ok {
		res := map[string]interface{}{
			"code":   http.StatusUnauthorized,
			"result": false,
			"msg":    msg,
		}
		Respond(w, res, http.StatusUnauthorized, req)
		return
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
	if user.UUID != uuid {
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
	}
	// Process
	err = model.PromoteUserToAdmin(database.DB, &user)
	if err != nil {
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred updating user privilege",
		}
		Respond(w, res, http.StatusInternalServerError, req)
		return
	}
	user.Password = ""
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   user,
	}
	Respond(w, res, http.StatusOK, req)
	return
}

func ListUsers(w http.ResponseWriter, req *http.Request)  {
	// Check Privilege
	ok, user, msg := auth.CheckAuthorization(req)
	if !ok {
		res := map[string]interface{}{
			"code":   http.StatusUnauthorized,
			"result": false,
			"msg":    msg,
		}
		Respond(w, res, http.StatusUnauthorized, req)
		return
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
	// Parse Request
	form := make(map[string]string)
	switch req.Header.Get("Content-Type") {
	case "application/json":
		err := json.NewDecoder(req.Body).Decode(&form)
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
	default:
		req.ParseForm()
		for k, v := range req.Form {
			if len(v) == 1 {
				form[k] = v[0]
			}
		}
	}
	var page, perPage uint = 1, 15
	var orderBy, order = "id", "desc"
	if val, ok := form["page"]; ok {
		page64, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			page = uint(page64)
		}
	}
	if val, ok := form["per_page"]; ok {
		perPage64, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			perPage = uint(perPage64)
		}
		if perPage == 0 {
			perPage = 15
		}
	}
	if val, ok := form["order_by"]; ok {
		orderBy = val
	}
	if val, ok := form["order"]; ok {
		order = val
	}
	// Process
	users, total, err := model.ListUsers(database.DB, page, perPage, orderBy, order)
	if err != nil {
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred querying users",
		}
		Respond(w, res, http.StatusInternalServerError, req)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data": map[string]interface{}{
			"current_page": page,
			"per_page": perPage,
			"total_page": total / perPage + 1,
			"total": total,
			"data":   users,
		},
	}
	Respond(w, res, http.StatusOK, req)
	return
}

func UpdateUser(w http.ResponseWriter, req *http.Request) {
	// Parse Request
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
	// Check Privilege
	ok, user, msg := auth.CheckAuthorization(req)
	if !ok {
		res := map[string]interface{}{
			"code":   http.StatusUnauthorized,
			"result": false,
			"msg":    msg,
		}
		Respond(w, res, http.StatusUnauthorized, req)
		return
	}
	if !user.IsAdmin && user.UUID != uuid {
		res := map[string]interface{}{
			"code":   http.StatusUnauthorized,
			"result": false,
			"msg":    "You have no privilege to do so",
		}
		Respond(w, res, http.StatusUnauthorized, req)
		return
	}
	if user.UUID != uuid {
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
	}
	// Parse Request
	patch := make(map[string]string)
	switch req.Header.Get("Content-Type") {
	case "application/json":
		err := json.NewDecoder(req.Body).Decode(&patch)
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
	default:
		req.ParseForm()
		for k, v := range req.Form {
			if len(v) == 1 {
				patch[k] = v[0]
			}
		}
	}
	// Process
	err = model.UpdateUser(database.DB, &user, patch)
	if err != nil {
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred updating user",
		}
		Respond(w, res, http.StatusInternalServerError, req)
		return
	}
	user.Password = ""
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   user,
	}
	Respond(w, res, http.StatusOK, req)
	return
}

func UpdateUserSetUsername(w http.ResponseWriter, req *http.Request) {
	// Parse Request
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
	// Check Privilege
	ok, user, msg := auth.CheckAuthorization(req)
	if !ok {
		res := map[string]interface{}{
			"code":   http.StatusUnauthorized,
			"result": false,
			"msg":    msg,
		}
		Respond(w, res, http.StatusUnauthorized, req)
		return
	}
	if !user.IsAdmin && user.UUID != uuid {
		res := map[string]interface{}{
			"code":   http.StatusUnauthorized,
			"result": false,
			"msg":    "You have no privilege to do so",
		}
		Respond(w, res, http.StatusUnauthorized, req)
		return
	}
	if user.UUID != uuid {
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
	}
	// Parse Request
	patch := make(map[string]string)
	switch req.Header.Get("Content-Type") {
	case "application/json":
		err := json.NewDecoder(req.Body).Decode(&patch)
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
	default:
		req.ParseForm()
		for k, v := range req.Form {
			if len(v) == 1 {
				patch[k] = v[0]
			}
		}
	}
	// Validate Request
	if val, ok := patch["username"]; !ok || val == "" {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid username",
		}
		Respond(w, res, http.StatusBadRequest, req)
		return
	}
	// Process
	err = model.UpdateUserSetUsername(database.DB, &user, patch["username"])
	if err != nil {
		if err.Error() == "UpdateUserSetUsername: UNIQUE constraint failed: users.username" {
			res := map[string]interface{}{
				"code":   http.StatusOK,
				"result": false,
				"msg":    "This username has already been used",
			}
			Respond(w, res, http.StatusOK, req)
			return
		}
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred updating username",
		}
		Respond(w, res, http.StatusInternalServerError, req)
		return
	}
	user.Password = ""
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   user,
	}
	Respond(w, res, http.StatusOK, req)
	return
}

func UpdateUserSetPassword(w http.ResponseWriter, req *http.Request) {
	// Parse Request
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
	// Check Privilege
	ok, user, msg := auth.CheckAuthorization(req)
	if !ok {
		res := map[string]interface{}{
			"code":   http.StatusUnauthorized,
			"result": false,
			"msg":    msg,
		}
		Respond(w, res, http.StatusUnauthorized, req)
		return
	}
	if !user.IsAdmin && user.UUID != uuid {
		res := map[string]interface{}{
			"code":   http.StatusUnauthorized,
			"result": false,
			"msg":    "You have no privilege to do so",
		}
		Respond(w, res, http.StatusUnauthorized, req)
		return
	}
	if user.UUID != uuid {
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
	}
	// Parse Request
	patch := make(map[string]string)
	switch req.Header.Get("Content-Type") {
	case "application/json":
		err := json.NewDecoder(req.Body).Decode(&patch)
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
	default:
		req.ParseForm()
		for k, v := range req.Form {
			if len(v) == 1 {
				patch[k] = v[0]
			}
		}
	}
	// Validate Request
	if val, ok := patch["password"]; !ok || val == "" {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Password cannot be blank",
		}
		Respond(w, res, http.StatusBadRequest, req)
		return
	}
	// Process
	err = model.UpdateUserSetPassword(database.DB, &user, patch["password"])
	if err != nil {
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred updating password",
		}
		Respond(w, res, http.StatusInternalServerError, req)
		return
	}
	user.Password = ""
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   user,
	}
	Respond(w, res, http.StatusOK, req)
	return
}

func DismissAdmin(w http.ResponseWriter, req *http.Request) {
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
		return
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
	err = model.DismissUserFromAdmin(database.DB, &user)
	if err != nil {
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred updating user privilege",
		}
		Respond(w, res, http.StatusInternalServerError, req)
		return
	}
	user.Password = ""
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   user,
	}
	Respond(w, res, http.StatusOK, req)
	return
}