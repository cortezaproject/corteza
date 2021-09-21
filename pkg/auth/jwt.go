package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/system/types"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/pkg/errors"
)

type (
	token struct {
		// Expiration time in minutes
		expiry    time.Duration
		tokenAuth *jwtauth.JWTAuth
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

func SetupDefault(secret string, expiry time.Duration) {
	// Use JWT secret for hmac signer for now
	DefaultSigner = HmacSigner(secret)
	DefaultJwtHandler, _ = JWT(secret, expiry)
}

func JWT(secret string, expiry time.Duration) (tkn *token, err error) {
	if len(secret) == 0 {
		return nil, errors.New("JWT secret missing")
	}

	tkn = &token{
		expiry:    expiry,
		tokenAuth: jwtauth.New(jwt.SigningMethodHS512.Alg(), []byte(secret), nil),
	}

	return tkn, nil
}

// SetJWTStore set store for JWT
// @todo find better way to initiate store,
// 			it mainly used for generating and storing accessToken for impersonate and corredor, Ref: j.Generate()
func SetJWTStore(store tokenStore) {
	DefaultJwtStore = store
}

func (t *token) Authenticate(token string) (jwt.MapClaims, error) {
	dt, err := t.tokenAuth.Decode(token)
	if err != nil {
		return nil, err
	}

	if dt == nil || !dt.Valid {
		return nil, jwtauth.ErrUnauthorized
	}

	if dt.Method != jwt.SigningMethodHS512 {
		return nil, jwtauth.ErrAlgoInvalid
	}

	if mc, is := dt.Claims.(jwt.MapClaims); is {
		return mc, nil
	}

	return nil, nil
}

// HttpVerifier returns a HTTP handler that verifies JWT and stores it into context
func (t *token) HttpVerifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(t.tokenAuth)
}

func (t *token) Encode(i Identifiable, scope ...string) string {
	var (
		// when possible, extend this with the client
		clientID uint64 = 0
	)

	if len(scope) == 0 {
		// for backward compatibility we default
		// unset scope to profile & api
		scope = []string{"profile", "api"}
	}

	return t.encode(i, clientID, scope...)
}

func (t *token) encode(i Identifiable, clientID uint64, scope ...string) string {
	roles := ""
	for _, r := range i.Roles() {
		roles += fmt.Sprintf(" %d", r)
	}

	_, tkn, _ := t.tokenAuth.Encode(jwt.MapClaims{
		"sub":   i.String(),
		"exp":   time.Now().Add(t.expiry).Unix(),
		"aud":   fmt.Sprintf("%d", clientID),
		"scope": strings.Join(scope, " "),
		"roles": strings.TrimSpace(roles),
	})

	return tkn
}

// HttpAuthenticator converts JWT claims into identity and stores it into context
func (t *token) HttpAuthenticator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			tkn, claims, err := jwtauth.FromContext(ctx)

			// When token is present, expect no errors and valid claims!
			if tkn != nil {
				if err != nil {
					// But if token is present, the shouldn't be an error
					api.Send(w, r, err)
					return
				}

				ctx = SetIdentityToContext(ctx, ClaimsToIdentity(claims))
				ctx = context.WithValue(ctx, scopeCtxKey{}, claims["scope"])

				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (t *token) Generate(ctx context.Context, i Identifiable) (tokenString string, err error) {
	var (
		eti  = GetExtraReqInfoFromContext(ctx)
		oa2t = &types.AuthOa2token{
			ID:         id.Next(),
			CreatedAt:  time.Now().Round(time.Second),
			RemoteAddr: eti.RemoteAddr,
			UserAgent:  eti.UserAgent,
		}

		acc = &types.AuthConfirmedClient{
			ConfirmedAt: oa2t.CreatedAt,
		}
	)

	tokenString = t.Encode(i)
	oa2t.Access = tokenString
	oa2t.ExpiresAt = oa2t.CreatedAt.Add(t.expiry)

	if oa2t.Data, err = json.Marshal(oa2t); err != nil {
		return
	}

	// extend this with the client
	oa2t.ClientID = 0

	// copy client id to auth client confirmation
	acc.ClientID = oa2t.ClientID

	if oa2t.UserID, _ = ExtractFromSubClaim(i.String()); oa2t.UserID == 0 {
		// UserID stores collection of IDs: user's ID and set of all roles' user is member of
		return "", fmt.Errorf("could not parse user ID from token")
	}

	// copy user id to auth client confirmation
	acc.UserID = oa2t.UserID

	if err = DefaultJwtStore.UpsertAuthConfirmedClient(ctx, acc); err != nil {
		return
	}

	return tokenString, DefaultJwtStore.CreateAuthOa2token(ctx, oa2t)
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
func ClaimsToIdentity(c jwt.MapClaims) (i *identity) {
	var (
		aux       interface{}
		ok        bool
		id, roles string
	)

	if aux, ok = c["sub"]; !ok {
		return
	}

	i = &identity{}
	if id, ok = aux.(string); ok {
		i.id, _ = strconv.ParseUint(id, 10, 64)
	}

	if i.id == 0 {
		// pointless to decode roles if id is 0
		return nil
	}

	if aux, ok = c["roles"]; !ok {
		return
	}

	if roles, ok = aux.(string); !ok {
		return
	}

	for _, role := range strings.Split(roles, " ") {
		if roleID, _ := strconv.ParseUint(role, 10, 64); roleID != 0 {
			i.memberOf = append(i.memberOf, roleID)
		}
	}

	return Authenticated(i.id, i.memberOf...)
}
