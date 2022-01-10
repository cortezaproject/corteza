package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/spf13/cast"
)

type (
	tokenManager struct {
		// Expiration time in minutes
		expiry time.Duration

		signAlgo jwa.SignatureAlgorithm
		signKey  jwk.Key
	}

	tokenStore interface {
		LookupUserByID(ctx context.Context, id uint64) (*types.User, error)
		LookupAuthOa2tokenByAccess(ctx context.Context, access string) (*types.AuthOa2token, error)
		SearchRoleMembers(ctx context.Context, f types.RoleMemberFilter) (types.RoleMemberSet, types.RoleMemberFilter, error)

		CreateAuthOa2token(ctx context.Context, rr ...*types.AuthOa2token) error
		DeleteAuthOA2TokenByUserID(ctx context.Context, _userID uint64) error

		UpsertAuthConfirmedClient(ctx context.Context, rr ...*types.AuthConfirmedClient) error
	}

	ExtraReqInfo struct {
		RemoteAddr string
		UserAgent  string
	}
)

var (
	DefaultJwtHandler TokenHandler
	DefaultJwtStore   tokenStore
)

func SetupDefault(secret string, expiry time.Duration) (err error) {
	// Use JWT secret for hmac signer for now
	DefaultSigner = HmacSigner(secret)
	DefaultJwtHandler, err = TokenManager(secret, expiry)
	return
}

// TokenManager returns token management facility
// @todo should be extended to accept different kinds of algorythms, private-keys etc.
func TokenManager(secret string, expiry time.Duration) (tm *tokenManager, err error) {
	tm = &tokenManager{
		expiry:   expiry,
		signAlgo: jwa.HS512,
	}

	if len(secret) == 0 {
		return nil, fmt.Errorf("JWK missing")
	}

	if tm.signKey, err = jwk.New([]byte(secret)); err != nil {
		return nil, fmt.Errorf("could not parse JWK: %w", err)
	}

	return

	//
	//var (
	//	//	tuukn  = jwt.New()
	//	//	signed []byte
	//	//)
	//	//
	//	//if err = tuukn.Set(jwt.ExpirationKey, expiry); err != nil {
	//	//	return
	//	//}
	//	//
	//	signed, err = jwt.Sign(tuukn, jwa.HS512, []byte(secret))
	//
	//	tkn = &tokenManager{
	//		expiry:    expiry,
	//		tokenAuth: jwtauth.New(jwt.SigningMethodHS512.Alg(), []byte(secret), nil),
	//		secret:    []byte(secret),
	//	}, nil
	//
	//return tkn, nil
}

// SetJWTStore set store for JWT
// @todo find better way to initiate store,
// 			it mainly used for generating and storing accessToken for impersonate and corredor, Ref: j.Generate()
func SetJWTStore(store tokenStore) {
	DefaultJwtStore = store
}

// Authenticate the token from the given string and return parsed token or error
func (tm *tokenManager) Authenticate(token string) (pToken jwt.Token, err error) {
	if pToken, err = jwt.Parse([]byte(token), jwt.WithVerify(tm.signAlgo, tm.signKey)); err != nil {
		return
	}

	if err = jwt.Validate(pToken); err != nil {
		return
	}

	return
}

//// Encode identity into a
//func (tm *tokenManager) Encode(identity Identifiable, scope ...string) ([]byte, error) {
//	var (
//		// when possible, extend this with the client
//		clientID uint64 = 0
//	)
//
//	if len(scope) == 0 {
//		// for backward compatibility we default
//		// unset scope to profile & api
//		scope = []string{"profile", "api"}
//	}
//
//	return tm.Encode(identity, clientID, scope...)
//}

// Encode give identity, clientID & scope into JWT access token (that can be use for API requests)
//
// @todo this follows implementation in auth/oauth2/jwt_access.go
//       and should be refactored accordingly (move both into the same location/pkg => here)
func (tm *tokenManager) Encode(identity Identifiable, clientID uint64, scope ...string) (_ []byte, err error) {
	var (
		token = jwt.New()
		roles = ""
	)

	if len(scope) == 0 {
		// for backward compatibility we default
		// unset scope to profile & api
		scope = []string{"profile", "api"}
	}

	for _, r := range identity.Roles() {
		roles += fmt.Sprintf(" %d", r)
	}

	// previous implementation had special a "salt" claim that ensured JWT uniquness
	// we're using more standard approach with JWT ID now.
	if err = token.Set(jwt.JwtIDKey, fmt.Sprintf("%d", id.Next())); err != nil {
		return
	}

	if err = token.Set(jwt.SubjectKey, identity.String()); err != nil {
		return
	}

	if err = token.Set(jwt.ExpirationKey, time.Now().Add(tm.expiry).Unix()); err != nil {
		return
	}

	if err = token.Set(jwt.AudienceKey, fmt.Sprintf("%d", clientID)); err != nil {
		return
	}

	if err = token.Set("scope", strings.Join(scope, " ")); err != nil {
		return
	}

	if err = token.Set("roles", strings.TrimSpace(roles)); err != nil {
		return
	}

	return jwt.Sign(token, tm.signAlgo, tm.signKey)

	//claims := jwt.MapClaims{
	//	"sub":   identity.String(),
	//	"exp":   time.Now().Add(tm.expiry).Unix(),
	//	"aud":   fmt.Sprintf("%d", clientID),
	//	"scope": strings.Join(scope, " "),
	//	"roles": strings.TrimSpace(roles),
	//}
	//
	//newToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	//newToken.Header["salt"] = string(rand.Bytes(32))
	//access, _ := newToken.SignedString(tm.secret)
	//return access
}

// HttpVerifier returns a HTTP handler that verifies JWT and stores it into context
func (tm *tokenManager) HttpVerifier() func(http.Handler) http.Handler {
	////jwt.WithHTTPClient()
	//return func(next http.Handler) http.Handler {
	//	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	//		token, err := jwt.ParseRequest(req)
	//		if err != nil {
	//
	//		}
	//
	//		next.ServeHTTP(w, req)
	//	})
	//}

	return jwtauth.Verifier(jwtauth.New(tm.signAlgo.String(), tm.signKey, nil))
}

// HttpAuthenticator converts JWT claims into identity and stores it into context
func (tm *tokenManager) HttpAuthenticator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			tkn, _, err := jwtauth.FromContext(ctx)

			// When token is present, expect no errors and valid claims!
			if tkn != nil {
				if err != nil {
					// But if token is present, there shouldn't be an error
					api.Send(w, r, err)
					return
				}

				ctx = SetIdentityToContext(ctx, IdentityFromToken(tkn))
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Generates JWT and stores alongside with client-confirmation entry,
func (tm *tokenManager) Generate(ctx context.Context, i Identifiable, clientID uint64, scope ...string) (token []byte, err error) {
	var (
		eti  = GetExtraReqInfoFromContext(ctx)
		oa2t = &types.AuthOa2token{
			ID:         id.Next(),
			CreatedAt:  time.Now().Round(time.Second),
			RemoteAddr: eti.RemoteAddr,
			UserAgent:  eti.UserAgent,
			ClientID:   clientID,
		}

		acc = &types.AuthConfirmedClient{
			ConfirmedAt: oa2t.CreatedAt,
			ClientID:    clientID,
		}
	)

	if token, err = tm.Encode(i, clientID, scope...); err != nil {
		return
	}

	oa2t.Access = string(token)

	// use the same expiration as on token
	oa2t.ExpiresAt = oa2t.CreatedAt.Add(tm.expiry)

	if oa2t.Data, err = json.Marshal(oa2t); err != nil {
		return
	}

	if oa2t.UserID, _ = ExtractFromSubClaim(i.String()); oa2t.UserID == 0 {
		// UserID stores collection of IDs: user's ID and set of all roles' user is member of
		return nil, fmt.Errorf("could not parse user ID from token")
	}

	// copy user id to auth client confirmation
	acc.UserID = oa2t.UserID

	if err = DefaultJwtStore.UpsertAuthConfirmedClient(ctx, acc); err != nil {
		return
	}

	return token, DefaultJwtStore.CreateAuthOa2token(ctx, oa2t)
}

func GetExtraReqInfoFromContext(ctx context.Context) ExtraReqInfo {
	eti := ctx.Value(ExtraReqInfo{})
	if eti != nil {
		return eti.(ExtraReqInfo)
	} else {
		return ExtraReqInfo{}
	}
}

// ClaimsToIdentity decodes sub & roles claims into identity
func IdentityFromToken(token jwt.Token) *identity {
	var (
		roles, _ = token.Get("roles")
	)

	return Authenticated(
		cast.ToUint64(token.Subject()),
		payload.ParseUint64s(strings.Split(cast.ToString(roles), " "))...,
	)
}
