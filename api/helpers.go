package api

// See: http://labs.omniti.com/labs/jsend

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
)

type M map[string]interface{}

func renderSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   data,
	})
}

func renderFailure(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "fail",
		"data":   data,
	})
}

func renderError(w http.ResponseWriter, msg string, data interface{}) {
	log.WithFields(logrus.Fields{
		"message": msg,
		"data":    data,
	}).Error("Error while handling request")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(500)

	resp := map[string]interface{}{
		"status":  "error",
		"message": msg,
	}
	if data != nil {
		resp["data"] = data
	}

	json.NewEncoder(w).Encode(resp)
}
