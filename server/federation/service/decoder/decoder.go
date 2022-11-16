package decoder

import (
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza/server/compose/types"
	ftypes "github.com/cortezaproject/corteza/server/federation/types"
)

type (
	ExposedRecord struct {
		ID     uint64               `json:"recordID,string"`
		Values types.RecordValueSet `json:"values"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	// Bits for decoding structure sync
	ModuleDocument struct {
		Response ComposeModule
	}

	ExposedModuleDocument struct {
		Response struct {
			Set ftypes.ExposedModuleSet
		}
	}

	ExposedRecordDocument struct {
		Response struct {
			Set []*ExposedRecord
		}
	}

	ExposedModule struct {
		ftypes.ExposedModule
	}

	// Bits for decoding data sync
	ComposeRecordResponse struct {
		Filter ComposeRecordFilter `json:",omitempty"`
		Set    ComposeRecordSet    `json:",omitempty"`
	}
	ComposeRecordDocument struct {
		Response ComposeRecordResponse
	}
)

// DecodeModuleSync decodes the response from the structure sync request.
//
// It returns a set of modules.
func DecodeModuleSync(in []byte) (types.ModuleSet, error) {
	out := &ModuleDocument{}

	err := json.Unmarshal(in, out)
	if err != nil {
		return nil, err
	}

	return types.ModuleSet{&out.Response.Module}, nil
}

func DecodeFederationModuleSync(in []byte) (ftypes.ExposedModuleSet, error) {
	out := &ExposedModuleDocument{}

	err := json.Unmarshal(in, out)
	if err != nil {
		return nil, err
	}

	return out.Response.Set, nil
}

func DecodeFederationRecordSync(in []byte) ([]*ExposedRecord, error) {
	out := &ExposedRecordDocument{}

	err := json.Unmarshal(in, out)
	if err != nil {
		return nil, err
	}

	return out.Response.Set, nil
}

// DecodeRecordSync decodes the response from the data sync request.
//
// It returns a set of records along with the used filter.
func DecodeRecordSync(in []byte) (types.RecordSet, *types.RecordFilter, error) {
	out := &ComposeRecordDocument{}

	err := json.Unmarshal(in, out)
	if err != nil {
		return nil, nil, err
	}

	ss := make(types.RecordSet, 0, len(out.Response.Set))
	for _, r := range out.Response.Set {
		ss = append(ss, &r.Record)
	}

	return ss, &out.Response.Filter.RecordFilter, nil
}
