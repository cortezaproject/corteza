package agentic

import (
	"strings"
	"testing"

	cmpTypes "github.com/crusttech/human/server/compose/types"
)

func TestFieldsParamDocNamesReservedFieldNames(t *testing.T) {
	for _, name := range cmpTypes.ModuleReservedFieldNames {
		if !strings.Contains(fieldsParamDoc, name) {
			t.Errorf("fields documentation does not name reserved field %q", name)
		}
	}
}
