package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"net/http"
	"net/http/httptest"
	"net/http/httputil"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/rbac"
	systemRepository "github.com/crusttech/crust/system/repository"
	systemTypes "github.com/crusttech/crust/system/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/namsral/flag"
)

type (
	jsonResponse struct {
		Error struct {
			Message string `json:"message"`
			Trace   string `json:"trace,omitempty"`
		} `json:"error"`
	}
)

func TestUsers(t *testing.T) {
	godotenv.Load("../.env")

	mountFlags("system", Flags, auth.Flags, rbac.Flags)

	// log to stdout not stderr
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Initialize routes and exit on failure.
	err := Init()
	assert(t, err == nil, "Error initializing: %+v", err)

	// Send check request with invalid JWT token.
	{
		req, err := http.NewRequest("GET", "http://127.0.0.1/auth/check", nil)
		assert(t, err == nil, "Error creating request: %+v", err)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":  "zblj",
			"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})
		tokenString, err := token.SignedString([]byte("secret"))
		assert(t, err == nil, "Error creating JWT token: %+v", err)

		req.AddCookie(&http.Cookie{
			Name:     "jwt",
			Value:    tokenString,
			Domain:   ".localhost",
			Expires:  time.Now().Add(time.Hour),
			HttpOnly: true,
			MaxAge:   50000,
			Path:     "/auth",
		})

		recorder := httptest.NewRecorder()
		Routes().ServeHTTP(recorder, req)
		resp := recorder.Result()

		fmt.Println(">>> (request)")
		fmt.Println(request(req))
		fmt.Println("----")
		fmt.Println("<<< (response)")
		fmt.Println(response(resp))

		jr, err := decodeJson(resp.Body)
		assert(t, err == nil, "Error decoding response body: %+v", err)
		assert(t, jr.Error.Message == "failed to authorize request: signature is invalid", "Expected error 'failed to authorize request: signature is invalid' got: %+v", jr.Error.Message)
	}

	// Send "Login" request without parameters.
	{
		req, err := http.NewRequest("POST", "http://localhost/auth/login", nil)
		assert(t, err == nil, "Error creating request: %+v", err)

		recorder := httptest.NewRecorder()
		Routes().ServeHTTP(recorder, req)
		resp := recorder.Result()

		fmt.Println(">>> (request)")
		fmt.Println(request(req))
		fmt.Println("----")
		fmt.Println("<<< (response)")
		fmt.Println(response(resp))

		jr, err := decodeJson(resp.Body)
		assert(t, err == nil, "Error decoding response body: %+v", err)
		assert(t, jr.Error.Message == "missing form body", "Expected error 'missing form body' got: %+v", jr.Error.Message)
	}

	// Send "Login" request with missing user.
	{
		jsonStr := `{"username":"test123","password":"test123"}`

		req, err := http.NewRequest("POST", "http://localhost/auth/login", strings.NewReader(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		assert(t, err == nil, "Error creating request: %+v", err)

		recorder := httptest.NewRecorder()
		Routes().ServeHTTP(recorder, req)
		resp := recorder.Result()

		fmt.Println(">>> (request)")
		fmt.Println(request(req))
		fmt.Println("----")
		fmt.Println("<<< (response)")
		fmt.Println(response(resp))

		jr, err := decodeJson(resp.Body)
		assert(t, err == nil, "Error decoding response body: %+v", err)
		assert(t, jr.Error.Message == "crust.auth.repository.UserNotFound", "Expected error 'crust.auth.repository.UserNotFound' got: %+v", jr.Error.Message)
	}

	// Create user.
	user := &systemTypes.User{
		ID:       1337,
		Username: "johndoe",
	}
	{
		err := user.GeneratePassword("johndoe123")
		assert(t, err == nil, "Error generating password: %+v", err)
	}
	{
		userAPI := systemRepository.User(context.Background(), nil)
		_, err := userAPI.Create(user)
		assert(t, err == nil, "Error when inserting user: %+v", err)
	}

	// Send "Login" request with existing user.
	{
		form := url.Values{}
		form.Add("username", "johndoe")
		form.Add("password", "johndoe123")

		req, err := http.NewRequest("POST", "http://localhost/auth/login", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		assert(t, err == nil, "Error creating request: %+v", err)

		recorder := httptest.NewRecorder()
		Routes().ServeHTTP(recorder, req)
		resp := recorder.Result()

		fmt.Println(">>> (request)")
		fmt.Println(request(req))
		fmt.Println("----")
		fmt.Println("<<< (response)")
		fmt.Println(response(resp))

		c := resp.Cookies()
		assert(t, len(c) == 1, "Expected 1 cookie value, got: %+v", len(c))
		assert(t, c[0].Value != "", "Expected non empty jwt token, got: %+v", c[0].Value)

		type jsonResponse struct {
			Response struct {
				UserID   string `json:"userID"`
				Username string `json:"username"`
			} `json:"response"`
		}

		var jr jsonResponse
		err = json.NewDecoder(resp.Body).Decode(&jr)
		assert(t, err == nil, "Error decoding response body: %+v", err)
		assert(t, jr.Response.UserID != "0", "Expected userID not to be 0, got: %+v", jr.Response.UserID)
		assert(t, jr.Response.Username == "johndoe", "Expected username 'johndoe', got: %+v", jr.Response.Username)

		// Check JWT token after successful login.
		req, err = http.NewRequest("GET", "http://localhost/auth/check", nil)
		assert(t, err == nil, "Error creating request: %+v", err)

		Routes().ServeHTTP(recorder, req)

		resp = recorder.Result()

		fmt.Println(">>> (request)")
		fmt.Println(request(req))
		fmt.Println("----")
		fmt.Println("<<< (response)")
		fmt.Println(response(resp))

		assert(t, resp.StatusCode == 200, "Expected http status code 200, got: %+v", resp.StatusCode)
		assert(t, len(c) == 1, "Expected 1 cookie value, got: %+v", len(c))
		assert(t, c[0].Value != "", "Expected non empty jwt token, got: %+v", c[0].Value)
	}

	// Send "Login" request with existing user.
	{
		jsonStr := `{"username": "johndoe", "password": "johndoe123"}`

		req, err := http.NewRequest("POST", "http://localhost/auth/login", strings.NewReader(jsonStr))
		req.Header.Add("Content-Type", "application/json")

		assert(t, err == nil, "Error creating request: %+v", err)

		recorder := httptest.NewRecorder()
		Routes().ServeHTTP(recorder, req)
		resp := recorder.Result()

		fmt.Println(">>> (request)")
		fmt.Println(request(req))
		fmt.Println("----")
		fmt.Println("<<< (response)")
		fmt.Println(response(resp))

		c := resp.Cookies()
		assert(t, len(c) == 1, "Expected 1 cookie value, got: %+v", len(c))
		assert(t, c[0].Value != "", "Expected non empty jwt token, got: %+v", c[0].Value)

		type jsonResponse struct {
			Response struct {
				UserID   string `json:"userID"`
				Username string `json:"username"`
			} `json:"response"`
		}

		var jr jsonResponse
		err = json.NewDecoder(resp.Body).Decode(&jr)
		assert(t, err == nil, "Error decoding response body: %+v", err)
		assert(t, jr.Response.UserID != "0", "Expected userID not to be 0, got: %+v", jr.Response.UserID)
		assert(t, jr.Response.Username == "johndoe", "Expected username 'johndoe', got: %+v", jr.Response.Username)

		// Check JWT token after successful login.
		req, err = http.NewRequest("GET", "http://localhost/auth/check", nil)
		assert(t, err == nil, "Error creating request: %+v", err)

		Routes().ServeHTTP(recorder, req)

		resp = recorder.Result()

		fmt.Println(">>> (request)")
		fmt.Println(request(req))
		fmt.Println("----")
		fmt.Println("<<< (response)")
		fmt.Println(response(resp))

		assert(t, resp.StatusCode == 200, "Expected http status code 200, got: %+v", resp.StatusCode)
		assert(t, len(c) == 1, "Expected 1 cookie value, got: %+v", len(c))
		assert(t, c[0].Value != "", "Expected non empty jwt token, got: %+v", c[0].Value)
	}

	// Send "Logout" request and expect empty jwt token.
	{
		req, err := http.NewRequest("GET", "http://localhost/auth/logout", nil)
		assert(t, err == nil, "Error creating request: %+v", err)

		recorder := httptest.NewRecorder()
		Routes().ServeHTTP(recorder, req)
		resp := recorder.Result()

		fmt.Println(">>> (request)")
		fmt.Println(request(req))
		fmt.Println("----")
		fmt.Println("<<< (response)")
		fmt.Println(response(resp))

		c := resp.Cookies()

		assert(t, resp.StatusCode == 200, "Expected http status code 200, got: %+v", resp.StatusCode)
		assert(t, len(c) == 1, "Expected 1 cookie value, got: %+v", len(c))
		assert(t, c[0].Value == "", "Expected empty jwt token, got: %+v", c[0].Value)
	}

	// Send check request without JWT token.
	{
		req, err := http.NewRequest("GET", "http://127.0.0.1/auth/check", nil)
		assert(t, err == nil, "Error creating request: %+v", err)

		recorder := httptest.NewRecorder()
		Routes().ServeHTTP(recorder, req)

		resp := recorder.Result()

		fmt.Println(">>> (request)")
		fmt.Println(request(req))
		fmt.Println("----")
		fmt.Println("<<< (response)")
		fmt.Println(response(resp))

		jr, err := decodeJson(resp.Body)
		assert(t, err == nil, "Error decoding response body: %+v", err)
		assert(t, jr.Error.Message == "http: named cookie not present", "Expected error 'http: named cookie not present' got: %+v", jr.Error.Message)
	}
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

func decodeJson(r io.Reader) (jsonResponse, error) {
	var ret jsonResponse
	err := json.NewDecoder(r).Decode(&ret)
	if err != nil {
		return jsonResponse{}, err
	}
	return ret, nil
}
