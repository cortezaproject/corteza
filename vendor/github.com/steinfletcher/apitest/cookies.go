package apitest

import (
	"fmt"
	"net/http"
	"time"
)

// Cookie used to represent an http cookie
type Cookie struct {
	name     *string
	value    *string
	path     *string
	domain   *string
	expires  *time.Time
	maxAge   *int
	secure   *bool
	httpOnly *bool
}

// NewCookie creates a new Cookie with the provided name
func NewCookie(name string) *Cookie {
	return &Cookie{
		name: &name,
	}
}

// Value sets the value of the Cookie
func (cookie *Cookie) Value(value string) *Cookie {
	cookie.value = &value
	return cookie
}

// Path sets the path of the Cookie
func (cookie *Cookie) Path(path string) *Cookie {
	cookie.path = &path
	return cookie
}

// Domain sets the domain of the Cookie
func (cookie *Cookie) Domain(domain string) *Cookie {
	cookie.domain = &domain
	return cookie
}

// Expires sets the expires time of the Cookie
func (cookie *Cookie) Expires(expires time.Time) *Cookie {
	cookie.expires = &expires
	return cookie
}

// MaxAge sets the maxage of the Cookie
func (cookie *Cookie) MaxAge(maxAge int) *Cookie {
	cookie.maxAge = &maxAge
	return cookie
}

// Secure sets the secure bool of the Cookie
func (cookie *Cookie) Secure(secure bool) *Cookie {
	cookie.secure = &secure
	return cookie
}

// HttpOnly sets the httpOnly bool of the Cookie
func (cookie *Cookie) HttpOnly(httpOnly bool) *Cookie {
	cookie.httpOnly = &httpOnly
	return cookie
}

// ToHttpCookie transforms the Cookie to an http cookie
func (cookie *Cookie) ToHttpCookie() *http.Cookie {
	httpCookie := http.Cookie{}

	if cookie.name != nil {
		httpCookie.Name = *cookie.name
	}

	if cookie.value != nil {
		httpCookie.Value = *cookie.value
	}

	if cookie.path != nil {
		httpCookie.Path = *cookie.path
	}

	if cookie.domain != nil {
		httpCookie.Domain = *cookie.domain
	}

	if cookie.expires != nil {
		httpCookie.Expires = *cookie.expires
	}

	if cookie.maxAge != nil {
		httpCookie.MaxAge = *cookie.maxAge
	}

	if cookie.secure != nil {
		httpCookie.Secure = *cookie.secure
	}

	if cookie.httpOnly != nil {
		httpCookie.HttpOnly = *cookie.httpOnly
	}

	return &httpCookie
}

// FromHTTPCookie transforms an http cookie into a Cookie
func FromHTTPCookie(httpCookie *http.Cookie) *Cookie {
	return NewCookie(httpCookie.Name).
		Value(httpCookie.Value).
		Path(httpCookie.Path).
		Domain(httpCookie.Domain).
		Expires(httpCookie.Expires).
		MaxAge(httpCookie.MaxAge).
		Secure(httpCookie.Secure).
		HttpOnly(httpCookie.HttpOnly)
}

// Compares cookies based on only the provided fields from Cookie.
// Supported fields are Name, Value, Domain, Path, Expires, MaxAge, Secure and HttpOnly
func compareCookies(expectedCookie *Cookie, actualCookie *http.Cookie) (bool, []string) {
	cookieFound := *expectedCookie.name == actualCookie.Name
	compareErrors := make([]string, 0)
	if cookieFound {
		compareErrors = compareValue(expectedCookie, actualCookie, compareErrors)
		compareErrors = compareDomain(expectedCookie, actualCookie, compareErrors)
		compareErrors = comparePath(expectedCookie, actualCookie, compareErrors)
		compareErrors = compareExpires(expectedCookie, actualCookie, compareErrors)
		compareErrors = compareMaxAge(expectedCookie, actualCookie, compareErrors)
		compareErrors = compareSecure(expectedCookie, actualCookie, compareErrors)
		compareErrors = compareHttpOnly(expectedCookie, actualCookie, compareErrors)
	}

	return cookieFound, compareErrors
}

func compareHttpOnly(expectedCookie *Cookie, actualCookie *http.Cookie, compareErrors []string) []string {
	if expectedCookie.httpOnly != nil && *expectedCookie.httpOnly != actualCookie.HttpOnly {
		compareErrors = append(compareErrors, formatError("HttpOnly", *expectedCookie.httpOnly, actualCookie.HttpOnly))
	}
	return compareErrors
}

func compareSecure(expectedCookie *Cookie, actualCookie *http.Cookie, compareErrors []string) []string {
	if expectedCookie.secure != nil && *expectedCookie.secure != actualCookie.Secure {
		compareErrors = append(compareErrors, formatError("Secure", *expectedCookie.secure, actualCookie.Secure))
	}
	return compareErrors
}

func compareMaxAge(expectedCookie *Cookie, actualCookie *http.Cookie, compareErrors []string) []string {
	if expectedCookie.maxAge != nil && *expectedCookie.maxAge != actualCookie.MaxAge {
		compareErrors = append(compareErrors, formatError("MaxAge", *expectedCookie.maxAge, actualCookie.MaxAge))
	}
	return compareErrors
}

func compareExpires(expectedCookie *Cookie, actualCookie *http.Cookie, compareErrors []string) []string {
	if expectedCookie.expires != nil && !(*expectedCookie.expires).Equal(actualCookie.Expires) {
		compareErrors = append(compareErrors, formatError("Expires", *expectedCookie.expires, actualCookie.Expires))
	}
	return compareErrors
}

func comparePath(expectedCookie *Cookie, actualCookie *http.Cookie, compareErrors []string) []string {
	if expectedCookie.path != nil && *expectedCookie.path != actualCookie.Path {
		compareErrors = append(compareErrors, formatError("Path", *expectedCookie.path, actualCookie.Path))
	}
	return compareErrors
}

func compareDomain(expectedCookie *Cookie, actualCookie *http.Cookie, compareErrors []string) []string {
	if expectedCookie.domain != nil && *expectedCookie.domain != actualCookie.Domain {
		compareErrors = append(compareErrors, formatError("Domain", *expectedCookie.domain, actualCookie.Domain))
	}
	return compareErrors
}

func compareValue(expectedCookie *Cookie, actualCookie *http.Cookie, compareErrors []string) []string {
	if expectedCookie.value != nil && *expectedCookie.value != actualCookie.Value {
		compareErrors = append(compareErrors, formatError("Value", *expectedCookie.value, actualCookie.Value))
	}
	return compareErrors
}

func formatError(name string, expectedValue, actualValue interface{}) string {
	return fmt.Sprintf("Mismatched field %s. Expected %v but received %v",
		name,
		expectedValue,
		actualValue)
}
