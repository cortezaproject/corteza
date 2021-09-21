package auth

import (
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-chi/jwtauth"
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/errors"
)

func MiddlewareValidOnly(next http.Handler) http.Handler {
	return AccessTokenCheck("api")(next)
}

func AccessTokenCheck(scope ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				ctx   = r.Context()
				roles []uint64
			)

			// retrieve token and claims from context
			tkn, claims, err := jwtauth.FromContext(ctx)
			if err != nil || !tkn.Valid {
				errors.ProperlyServeHTTP(w, r, ErrUnauthorized(), false)
				return
			}

			i := ClaimsToIdentity(claims)
			if !i.Valid() {
				errors.ProperlyServeHTTP(w, r, ErrUnauthorized(), false)
				return
			}

			// check valid scope
			for _, s := range scope {
				if !CheckScope(ctx.Value(scopeCtxKey{}), s) {
					errors.ProperlyServeHTTP(w, r, ErrUnauthorizedScope(), false)
					return
				}
			}

			// verify JWT from store
			_, err = DefaultJwtStore.LookupAuthOa2tokenByAccess(ctx, tkn.Raw)
			if err != nil {
				errors.ProperlyServeHTTP(w, r, ErrUnauthorized(), false)
				return
			}

			u, err := DefaultJwtStore.LookupUserByID(ctx, i.Identity())
			if err != nil {
				errors.ProperlyServeHTTP(w, r, ErrUnauthorized(), false)
				return
			}

			deleteTokens := func() {
				_ = DefaultJwtStore.DeleteAuthOA2TokenByUserID(ctx, u.ID)
			}

			// check if user is not suspended or deleted otherwise remove their all tokens
			if u.SuspendedAt != nil || u.DeletedAt != nil {
				deleteTokens()
				errors.ProperlyServeHTTP(w, r, ErrUnauthorized(), false)
				return
			}

			// check if user's role haven't changed otherwise remove their all tokens
			set, _, _ := DefaultJwtStore.SearchRoleMembers(ctx, types.RoleMemberFilter{UserID: u.ID})
			_ = set.Walk(func(member *types.RoleMember) error {
				roles = append(roles, member.RoleID)
				return nil
			})
			if !equal(roles, i.memberOf) {
				deleteTokens()
				errors.ProperlyServeHTTP(w, r, ErrUnauthorized(), false)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
// fixme maybe move to utils
func equal(a, b []uint64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
