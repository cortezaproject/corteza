package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/go-chi/jwtauth"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type (
	MiddlewareValidator interface {
		HttpValidator(scope ...string) func(http.Handler) http.Handler
		Generate(ctx context.Context, i Identifiable, clientID uint64, scope ...string) (signed []byte, err error)
	}

	oauth2manager interface {
		LoadAccessToken(ctx context.Context, access string) (ti oauth2.TokenInfo, err error)
		GenerateAccessToken(ctx context.Context, gt oauth2.GrantType, tgr *oauth2.TokenGenerateRequest) (oauth2.TokenInfo, error)
	}

	jwtManager struct {
		// Expiration time in minutes
		expiry time.Duration

		signAlgo jwa.SignatureAlgorithm
		signKey  jwk.Key

		log *zap.Logger

		oa2m oauth2manager

		issuerClaim string
	}
)

var (
	defaultJWTManager *jwtManager
)

// JWT returns d
func JWT() *jwtManager {
	return defaultJWTManager
}

func SetupDefault(oa2m oauth2manager, secret string, expiry time.Duration) (err error) {
	// Use JWT secret for hmac signer for now
	DefaultSigner = HmacSigner(secret)
	defaultJWTManager, err = NewJWTManager(oa2m, jwa.HS512, secret, expiry)
	return
}

// NewJWTManager initializes and returns new instance of JWT manager
// @todo should be extended to accept different kinds of algorythms, private-keys etc.
func NewJWTManager(oa2m oauth2manager, algo jwa.SignatureAlgorithm, secret string, expiry time.Duration) (tm *jwtManager, err error) {
	tm = &jwtManager{
		expiry:      expiry,
		signAlgo:    algo,
		issuerClaim: "cortezaproject.org",
		log:         logger.Default(),
		oa2m:        oa2m,
	}

	if len(secret) == 0 {
		return nil, fmt.Errorf("JWK missing")
	}

	if tm.signKey, err = jwk.New([]byte(secret)); err != nil {
		return nil, fmt.Errorf("could not parse JWK: %w", err)
	}

	return
}

// Sign takes security information and returns signed JWT
//
// Access token is expected to be issued by OAuth2 token manager and we want to
// transport access-token one of the JWT claims (JWT ID!).
//
// This way we can perform static checks (origin, validity, expiration)
// before doing any storage lookups.
//f
// Identity holds user ID and all roles that go into this security context
// Client ID represents the auth client that was used
func (m *jwtManager) Sign(accessToken string, identity Identifiable, clientID uint64, scope ...string) (signed []byte, err error) {
	var (
		roles string
		token = jwt.New()
	)

	if len(scope) == 0 {
		// for backward compatibility we default
		// unset scope to profile & api
		scope = []string{"profile", "api"}
	}

	for _, r := range identity.Roles() {
		roles += strconv.FormatUint(r, 10)
	}

	// this is the key part
	// here we put access token to the JWT ID claim
	if err = token.Set(jwt.JwtIDKey, accessToken); err != nil {
		return
	}

	if err = token.Set(jwt.SubjectKey, identity.String()); err != nil {
		return
	}

	if err = token.Set(jwt.ExpirationKey, time.Now().Add(m.expiry).Unix()); err != nil {
		return
	}

	if err = token.Set(jwt.IssuerKey, m.issuerClaim); err != nil {
		return
	}

	if err = token.Set(jwt.IssuedAtKey, time.Now().Unix()); err != nil {
		return
	}

	if err = token.Set("clientID", strconv.FormatUint(clientID, 10)); err != nil {
		return
	}

	if err = token.Set("scope", strings.Join(scope, " ")); err != nil {
		return
	}

	if err = token.Set("roles", strings.TrimSpace(roles)); err != nil {
		return
	}

	if signed, err = jwt.Sign(token, m.signAlgo, m.signKey); err != nil {
		return
	}

	return signed, nil
}

// Generate new access-token and JWT
//
// Why so much effort and not just return the access token?
// We want to transport access-token one of the JWT claims (JWT ID!).
//
// This way we can perform static checks (origin, validity, expiration)
// before doing any storage lookups.
func (m *jwtManager) Generate(ctx context.Context, i Identifiable, clientID uint64, scope ...string) (signed []byte, err error) {
	var (
		ti oauth2.TokenInfo
	)

	ti, err = m.oa2m.GenerateAccessToken(ctx, oauth2.Implicit, &oauth2.TokenGenerateRequest{
		ClientID:       strconv.FormatUint(clientID, 10),
		UserID:         i.String(),
		Scope:          strings.Join(scope, " "),
		Refresh:        "??????????",
		AccessTokenExp: m.expiry,
	})

	if err != nil {
		return
	}

	return m.Sign(ti.GetAccess(), i, clientID, scope...)
}

// HttpVerifier http middleware handler will verify a JWT string from a http request.
func (m *jwtManager) HttpVerifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(jwtauth.New(m.signAlgo.String(), m.signKey, nil))
}

func (m *jwtManager) HttpValidator(scope ...string) func(http.Handler) http.Handler {
	if len(scope) == 0 {
		// ensure that scope is not empty
		scope = []string{"api"}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := ValidateContext(r.Context(), m.oa2m, scope...); err != nil {
				errors.ProperlyServeHTTP(w, r, err, false)
				return
			} else {
				token, _, _ := jwtauth.FromContext(r.Context())
				r = r.WithContext(SetIdentityToContext(r.Context(), IdentityFromToken(token)))
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (m *jwtManager) ValidateContext(ctx context.Context, scope ...string) error {
	return ValidateContext(ctx, m.oa2m, scope...)
}

func (m *jwtManager) Validate(ctx context.Context, token jwt.Token, scope ...string) error {
	return Validate(ctx, token, m.oa2m, scope...)
}

// ValidateContext gets JWT & claims from context
//
// It's chi middleware that puts it there
func ValidateContext(ctx context.Context, oa2m oauth2manager, scope ...string) (err error) {
	var (
		token jwt.Token
	)

	if token, _, err = jwtauth.FromContext(ctx); err != nil {
		return ErrUnauthorized()
	}

	return Validate(ctx, token, oa2m, scope...)
}

// Validate performs token validation
//
// Steps:
//   - check scope in the JWT
//   - check if JWT ID is set (where the access-token string is stored)
//   - check if access-token exists in the DB
//
//
func Validate(ctx context.Context, token jwt.Token, oa2m oauth2manager, scope ...string) (err error) {
	if len(scope) > 0 && !CheckJwtScope(token, scope...) {
		return ErrUnauthorizedScope()
	}

	// Extract the JWT id from the token (string) and convert it to uint64
	// to be compatible with the lookup function
	if len(token.JwtID()) < 10 {
		return ErrMalformedToken("missing or malformed JWT ID")
	}

	// @todo we could use a simple caching mechanism here
	//       1. if lookup is successful, add a JWT ID to the list
	//       2. add short exp time (that should not last longer than token's exp time)
	//       3. check against the list first; if JWT ID is not present there check in storage
	//
	if _, err = oa2m.LoadAccessToken(ctx, token.JwtID()); err != nil {
		return ErrUnauthorized()
	}

	return nil
}

// IdentityFromToken decodes sub & roles claims into identity
func IdentityFromToken(token jwt.Token) *identity {
	var (
		roles, _ = token.Get("roles")
	)

	return Authenticated(
		cast.ToUint64(token.Subject()),
		payload.ParseUint64s(strings.Split(cast.ToString(roles), " "))...,
	)
}
