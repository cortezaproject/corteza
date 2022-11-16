package service

import (
	"context"
	"fmt"

	ct "github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/federation/service/decoder"
	"github.com/cortezaproject/corteza/server/federation/types"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	st "github.com/cortezaproject/corteza/server/system/types"
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

	ctx = auth.SetIdentityToContext(ctx, auth.FederationUser())

	for _, er := range o {
		var (
			rec *ct.Record
			err error
		)

		dp.SyncService.mapper.Merge(&er.Values, dp.ModuleMappingValues, dp.ModuleMappings)

		if er.DeletedAt != nil {
			// find the record
			if rec, err = dp.findRecordByFederationID(ctx, er.ID, dp.ComposeModuleID, dp.ComposeNamespaceID); err != nil {
				continue
			}

			// Handle edge cases where the data doesn't exist anymore
			if rec != nil {
				dp.SyncService.DeleteRecord(ctx, rec)
			}
			processed++

			continue
		}

		if er.UpdatedAt != nil {
			if rec, err = dp.findRecordByFederationID(ctx, er.ID, dp.ComposeModuleID, dp.ComposeNamespaceID); err != nil {
				// could not find existing record
				continue
			}

			if rec != nil {
				rec.Values = *dp.ModuleMappingValues
			}
		}

		// if the record was updated on origin, but we somehow do not have it
		// create it anyway
		if rec == nil {
			rec = &ct.Record{
				ModuleID:    dp.ComposeModuleID,
				NamespaceID: dp.ComposeNamespaceID,
				Values:      *dp.ModuleMappingValues,
				Meta: map[string]any{
					"federation":           dp.NodeBaseURL,
					"federation_extrecord": fmt.Sprintf("%d", er.ID),
				},
			}
		}

		if rec.ID != 0 {
			_, err = dp.SyncService.UpdateRecord(ctx, rec)
		} else {
			_, err = dp.SyncService.CreateRecord(ctx, rec)
		}

		if err != nil {
			continue
		}

		processed++
	}

	return dataProcesserResponse{
		Processed: processed,
	}, nil
}

// findRecordByFederationID finds any already existing records via
// federation label
func (dp *dataProcesser) findRecordByFederationID(ctx context.Context, recordID, moduleID, namespaceID uint64) (r *ct.Record, err error) {
	filter := ct.RecordFilter{
		NamespaceID: namespaceID,
		ModuleID:    moduleID,
		Meta:        map[string]any{"federation_extrecord": fmt.Sprintf("%d", recordID)}}

	if s, err := dp.SyncService.FindRecords(ctx, filter); err == nil {
		if len(s) == 1 {
			r = s[0]
		}
	}

	return
}
