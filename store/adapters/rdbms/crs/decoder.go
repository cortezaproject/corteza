package crs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/data"
	"github.com/google/uuid"
	"github.com/spf13/cast"
)

// decodes values into record by iterating over all columns
// and running decode handler on all
func decode(columns []*column, values []any, r *types.Record) (err error) {
	for i, c := range columns {
		if err = c.decode(c.attributes, values[i], r); err != nil {
			return
		}
	}
	return nil
}

func decodeStdRecordValueJSON(aa []*data.Attribute, value any, r *types.Record) (err error) {
	rawJson, is := value.(*sql.RawBytes)
	if !is {
		return fmt.Errorf("incompatible input value type (%T), expecting *sql.RawBytes", value)
	}

	buf := make(map[string][]any)
	if err = json.Unmarshal(*rawJson, &buf); err != nil {
		return
	}

	rvs := make([]*types.RecordValue, 0, len(buf))

	// @todo this is too naive
	for f, vv := range buf {
		if isSystemField(f) {
			if err = decodeSystemField(f, value, r); err != nil {
				return
			}
		} else {
			for i, v := range vv {
				rvs = append(rvs, &types.RecordValue{Name: f, Value: cast.ToString(v), Place: uint(i)})
			}
		}
	}

	r.Values = rvs

	return
}

// decode non-serialized record value
func decodeRecordValue(aa []*data.Attribute, raw any, r *types.Record) (err error) {
	if len(aa) == 0 {
		return fmt.Errorf("can not decode value from column without attributes")
	}

	if len(aa) > 1 {
		return fmt.Errorf("can not decode value from a non-embedded column with more than one attribute")
	}

	var (
		rv = &types.RecordValue{Name: aa[0].Ident}
	)

	switch aux := raw.(type) {
	case *uint64:
		if aux != nil {
			rv.Ref = *aux
			rv.Value = strconv.FormatUint(*aux, 10)
		}

	case *sql.NullTime:
		if aux == nil || !aux.Valid {
			return
		}

		rv.Value = aux.Time.Format(time.RFC3339)

	case *sql.NullString:
		if aux == nil || !aux.Valid {
			return
		}

		rv.Value = aux.String

	case *sql.NullBool:
		if aux == nil || !aux.Valid {
			return
		}

		if aux.Bool {
			rv.Value = "true"
		} else {
			rv.Value = "false"
		}

	case *sql.RawBytes:
		if raw == nil {
			return
		}

		rv.Value = string(*aux)

	case *uuid.UUID:
		if raw == nil {
			return
		}

		rv.Value = aux.String()

	default:
		return
	}

	r.Values = r.Values.Set(rv)
	return nil
}

func decodeSystemField(ident string, raw any, r *types.Record) error {
	switch ident {
	case sysID:
		if tmp, is := raw.(*uint64); is && tmp != nil {
			r.ID = *tmp
		} else {
			return fmt.Errorf("incompatible type %T for system field %s", raw, sysID)
		}
	case sysNamespaceID:
		if tmp, is := raw.(*uint64); is && tmp != nil {
			r.ID = *tmp
		} else {
			return fmt.Errorf("incompatible type %T for system field %s", raw, sysNamespaceID)
		}
	case sysModuleID:
		if tmp, is := raw.(*uint64); is && tmp != nil {
			r.ID = *tmp
		} else {
			return fmt.Errorf("incompatible type %T for system field %s", raw, sysModuleID)
		}
	case sysCreatedAt:
		if tmp, is := raw.(*sql.NullTime); is {
			if tmp.Valid {
				r.CreatedAt = tmp.Time
			}
		} else {
			return fmt.Errorf("incompatible type %T for system field %s", raw, sysCreatedAt)
		}

	case sysCreatedBy:
		if tmp, is := raw.(*uint64); is && tmp != nil {
			r.CreatedBy = *tmp
		} else {
			return fmt.Errorf("incompatible type %T for system field %s", raw, sysCreatedBy)
		}

	case sysUpdatedAt:
		if tmp, is := raw.(*sql.NullTime); is {
			if tmp.Valid {
				r.UpdatedAt = &tmp.Time
			}
		} else {
			return fmt.Errorf("incompatible type %T for system field %s", raw, sysUpdatedAt)
		}

	case sysUpdatedBy:
		if tmp, is := raw.(*uint64); is && tmp != nil {
			r.UpdatedBy = *tmp
		} else {
			return fmt.Errorf("incompatible type %T for system field %s", raw, sysUpdatedBy)
		}

	case sysDeletedAt:
		if tmp, is := raw.(*sql.NullTime); is {
			if tmp.Valid {
				r.DeletedAt = &tmp.Time
			}
		} else {
			return fmt.Errorf("incompatible type %T for system field %s", raw, sysDeletedAt)
		}

	case sysDeletedBy:
		if tmp, is := raw.(*uint64); is && tmp != nil {
			r.DeletedBy = *tmp
		} else {
			return fmt.Errorf("incompatible type %T for system field %s", raw, sysDeletedBy)
		}

	case sysOwnedBy:
		if tmp, is := raw.(*uint64); is && tmp != nil {
			r.OwnedBy = *tmp
		} else {
			return fmt.Errorf("incompatible type %T for system field %s", raw, sysOwnedBy)
		}
	}

	return nil
}
