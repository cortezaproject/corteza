package service

import (
	"context"

	ct "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/decoder"
	st "github.com/cortezaproject/corteza-server/system/types"
)

type (
	dataProcesser struct {
		ID                  uint64
		ComposeModuleID     uint64
		ComposeNamespaceID  uint64
		NodeBaseURL         string
		ModuleMappings      *types.ModuleFieldMappingSet
		ModuleMappingValues *ct.RecordValueSet
		SyncService         *Sync
		Node                *types.Node
		User                *st.User
	}

	dataProcesserResponse struct {
		Processed int
	}
)

// Process gets the payload from syncer and
// uses the decode package to decode the whole set, depending on
// the filtering that was used (limit)
func (dp *dataProcesser) Process(ctx context.Context, payload []byte) (ProcesserResponse, error) {
	processed := 0
	o, err := decoder.DecodeFederationRecordSync([]byte(payload))

	if err != nil {
		return dataProcesserResponse{
			Processed: processed,
		}, err
	}

	if len(o) == 0 {
		return dataProcesserResponse{
			Processed: processed,
		}, nil
	}

	// get the user that is tied to this node
	ctx = auth.SetIdentityToContext(ctx, dp.User)

	for _, er := range o {
		dp.SyncService.mapper.Merge(&er.Values, dp.ModuleMappingValues, dp.ModuleMappings)

		rec := &ct.Record{
			ModuleID:    dp.ComposeModuleID,
			NamespaceID: dp.ComposeNamespaceID,
			Values:      *dp.ModuleMappingValues,
		}

		AddFederationLabel(rec, dp.NodeBaseURL)

		_, err := dp.SyncService.CreateRecord(ctx, rec)

		if err != nil {
			continue
		}

		processed++
	}

	return dataProcesserResponse{
		Processed: processed,
	}, nil
}
