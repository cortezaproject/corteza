package auth

//import (
//	"net/http"
//)
//
//func MiddlewareValidOnly(next http.Handler) http.Handler {
//	return AccessTokenCheck("api")(next)
//}

//func AccessTokenCheck(s store.AuthOa2tokens, scope ...string) func(http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			if err := validateContextToken(r.Context(), s, scope); err != nil {
//				errors.ProperlyServeHTTP(w, r, err, false)
//				return
//			}
//
//			next.ServeHTTP(w, r)
//		})
//	}
//}
//
//func validateContextToken(ctx context.Context, s store.AuthOa2tokens, scope []string) (err error) {
//	var (
//		token jwt.Token
//	)
//
//	if token, _, err = jwtauth.FromContext(ctx); err != nil {
//		return ErrUnauthorized()
//	}
//
//	if !CheckJwtScope(token, scope...) {
//		return ErrUnauthorizedScope()
//	}
//
//	// Extract the JWT id from the token (string) and convert it to uint64
//	// to be compatible with the lookup function
//	if len(token.JwtID()) < 10 {
//		return ErrMalformedToken("missing or malformed JWT ID")
//	}
//
//	// check if token exists in our DB
//	// there is no need to check for anything beyond existence
//	// because
//	//
//	// @todo we could use a simple caching mechanism here
//	//       1. if lookup is successful, add a JWT ID to the list
//	//       2. add short exp time (that should not last onger than token's exp time)
//	//       3. check against the list first; if JWT ID is not present there check in storage
//	//
//	if _, err = store.LookupAuthOa2tokenByAccess(ctx, s, token.JwtID()); err != nil {
//		return ErrUnauthorized()
//	}
//
//	return nil
//}
