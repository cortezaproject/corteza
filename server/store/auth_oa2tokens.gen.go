package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/auth_oa2tokens.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	AuthOa2tokens interface {
		SearchAuthOa2tokens(ctx context.Context, f types.AuthOa2tokenFilter) (types.AuthOa2tokenSet, types.AuthOa2tokenFilter, error)
		LookupAuthOa2tokenByID(ctx context.Context, id uint64) (*types.AuthOa2token, error)
		LookupAuthOa2tokenByCode(ctx context.Context, code string) (*types.AuthOa2token, error)
		LookupAuthOa2tokenByAccess(ctx context.Context, access string) (*types.AuthOa2token, error)
		LookupAuthOa2tokenByRefresh(ctx context.Context, refresh string) (*types.AuthOa2token, error)

		CreateAuthOa2token(ctx context.Context, rr ...*types.AuthOa2token) error

		DeleteAuthOa2token(ctx context.Context, rr ...*types.AuthOa2token) error
		DeleteAuthOa2tokenByID(ctx context.Context, ID uint64) error

		TruncateAuthOa2tokens(ctx context.Context) error

		// Additional custom functions

		// DeleteExpiredAuthOA2Tokens (custom function)
		DeleteExpiredAuthOA2Tokens(ctx context.Context) error

		// DeleteAuthOA2TokenByCode (custom function)
		DeleteAuthOA2TokenByCode(ctx context.Context, _code string) error

		// DeleteAuthOA2TokenByAccess (custom function)
		DeleteAuthOA2TokenByAccess(ctx context.Context, _access string) error

		// DeleteAuthOA2TokenByRefresh (custom function)
		DeleteAuthOA2TokenByRefresh(ctx context.Context, _refresh string) error

		// DeleteAuthOA2TokenByUserID (custom function)
		DeleteAuthOA2TokenByUserID(ctx context.Context, _userID uint64) error
	}
)

var _ *types.AuthOa2token
var _ context.Context

// SearchAuthOa2tokens returns all matching AuthOa2tokens from store
func SearchAuthOa2tokens(ctx context.Context, s AuthOa2tokens, f types.AuthOa2tokenFilter) (types.AuthOa2tokenSet, types.AuthOa2tokenFilter, error) {
	return s.SearchAuthOa2tokens(ctx, f)
}

// LookupAuthOa2tokenByID
func LookupAuthOa2tokenByID(ctx context.Context, s AuthOa2tokens, id uint64) (*types.AuthOa2token, error) {
	return s.LookupAuthOa2tokenByID(ctx, id)
}

// LookupAuthOa2tokenByCode
func LookupAuthOa2tokenByCode(ctx context.Context, s AuthOa2tokens, code string) (*types.AuthOa2token, error) {
	return s.LookupAuthOa2tokenByCode(ctx, code)
}

// LookupAuthOa2tokenByAccess
func LookupAuthOa2tokenByAccess(ctx context.Context, s AuthOa2tokens, access string) (*types.AuthOa2token, error) {
	return s.LookupAuthOa2tokenByAccess(ctx, access)
}

// LookupAuthOa2tokenByRefresh
func LookupAuthOa2tokenByRefresh(ctx context.Context, s AuthOa2tokens, refresh string) (*types.AuthOa2token, error) {
	return s.LookupAuthOa2tokenByRefresh(ctx, refresh)
}

// CreateAuthOa2token creates one or more AuthOa2tokens in store
func CreateAuthOa2token(ctx context.Context, s AuthOa2tokens, rr ...*types.AuthOa2token) error {
	return s.CreateAuthOa2token(ctx, rr...)
}

// DeleteAuthOa2token Deletes one or more AuthOa2tokens from store
func DeleteAuthOa2token(ctx context.Context, s AuthOa2tokens, rr ...*types.AuthOa2token) error {
	return s.DeleteAuthOa2token(ctx, rr...)
}

// DeleteAuthOa2tokenByID Deletes AuthOa2token from store
func DeleteAuthOa2tokenByID(ctx context.Context, s AuthOa2tokens, ID uint64) error {
	return s.DeleteAuthOa2tokenByID(ctx, ID)
}

// TruncateAuthOa2tokens Deletes all AuthOa2tokens from store
func TruncateAuthOa2tokens(ctx context.Context, s AuthOa2tokens) error {
	return s.TruncateAuthOa2tokens(ctx)
}

func DeleteExpiredAuthOA2Tokens(ctx context.Context, s AuthOa2tokens) error {
	return s.DeleteExpiredAuthOA2Tokens(ctx)
}

func DeleteAuthOA2TokenByCode(ctx context.Context, s AuthOa2tokens, _code string) error {
	return s.DeleteAuthOA2TokenByCode(ctx, _code)
}

func DeleteAuthOA2TokenByAccess(ctx context.Context, s AuthOa2tokens, _access string) error {
	return s.DeleteAuthOA2TokenByAccess(ctx, _access)
}

func DeleteAuthOA2TokenByRefresh(ctx context.Context, s AuthOa2tokens, _refresh string) error {
	return s.DeleteAuthOA2TokenByRefresh(ctx, _refresh)
}

func DeleteAuthOA2TokenByUserID(ctx context.Context, s AuthOa2tokens, _userID uint64) error {
	return s.DeleteAuthOA2TokenByUserID(ctx, _userID)
}
