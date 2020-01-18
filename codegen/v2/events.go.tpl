package {{ .Package }}

import "github.com/cortezaproject/corteza-server/pkg/eventbus"

// Match returns false if given conditions do not match event & resource internals
func (res {{ camelCase .ResourceIdent "base" }}) Match(c constraint) bool {
	// By default we match no mather what kind of constraints we receive
	//
	// Function will be called multiple times - once for every trigger constraint
	// All should match (return true):
	//   constraint#1 AND constraint#2 AND constraint#3 ...
	//
	// When there are multiple values, Match() can decide how to treat them (OR, AND...)
	return true
}
