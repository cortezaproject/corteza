package dal

import "regexp"

type (
	ConnectionWrap struct {
		ID uint64

		connection Connection
		params     ConnectionParams
		Config     ConnectionConfig
	}

	ConnectionConfig struct {
		SensitivityLevelID uint64
		Label              string

		// When model does not specify the ident (table name for example), fallback to this
		ModelIdent string

		// when a new model is added on a connection, it's ident
		// is verified against this regexp
		//
		// ident is considered valid if it matches one of the expressions
		// or if the list of checks is empty
		ModelIdentCheck []*regexp.Regexp
	}
)

func checkIdent(ident string, rr ...*regexp.Regexp) bool {
	if len(rr) == 0 {
		return true
	}

	for _, r := range rr {
		if r.MatchString(ident) {
			return true
		}
	}

	return false
}
