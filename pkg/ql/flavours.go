package ql

type (
	Encoder interface {
		// CaseInsensitiveLike translates "like" or "not like" op to case-insensitive version
		// supported by the db
		CaseInsensitiveLike(neg bool) string
	}
)

var (
	// QueryEncoder
	// This is a temp solution to enable multi-db support in the ql package
	// that is not aware of the SQL flavour
	//
	// We can get away with this kind of solution right now because we are currently
	// only supporting one single db at the time.
	//
	// This will change in the future and so will the pkg/ql logic
	QueryEncoder Encoder
)
