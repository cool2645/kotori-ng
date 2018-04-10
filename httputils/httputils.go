package httputils

import (
	"net/http"
	"encoding/json"
	"github.com/yanzay/log"
	"strings"
	"errors"
	"strconv"
	"regexp"
	"fmt"
	"github.com/gorilla/mux"
)

func RespondJson(w http.ResponseWriter, data map[string]interface{}, httpStatusCode int) {
	resJson, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred encoding response.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	w.Write(resJson)
	return
}

func RespondJsonp(w http.ResponseWriter, data map[string]interface{}, callback string) {
	resJson, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred encoding response.", http.StatusInternalServerError)
		return
	}
	resJson = []byte(fmt.Sprintf("%s(%s)", callback, resJson))
	w.Header().Set("Content-Type", "application/javascript")
	w.Write(resJson)
	return
}

func RespondFormattedText(w http.ResponseWriter, data map[string]interface{}, httpStatusCode int) {
	res, err := json.MarshalIndent(data, "", "\t")
	r, _ := regexp.Compile(`"(\w+)"\s*:`)
	res = r.ReplaceAll(res, []byte("$1:"))
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred encoding response.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(httpStatusCode)
	w.Write(res)
	return
}

func RespondPlainText(w http.ResponseWriter, data string, httpStatusCode int) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(httpStatusCode)
	w.Write([]byte(data))
	return
}

func Respond(w http.ResponseWriter, data map[string]interface{}, httpStatusCode int, r *http.Request) {
	acc := r.Header.Get("Accept")
	if strings.Contains(acc, "application/json") {
		RespondJson(w, data, httpStatusCode)
	} else if strings.Contains(acc, "text/plain") || strings.Contains(acc, "text/html") {
		RespondFormattedText(w, data , httpStatusCode)
	} else {
		RespondPlainText(w, "406 Not Acceptable" , http.StatusNotAcceptable)
	}
	return
}

func GetInt64UrlParameter(r *http.Request, key string) (value int64, err error) {
	val, err := GetUrlParameter(r, key)
	if err != nil {
		return
	}
	value, err = strconv.ParseInt(val, 10, 32)
	if err != nil {
		return
	}
	return
}

func GetUint64UrlParameter(r *http.Request, key string) (value uint64, err error) {
	val, err := GetUrlParameter(r, key)
	if err != nil {
		return
	}
	value, err = strconv.ParseUint(val, 10, 32)
	if err != nil {
		return
	}
	return
}

func GetUintUrlParameter(r *http.Request, key string) (value uint, err error) {
	val, err := GetUint64UrlParameter(r, key)
	if err != nil {
		return
	}
	value = uint(val)
	return
}

func GetIntUrlParameter(r *http.Request, key string) (value int, err error) {
	val, err := GetUint64UrlParameter(r, key)
	if err != nil {
		return
	}
	value = int(val)
	return
}

func GetUrlParameter(r *http.Request, key string) (value string, err error) {
	query := mux.Vars(r)
	value, ok := query[key]
	if !ok || value == "" {
		err = errors.New("cannot retrieve requested key")
		return
	}
	return
}