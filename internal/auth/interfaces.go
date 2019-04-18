package auth

type (
	Identifiable interface {
		Identity() uint64
		Valid() bool
	}

	TokenEncoder interface {
		Encode(identity Identifiable) string
	}
)
