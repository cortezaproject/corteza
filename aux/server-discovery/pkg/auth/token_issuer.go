package auth

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/spf13/cast"
)

type (
	tokenIssuer struct {
		defaultRequest *TokenRequest

		// store issued tokens
		store tokenIssuerStore

		// lookup for issued tokens
		lookup tokenIssuerLookup

		// generator for issued tokens
		generator tokenIssuerGenerator

		// signer for issued tokens
		signer tokenIssuerSigner
	}

	// IssuerOptFn modify toeknIssuer
	IssuerOptFn func(*tokenIssuer) error

	TokenRequest struct {
		AccessToken  string
		RefreshToken string
		Expiration   time.Duration
		Audience     string
		Issuer       string
		IssuedAt     time.Time
		ClientID     uint64
		UserID       uint64
		Roles        []uint64
		Scope        []string
	}

	// IssueOptFn functions modify TokenRequest
	IssueOptFn func(*TokenRequest) error

	tokenIssuerStore     func(context.Context, TokenRequest) error
	tokenIssuerLookup    func(context.Context, string) error
	tokenIssuerGenerator func(context.Context, TokenRequest) (string, string, error)
	tokenIssuerSigner    func(token jwt.Token) ([]byte, error)
)

var (
	TokenIssuer *tokenIssuer

	// wrapper around time.Now() that will aid service testing
	now = func() *time.Time {
		c := time.Now().Truncate(time.Second)
		return &c
	}
)

// NewTokenIssuer initializes and returns new instance of JWT manager
func NewTokenIssuer(opt ...IssuerOptFn) (issuer *tokenIssuer, err error) {
	issuer = &tokenIssuer{
		defaultRequest: &TokenRequest{Issuer: "cortezaproject.org"},

		store: func(ctx context.Context, request TokenRequest) error {
			// elegantly handle unconfigured store
			return fmt.Errorf("token issuer store not configured")
		},

		lookup: func(context.Context, string) error {
			// elegantly handle unconfigured lookup
			return fmt.Errorf("token issuer lookup not configured")
		},

		signer: func(token jwt.Token) ([]byte, error) {
			// elegantly handle unconfigured signer
			return nil, fmt.Errorf("token issuer signer not configured")
		},

		generator: DefaultAccessTokenGenerator,
	}

	return issuer, applyToIssuer(issuer, opt...)
}

// Issue issues new access token, stores it and returns signed JWT.
func (i *tokenIssuer) Issue(ctx context.Context, opt ...IssueOptFn) (_ []byte, err error) {
	var req = i.newTokenRequest()
	if err = req.apply(opt...); err != nil {
		return
	}

	if len(req.AccessToken+req.RefreshToken) > 0 {
		panic("can not issue new token with preset access and refresh tokens, " +
			"this is most likely an implementation mistake")
	}

	if req.AccessToken, req.RefreshToken, err = i.generator(ctx, *req); err != nil {
		return
	}

	if err = i.store(ctx, *req); err != nil {
		return
	}

	return i.sign(req)
}

func (i *tokenIssuer) Sign(opt ...IssueOptFn) (_ []byte, err error) {
	var req = i.newTokenRequest()
	if err = req.apply(opt...); err != nil {
		return
	}

	return i.sign(req)
}

func (i *tokenIssuer) sign(req *TokenRequest) ([]byte, error) {
	if token, err := makeToken(req); err != nil {
		return nil, err
	} else {
		return i.signer(token)
	}
}

// Returns new token request and copies relevant values from the default token request on the issuer
func (i *tokenIssuer) newTokenRequest() *TokenRequest {
	return &TokenRequest{
		Issuer:     i.defaultRequest.Issuer,
		ClientID:   i.defaultRequest.ClientID,
		Expiration: i.defaultRequest.Expiration,
		IssuedAt:   *now(),
	}
}

// Validate performs token validation by checking existence of access-token in the store
func (i *tokenIssuer) Validate(ctx context.Context, token jwt.Token) (err error) {
	if err = i.lookup(ctx, token.JwtID()); err != nil {
		return fmt.Errorf("unauthorized")
	}

	return nil
}

func makeToken(req *TokenRequest) (_ jwt.Token, err error) {
	var (
		roles = make([]string, len(req.Roles))
		token = jwt.New()

		toString = func(i uint64) string {
			return strconv.FormatUint(i, 10)
		}
	)

	if len(req.Scope) == 0 {
		// for backward compatibility we default
		// unset scope to profile & api
		req.Scope = []string{"profile", "api"}
	}

	if req.IssuedAt.IsZero() {
		req.IssuedAt = *now()
	}

	for i, r := range req.Roles {
		roles[i] = toString(r)
	}

	// The key part: store access token as JWT ID claim.
	// Claim will be extracted when JWT is validated and checked
	if err = token.Set(jwt.JwtIDKey, req.AccessToken); err != nil {
		return
	}

	if err = token.Set(jwt.SubjectKey, toString(req.UserID)); err != nil {
		return
	}

	if err = token.Set(jwt.ExpirationKey, now().Add(req.Expiration).Unix()); err != nil {
		return
	}

	if req.Audience != "" {
		if err = token.Set(jwt.AudienceKey, req.Audience); err != nil {
			return
		}
	}

	if err = token.Set(jwt.IssuerKey, req.Issuer); err != nil {
		return
	}

	if err = token.Set(jwt.IssuedAtKey, req.IssuedAt.Unix()); err != nil {
		return
	}

	if err = token.Set("clientID", toString(req.ClientID)); err != nil {
		return
	}

	if err = token.Set("scope", strings.Join(req.Scope, " ")); err != nil {
		return
	}

	if err = token.Set("roles", roles); err != nil {
		return
	}

	return token, nil
}

// IdentityFromToken decodes sub & roles claims into identity
func IdentityFromToken(token jwt.Token) *identity {
	if token == nil {
		return Anonymous()
	}

	var (
		roles, _ = token.Get("roles")
	)

	return Authenticated(
		cast.ToUint64(token.Subject()),
		payload.ParseUint64s(cast.ToStringSlice(roles))...,
	)
}

// DefaultAccessTokenGenerator uses token generator from oauth2 lib
func DefaultAccessTokenGenerator(ctx context.Context, req TokenRequest) (string, string, error) {
	return generates.NewAccessGenerate().Token(
		ctx,
		&oauth2.GenerateBasic{
			Client: &models.Client{
				ID: strconv.FormatUint(req.ClientID, 10),
			},
			UserID:    strconv.FormatUint(req.UserID, 10),
			CreateAt:  *now(),
			TokenInfo: nil,
			Request:   nil,
		},
		true,
	)
}

// WithSecretSigner configures token issuer with
func WithSecretSigner(secret string) IssuerOptFn {
	return func(tm *tokenIssuer) (err error) {
		if len(secret) == 0 {
			return fmt.Errorf("JWK missing")
		}

		var key jwk.Key
		if key, err = jwk.New([]byte(secret)); err != nil {
			return fmt.Errorf("could not parse JWK: %w", err)
		}

		tm.signer = func(token jwt.Token) ([]byte, error) {
			return jwt.Sign(token, jwa.HS512, key)
		}

		return nil
	}
}

// WithDefaultClientID configures ID of the default auth client
func WithDefaultClientID(ID uint64) IssuerOptFn {
	return func(tm *tokenIssuer) (err error) {
		tm.defaultRequest.ClientID = ID
		return
	}
}

// WithDefaultExpiration configures default token expiration time
func WithDefaultExpiration(exp time.Duration) IssuerOptFn {
	return func(tm *tokenIssuer) (err error) {
		tm.defaultRequest.Expiration = exp
		return
	}
}

// WithDefaultIssuer configures default issuer claim
func WithDefaultIssuer(iss string) IssuerOptFn {
	return func(tm *tokenIssuer) (err error) {
		tm.defaultRequest.Issuer = iss
		return
	}
}

// WithStore configures store function
func WithStore(fn tokenIssuerStore) IssuerOptFn {
	return func(tm *tokenIssuer) (err error) {
		tm.store = fn
		return
	}
}

// WithLookup configures lookup function
func WithLookup(fn tokenIssuerLookup) IssuerOptFn {
	return func(tm *tokenIssuer) (err error) {
		tm.lookup = fn
		return
	}
}

// WithGenerator configures generator function
func WithGenerator(fn tokenIssuerGenerator) IssuerOptFn {
	return func(tm *tokenIssuer) (err error) {
		tm.generator = fn
		return
	}
}

// WithSigner configures signer function
func WithSigner(fn tokenIssuerSigner) IssuerOptFn {
	return func(tm *tokenIssuer) (err error) {
		tm.signer = fn
		return
	}
}

func (req *TokenRequest) apply(opt ...IssueOptFn) (err error) {
	for _, fn := range opt {
		if err = fn(req); err != nil {
			return
		}
	}

	return
}

func WithExpiration(e time.Duration) IssueOptFn {
	return func(t *TokenRequest) (err error) {
		t.Expiration = e
		return
	}
}

func WithIdentity(i Identifiable) IssueOptFn {
	return func(t *TokenRequest) (err error) {
		t.UserID = i.Identity()
		t.Roles = i.Roles()
		return
	}
}

func WithAccessToken(access string) IssueOptFn {
	return func(t *TokenRequest) (err error) {
		t.AccessToken = access
		return
	}
}

func WithScope(ss ...string) IssueOptFn {
	return func(t *TokenRequest) (err error) {
		t.Scope = ss
		return
	}
}

func WithAudience(aud string) IssueOptFn {
	return func(t *TokenRequest) (err error) {
		t.Audience = aud
		return
	}
}

func WithClientID(id uint64) IssueOptFn {
	return func(t *TokenRequest) (err error) {
		t.ClientID = id
		return
	}
}

func applyToIssuer(tm *tokenIssuer, opt ...IssuerOptFn) (err error) {
	for _, fn := range opt {
		if err = fn(tm); err != nil {
			return
		}
	}

	return
}
