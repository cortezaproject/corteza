package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"io/ioutil"
	"net/http/httptest"
)

func TestTests(t *testing.T) {
	testResponse := func(output interface{}) string {
		w := httptest.NewRecorder()
		r := &http.Request{Header: http.Header{}}
		r.Header.Add("accept", "application/json")
		Send(w, r, output)
		body, _ := ioutil.ReadAll(w.Result().Body)
		return string(body)
	}

	var cc = []struct {
		name string
		inp  interface{}
		out  string
	}{
		{"nil", nil, `{"response":false}`},
		{"bool true", true, `{"response":true}`},
		{"bool false", false, `{"response":false}`},
		{"string empty", "", `{"response":false}`},
		{"string", "string", `{"response":"string"}`},
		{"int zero", 0, `{"response":0}`},
		{"int non-zero", 1337, `{"response":1337}`},
		{"int sub-zero", -1, `{"response":-1}`},
		{"error nil", func() error { return nil }, `{"response":false}`},
		{"error", func() error { return fmt.Errorf("error response") }, `{"error":{"message":"error response"}}`},
		{"value + error", func() (interface{}, error) { return "string response", fmt.Errorf("error response") }, `{"error":{"message":"error response"}}`},
		{"empty value + error", func() (interface{}, error) { return "", fmt.Errorf("error response") }, `{"error":{"message":"error response"}}`},
		{"value + empty error", func() (interface{}, error) { return "string response", nil }, `{"response":"string response"}`},
		{"success default", Success(), `{"success":{"message":"OK"}}`},
		{"ok", OK(), `{"success":{"message":"OK"}}`},
		{"success custom", Success("string"), `{"success":{"message":"string"}}`},
		{"error stdlib", fmt.Errorf("string"), `{"error":{"message":"string"}}`},
		{"error stdlib nil", func() interface{} { return func() error { return nil }() }(), `{"response":false}`},
		{"func json nil", func() ([]byte, error) { return json.Marshal(nil) }, `null`},
		{"func json false", func() ([]byte, error) { return json.Marshal(false) }, `false`},
		{"func json 0", func() ([]byte, error) { return json.Marshal(0) }, `0`},
		{"func json empty string", func() ([]byte, error) { return json.Marshal("") }, `""`},
		{"func writer/req", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("foo")) }, `foo`},
		{"custom struct", struct {
			Name string `json:"name"`
		}{"Corteza"}, `{"response":{"name":"Corteza"}}`},
	}

	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			got := strings.TrimSpace(testResponse(c.inp))
			if got != strings.TrimSpace(c.out) {
				t.Errorf("got %#v, expected %#v", got, c.out)
			}
		})
	}

}
