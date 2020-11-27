package scim

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

func send(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Error("could not encode payload", zap.Error(err))
	}
}
