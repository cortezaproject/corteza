package store

import (
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
)

type (
	composeRecordAux struct {
		refMod string
		relMod *types.Module

		refNs    string
		relUsers resource.UserstampIndex
		walker   resource.CrsWalker
	}

	composeRecord struct {
		cfg *EncoderConfig

		res *resource.ComposeRecord
		rec *composeRecordAux

		relNS  *types.Namespace
		relMod *types.Module

		fieldModRef map[string]resource.Identifiers
		// module identifier -> record identifier -> recordID
		externalRef map[string]map[string]uint64
		recMap      map[string]*types.Record

		// Little helper flag for conditional encoding
		missing bool
	}
)
