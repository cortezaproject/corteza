package scim

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

func send(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Error("could not encode payload", zap.Error(err))
	}
}

func sendError(w http.ResponseWriter, err error) {
	var (
		status = http.StatusInternalServerError
	)

	if er, ok := err.(*errorResponse); ok {
		status = er.Status
	}

	send(w, status, err)
}
