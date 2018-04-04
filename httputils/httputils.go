package httputils

import (
	"net/http"
	"encoding/json"
	"github.com/yanzay/log"
	"strings"
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

func Respond(w http.ResponseWriter, data map[string]interface{}, httpStatusCode int, r *http.Request) {
	acc := r.Header.Get("Accept")
	if strings.Contains(acc, "application/json") {
		RespondJson(w, data, httpStatusCode)
	} else {
		RespondJson(w, data, httpStatusCode)
	}
	return
}
