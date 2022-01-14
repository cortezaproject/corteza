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
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-chi/jwtauth"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type (
	signer interface {
		Sign(accessToken string, identity Identifiable, clientID uint64, scope ...string) (signed []byte, err error)
	}

	MiddlewareValidator interface {
		HttpValidator(scope ...string) func(http.Handler) http.Handler
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

	// @todo remove
	tokenStore interface {
		CreateAuthOa2token(ctx context.Context, rr ...*types.AuthOa2token) error
		UpsertAuthConfirmedClient(ctx context.Context, rr ...*types.AuthConfirmedClient) error
	}

	// @todo remove
	//tokenLookup interface {
	//	LookupAuthOa2tokenByID(ctx context.Context, id uint64) (*types.AuthOa2token, error)
	//}
	//
	//tokenStoreWithLookup interface {
	//	tokenStore
	//	tokenLookup
	//}
)

var (
	defaultJWTManager *jwtManager
	//DefaultJwtStore   tokenStoreWithLookup
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

//// @todo remove
////// SetJWTStore set store for JWT
////// @todo find better way to initiate store,
////// 			it mainly used for generating and storing accessToken for impersonate and corredor, Ref: j.Generate()
////func SetJWTStore(store tokenStoreWithLookup) {
////	DefaultJwtStore = store
////}
//
//// Authenticate the token from the given string and return parsed token or error
//func (m *jwtManager) Authenticate(s string) (pToken jwt.Token, err error) {
//	if pToken, err = jwt.Parse([]byte(s), jwt.WithVerify(m.signAlgo, m.signKey)); err != nil {
//		return
//	}
//
//	if err = jwt.Validate(pToken); err != nil {
//		return
//	}
//
//	return
//}

// Sign takes security information and returns signed JWT
//
// Access token is expected to be issued by OAuth2 token manager
// without it, we can only do static (JWT itself) validation
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

func (m *jwtManager) Generate(ctx context.Context, i Identifiable, clientID uint64, scope ...string) (signed []byte, err error) {
	var (
		ti oauth2.TokenInfo
	)

	ti, err = m.oa2m.GenerateAccessToken(ctx, oauth2.Implicit, &oauth2.TokenGenerateRequest{
		ClientID:       strconv.FormatUint(clientID, 10),
		UserID:         i.String(),
		Scope:          strings.Join(scope, " "),
		Refresh:        "cli?",
		AccessTokenExp: m.expiry,
	})

	if err != nil {
		return
	}

	return m.Sign(ti.GetAccess(), i, 0, scope...)
}

func ValidateContext(ctx context.Context, oa2m oauth2manager, scope ...string) (err error) {
	var (
		token jwt.Token
	)

	if token, _, err = jwtauth.FromContext(ctx); err != nil {
		return ErrUnauthorized()
	}

	return Validate(ctx, token, oa2m, scope...)
}

func Validate(ctx context.Context, token jwt.Token, oa2m oauth2manager, scope ...string) (err error) {
	if !CheckJwtScope(token, scope...) {
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

//// HttpAuthenticator converts JWT claims into identity and stores it into context
//func (m *jwtManager) HttpAuthenticator(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		ctx := r.Context()
//
//		tkn, _, err := jwtauth.FromContext(ctx)
//
//		// Requests w/o token should not yield an error
//		// there are parts of the system that can be access without it
//		// and/or handle such situation internally
//		if err != nil && !errors.Is(err, jwtauth.ErrNoTokenFound) {
//			api.Send(w, r, err)
//			return
//		}
//
//		// If token is present extract identity
//		if tkn != nil {
//			ctx = SetIdentityToContext(ctx, IdentityFromToken(tkn))
//			r = r.WithContext(ctx)
//
//			// @todo verify JWT ID (access-token!!
//			tkn.JwtID()
//
//		}
//
//		next.ServeHTTP(w, r)
//	})
//}

//
//// Generate makes a new token and stores it in the database
//func (tm *tokenManager) Generate(ctx context.Context, i Identifiable, clientID uint64, scope ...string) (token []byte, err error) {
//	var (
//		//	eti  = GetExtraReqInfoFromContext(ctx)
//		//	oa2t = &types.AuthOa2token{
//		//		ID:         id.Next(),
//		//		CreatedAt:  time.Now().Round(time.Second),
//		//		RemoteAddr: eti.RemoteAddr,
//		//		UserAgent:  eti.UserAgent,
//		//		ClientID:   clientID,
//		//	}
//		//
//		//	acc = &types.AuthConfirmedClient{
//		//		ConfirmedAt: oa2t.CreatedAt,
//		//		ClientID:    clientID,
//		//	}
//		oa2t *types.AuthOa2token
//		acc  *types.AuthConfirmedClient
//
//		jwtID = id.Next()
//	)
//
//	if oa2t, acc, err = MakeAuthStructs(ctx, jwtID, i.Identity(), clientID, nil, tm.expiry); err != nil {
//		return
//	}
//
//	if token, err = tm.make(jwtID, i, clientID, scope...); err != nil {
//		return nil, err
//	}
//
//	oa2t.Access = string(token)
//
//	// use the same expiration as on token
//	//oa2t.ExpiresAt = oa2t.CreatedAt.Add(tm.expiry)
//
//	//if oa2t.Data, err = json.Marshal(oa2t); err != nil {
//	//	return
//	//}
//
//	//if oa2t.UserID, _ = ExtractFromSubClaim(i.String()); oa2t.UserID == 0 {
//	//	// UserID stores collection of IDs: user's ID and set of all roles' user is member of
//	//	return nil, fmt.Errorf("could not parse user ID from token")
//	//}
//	//
//	//// copy user id to auth client confirmation
//	//acc.UserID = oa2t.UserID
//
//	return token, StoreAuthToken(ctx, DefaultJwtStore, oa2t, acc)
//}

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
