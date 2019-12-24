package corredor

// Authentication token maker must be able to convert user's handle or email
// into valid authentication token with short expiration

// Used by non-system services
func CrossServiceAuthTokenMaker() AuthTokenMaker {
	return func(user string) (s string, err error) {
		panic("not implemented")
		return "", nil
	}
}

// InternalAuthTokenMaker used by system or by all services when running in monolith mode
func InternalAuthTokenMaker() AuthTokenMaker {
	return func(user string) (s string, err error) {
		panic("not implemented")
		// @todo implementation
		//
		// DefaultUser.FindByAny(user)
		// auth.TokenEncoder

		return "", nil
	}
}
