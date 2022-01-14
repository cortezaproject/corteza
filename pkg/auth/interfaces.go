package auth

type (
	Identifiable interface {
		Identity() uint64
		Roles() []uint64
		Valid() bool
		String() string
	}

	//TokenGenerator interface {
	//	Encode(i Identifiable, clientID uint64, scope ...string) (token []byte, err error)
	//	Generate(ctx context.Context, i Identifiable, clientID uint64, scope ...string) (token []byte, err error)
	//}

	//TokenHandler interface {
	//	TokenGenerator
	//	HttpVerifier() func(http.Handler) http.Handler
	//	HttpAuthenticator() func(http.Handler) http.Handler
	//}

	Signer interface {
		Sign(userID uint64, pp ...interface{}) string
		Verify(signature string, userID uint64, pp ...interface{}) bool
	}
)
