package encoder

import (
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
)

// Time formatter
//
// Takes ptr to time.Time so we can conver both cases (value + ptr)
func fmtTime(tp *time.Time) string {
	if tp == nil {
		return ""
	}

	return tp.UTC().Format(time.RFC3339)
}

func fmtUint64(u uint64) string {
	return strconv.FormatUint(u, 10)
}

func (enc flatWriter) Record(r *types.Record) error {
	var out = make([]string, len(enc.ff))

	for f, field := range enc.ff {
		switch field.name {
		case "recordID", "ID":
			out[f] = fmtUint64(r.ID)
		case "moduleID":
			out[f] = fmtUint64(r.ModuleID)
		case "namespaceID":
			out[f] = fmtUint64(r.NamespaceID)
		case "ownedBy":
			out[f] = fmtUint64(r.OwnedBy)
		case "createdBy":
			out[f] = fmtUint64(r.CreatedBy)
		case "createdAt":
			out[f] = fmtTime(&r.CreatedAt)
		case "updatedBy":
			out[f] = fmtUint64(r.UpdatedBy)
		case "updatedAt":
			out[f] = fmtTime(r.UpdatedAt)
		case "deletedBy":
			out[f] = fmtUint64(r.DeletedBy)
		case "deletedAt":
			out[f] = fmtTime(r.DeletedAt)
		default:
			vv := r.Values.FilterByName(field.name)
			// @todo support for field.encodeAllMulti
			if len(vv) > 0 {
				out[f] = vv[0].Value
			}
		}
	}

	defer enc.w.Flush()

	return enc.w.Write(out)
}

func (enc structuredEncoder) Record(r *types.Record) error {
	var (
		// Exporter can choose fields so we need this buffer
		// to hold just what we need
		out = make(map[string]interface{})
		vv  types.RecordValueSet
		c   int
	)

	for _, f := range enc.ff {
		switch f.name {
		case "recordID", "ID":
			out[f.name] = r.ID
		case "moduleID":
			out[f.name] = r.ModuleID
		case "namespaceID":
			out[f.name] = r.NamespaceID
		case "ownedBy":
			out[f.name] = r.OwnedBy
		case "createdBy":
			out[f.name] = r.CreatedBy
		case "createdAt":
			out[f.name] = fmtTime(&r.CreatedAt)
		case "updatedBy":
			out[f.name] = r.UpdatedBy
		case "updatedAt":
			if r.UpdatedAt == nil {
				out[f.name] = nil
			} else {
				out[f.name] = fmtTime(r.UpdatedAt)
			}

		case "deletedBy":
			out[f.name] = r.DeletedBy
		case "deletedAt":
			if r.DeletedAt == nil {
				out[f.name] = nil
			} else {
				out[f.name] = fmtTime(r.DeletedAt)
			}

		default:
			vv = r.Values.FilterByName(f.name)
			c = len(vv)

			if c == 0 {
				break
			}

			if c == 1 {
				out[f.name] = vv[0].Value
			} else {
				multi := make([]string, c)

				for n := range vv {
					multi[n] = vv[n].Value
				}

				out[f.name] = multi
			}
		}
	}

	return enc.w.Encode(out)
}

func (enc *excelizeEncoder) Record(r *types.Record) error {
	enc.row++

	for p, f := range enc.ff {
		p++
		switch f.name {
		case "recordID", "ID":
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), fmtUint64(r.ID))
		case "moduleID":
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), fmtUint64(r.ModuleID))
		case "namespaceID":
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), fmtUint64(r.NamespaceID))
		case "ownedBy":
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), fmtUint64(r.OwnedBy))
		case "createdBy":
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), fmtUint64(r.CreatedBy))
		case "createdAt":
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), fmtTime(&r.CreatedAt))
		case "updatedBy":
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), fmtUint64(r.UpdatedBy))
		case "updatedAt":
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), fmtTime(r.UpdatedAt))
		case "deletedBy":
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), fmtUint64(r.DeletedBy))
		case "deletedAt":
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), fmtTime(r.DeletedAt))
		default:
			vv := r.Values.FilterByName(f.name)
			if len(vv) > 0 {
				_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), vv[0].Value)
			}
		}
	}

	return nil
}
