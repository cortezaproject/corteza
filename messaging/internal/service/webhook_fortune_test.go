// +build integration

package service

import (
	"fmt"

	"encoding/json"
	"net/http"
)

type Fortune struct{}

func (*Fortune) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fortune := "Fortune favors the prepared mind. - Louis Pasteur"
	username := r.URL.Query()["username"]
	if len(username) > 0 {
		response := struct {
			Username string `json:"username"`
			Text     string `json:"text"`
		}{
			username[0],
			fortune,
		}
		b, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}
	fmt.Fprintf(w, fortune)
}
