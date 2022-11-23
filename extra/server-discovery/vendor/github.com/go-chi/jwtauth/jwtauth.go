package jwtauth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

type JWTAuth struct {
	alg       jwa.SignatureAlgorithm
	signKey   interface{} // private-key
	verifyKey interface{} // public-key, only used by RSA and ECDSA algorithms
	verifier  jwt.ParseOption
}

var (
	TokenCtxKey = &contextKey{"Token"}
	ErrorCtxKey = &contextKey{"Error"}
)

var (
	ErrUnauthorized = errors.New("token is unauthorized")
	ErrExpired      = errors.New("token is expired")
	ErrNBFInvalid   = errors.New("token nbf validation failed")
	ErrIATInvalid   = errors.New("token iat validation failed")
	ErrNoTokenFound = errors.New("no token found")
	ErrAlgoInvalid  = errors.New("algorithm mismatch")
)

func New(alg string, signKey interface{}, verifyKey interface{}) *JWTAuth {
	ja := &JWTAuth{alg: jwa.SignatureAlgorithm(alg), signKey: signKey, verifyKey: verifyKey}

	if ja.verifyKey != nil {
		ja.verifier = jwt.WithVerify(ja.alg, ja.verifyKey)
	} else {
		ja.verifier = jwt.WithVerify(ja.alg, ja.signKey)
	}

	return ja
}

// Verifier http middleware handler will verify a JWT string from a http request.
//
// Verifier will search for a JWT token in a http request, in the order:
//   1. 'jwt' URI query parameter
//   2. 'Authorization: BEARER T' request header
//   3. Cookie 'jwt' value
//
// The first JWT string that is found as a query parameter, authorization header
// or cookie header is then decoded by the `jwt-go` library and a *jwt.Token
// object is set on the request context. In the case of a signature decoding error
// the Verifier will also set the error on the request context.
//
// The Verifier always calls the next http handler in sequence, which can either
// be the generic `jwtauth.Authenticator` middleware or your own custom handler
// which checks the request context jwt token and error to prepare a custom
// http response.
func Verifier(ja *JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Verify(ja, TokenFromHeader, TokenFromCookie)(next)
	}
}

func Verify(ja *JWTAuth, findTokenFns ...func(r *http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token, err := VerifyRequest(ja, r, findTokenFns...)
			ctx = NewContext(ctx, token, err)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

func VerifyRequest(ja *JWTAuth, r *http.Request, findTokenFns ...func(r *http.Request) string) (jwt.Token, error) {
	var tokenString string

	// Extract token string from the request by calling token find functions in
	// the order they where provided. Further extraction stops if a function
	// returns a non-empty string.
	for _, fn := range findTokenFns {
		tokenString = fn(r)
		if tokenString != "" {
			break
		}
	}
	if tokenString == "" {
		return nil, ErrNoTokenFound
	}

	return VerifyToken(ja, tokenString)
}

func VerifyToken(ja *JWTAuth, tokenString string) (jwt.Token, error) {
	// Decode & verify the token
	token, err := ja.Decode(tokenString)
	if err != nil {
		return token, ErrorReason(err)
	}

	if token == nil {
		return nil, ErrUnauthorized
	}

	if err := jwt.Validate(token); err != nil {
		return token, ErrorReason(err)
	}

	// Valid!
	return token, nil
}

func (ja *JWTAuth) Encode(claims map[string]interface{}) (t jwt.Token, tokenString string, err error) {
	t = jwt.New()
	for k, v := range claims {
		t.Set(k, v)
	}
	payload, err := ja.sign(t)
	if err != nil {
		return nil, "", err
	}
	tokenString = string(payload)
	return
}

func (ja *JWTAuth) Decode(tokenString string) (jwt.Token, error) {
	return ja.parse([]byte(tokenString))
}

func (ja *JWTAuth) sign(token jwt.Token) ([]byte, error) {
	return jwt.Sign(token, ja.alg, ja.signKey)
}

func (ja *JWTAuth) parse(payload []byte) (jwt.Token, error) {
	return jwt.Parse(payload, ja.verifier)
}

// ErrorReason will normalize the error message from the underlining
// jwt library
func ErrorReason(err error) error {
	switch err.Error() {
	case "exp not satisfied", ErrExpired.Error():
		return ErrExpired
	case "iat not satisfied", ErrIATInvalid.Error():
		return ErrIATInvalid
	case "nbf not satisfied", ErrNBFInvalid.Error():
		return ErrNBFInvalid
	default:
		return ErrUnauthorized
	}
}

// Authenticator is a default authentication middleware to enforce access from the
// Verifier middleware request context values. The Authenticator sends a 401 Unauthorized
// response for any unverified tokens and passes the good ones through. It's just fine
// until you decide to write something similar and customize your client response.
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := FromContext(r.Context())

		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

func NewContext(ctx context.Context, t jwt.Token, err error) context.Context {
	ctx = context.WithValue(ctx, TokenCtxKey, t)
	ctx = context.WithValue(ctx, ErrorCtxKey, err)
	return ctx
}

func FromContext(ctx context.Context) (jwt.Token, map[string]interface{}, error) {
	token, _ := ctx.Value(TokenCtxKey).(jwt.Token)

	var err error
	var claims map[string]interface{}

	if token != nil {
		claims, err = token.AsMap(context.Background())
		if err != nil {
			return token, nil, err
		}
	} else {
		claims = map[string]interface{}{}
	}

	err, _ = ctx.Value(ErrorCtxKey).(error)

	return token, claims, err
}

// UnixTime returns the given time in UTC milliseconds
func UnixTime(tm time.Time) int64 {
	return tm.UTC().Unix()
}

// EpochNow is a helper function that returns the NumericDate time value used by the spec
func EpochNow() int64 {
	return time.Now().UTC().Unix()
}

// ExpireIn is a helper function to return calculated time in the future for "exp" claim
func ExpireIn(tm time.Duration) int64 {
	return EpochNow() + int64(tm.Seconds())
}

// Set issued at ("iat") to specified time in the claims
func SetIssuedAt(claims map[string]interface{}, tm time.Time) {
	claims["iat"] = tm.UTC().Unix()
}

// Set issued at ("iat") to present time in the claims
func SetIssuedNow(claims map[string]interface{}) {
	claims["iat"] = EpochNow()
}

// Set expiry ("exp") in the claims
func SetExpiry(claims map[string]interface{}, tm time.Time) {
	claims["exp"] = tm.UTC().Unix()
}

// Set expiry ("exp") in the claims to some duration from the present time
func SetExpiryIn(claims map[string]interface{}, tm time.Duration) {
	claims["exp"] = ExpireIn(tm)
}

// TokenFromCookie tries to retreive the token string from a cookie named
// "jwt".
func TokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// TokenFromHeader tries to retreive the token string from the
// "Authorization" reqeust header: "Authorization: BEARER T".
func TokenFromHeader(r *http.Request) string {
	// Get token from authorization header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

// TokenFromQuery tries to retreive the token string from the "jwt" URI
// query parameter.
//
// To use it, build our own middleware handler, such as:
//
// func Verifier(ja *JWTAuth) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return Verify(ja, TokenFromQuery, TokenFromHeader, TokenFromCookie)(next)
// 	}
// }
func TokenFromQuery(r *http.Request) string {
	// Get token from query param named "jwt".
	return r.URL.Query().Get("jwt")
}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "jwtauth context value " + k.name
}
