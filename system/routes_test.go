package service

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"
	"net/http/httputil"

	"github.com/joho/godotenv"
	"github.com/namsral/flag"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/rbac"
)

func TestUsers(t *testing.T) {
	godotenv.Load("../.env")

	mountFlags("system", Flags, auth.Flags, rbac.Flags)

	// log to stdout not stderr
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := Init()
	assert(t, err == nil, "Error initializing: %+v", err)

	req, err := http.NewRequest("GET", "http://localhost/auth/check", nil)
	assert(t, err == nil, "Error creating request: %+v", err)

	recorder := httptest.NewRecorder()
	Routes().ServeHTTP(recorder, req)

	resp := recorder.Result()

	fmt.Println(">>> (request)")
	fmt.Println(request(req))
	fmt.Println("----")
	fmt.Println("<<< (response)")
	fmt.Println(response(resp))
}

func request(req *http.Request) string {
	b, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return ">>> Error: " + err.Error()
	}
	if b != nil {
		return strings.TrimSpace(string(b))
	}
	return ""
}

func response(resp *http.Response) string {
	b, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return "<<< Error: " + err.Error()
	}
	if b != nil {
		return strings.TrimSpace(string(b))
	}
	return ""
}

func assert(t *testing.T, ok bool, format string, args ...interface{}) bool {
	if !ok {
		t.Fatalf(format, args...)
	}
	return ok
}

func mountFlags(prefix string, mountFlags ...func(...string)) {
	for _, mount := range mountFlags {
		mount(prefix)
	}
	flag.Parse()
}
