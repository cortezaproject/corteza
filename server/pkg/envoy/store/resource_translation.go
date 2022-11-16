package store

import (
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	resourceTranslation struct {
		cfg *EncoderConfig

		locales types.ResourceTranslationSet

		// point to the resource translation
		refResourceTranslation string
		refLocaleRes           *resource.Ref

		refPathRes []*resource.Ref
	}
)

func resourceTranslationErrUnidentifiable(ii resource.Identifiers) error {
	return fmt.Errorf("resource translation unidentifiable %v", ii.StringSlice())
}
