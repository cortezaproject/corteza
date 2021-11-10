package oauth2

import (
	"context"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
)

// NewJWTAccessGenerate create to generate the jwt access token instance
//
// @todo move this to pkg/auth (??) so it can be re-used
func NewJWTAccessGenerate(kid string, key []byte, method jwt.SigningMethod) *JWTAccessGenerate {
	return &JWTAccessGenerate{
		SignedKeyID:  kid,
		SignedKey:    key,
		SignedMethod: method,
	}
}

// JWTAccessGenerate generate the jwt access token
type JWTAccessGenerate struct {
	SignedKeyID  string
	SignedKey    []byte
	SignedMethod jwt.SigningMethod
}

// Token based on the UUID generated token
func (a *JWTAccessGenerate) Token(ctx context.Context, data *oauth2.GenerateBasic, isGenRefresh bool) (string, string, error) {
	// extract user ID and roles from a space-delimited list of IDs stored in userID
	userIdWithRoles := strings.SplitN(data.TokenInfo.GetUserID(), " ", 2)
	if len(userIdWithRoles) == 1 {
		userIdWithRoles = append(userIdWithRoles, "")
	}

	// using jwt.MapClaims is good enough, it's validation rules ae
	claims := jwt.MapClaims{
		"aud":   data.Client.GetID(),
		"sub":   userIdWithRoles[0],
		"exp":   data.TokenInfo.GetAccessCreateAt().Add(data.TokenInfo.GetAccessExpiresIn()).Unix(),
		"scope": data.TokenInfo.GetScope(),
		"roles": userIdWithRoles[1],
	}

	token := jwt.NewWithClaims(a.SignedMethod, claims)
	token.Header["salt"] = string(rand.Bytes(32))
	if a.SignedKeyID != "" {
		token.Header["kid"] = a.SignedKeyID
	}
	var key interface{}
	if a.isEs() {
		v, err := jwt.ParseECPrivateKeyFromPEM(a.SignedKey)
		if err != nil {
			return "", "", err
		}
		key = v
	} else if a.isRsOrPS() {
		v, err := jwt.ParseRSAPrivateKeyFromPEM(a.SignedKey)
		if err != nil {
			return "", "", err
		}
		key = v
	} else if a.isHs() {
		key = a.SignedKey
	} else {
		return "", "", errors.New("unsupported sign method")
	}

	access, err := token.SignedString(key)
	if err != nil {
		return "", "", err
	}

	refresh := ""
	if isGenRefresh {
		refresh = string(rand.Bytes(48))
	}

	return access, refresh, nil
}

func (a *JWTAccessGenerate) isEs() bool {
	return strings.HasPrefix(a.SignedMethod.Alg(), "ES")
}

func (a *JWTAccessGenerate) isRsOrPS() bool {
	isRs := strings.HasPrefix(a.SignedMethod.Alg(), "RS")
	isPs := strings.HasPrefix(a.SignedMethod.Alg(), "PS")
	return isRs || isPs
}

func (a *JWTAccessGenerate) isHs() bool {
	return strings.HasPrefix(a.SignedMethod.Alg(), "HS")
}
