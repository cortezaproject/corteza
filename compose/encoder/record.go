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

func fmtSysUser(u uint64, finder userFinder) (string, error) {
	if u <= 0 || finder == nil {
		return fmtUint64(u), nil
	}
	su, err := finder(u)
	if err != nil {
		return "", err
	}
	return su.Name, nil
}

func (enc flatWriter) Record(r *types.Record) (err error) {
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
			out[f], err = fmtSysUser(r.OwnedBy, enc.u)
			if err != nil {
				return err
			}
		case "createdBy":
			out[f], err = fmtSysUser(r.CreatedBy, enc.u)
			if err != nil {
				return err
			}
		case "createdAt":
			out[f] = fmtTime(&r.CreatedAt)
		case "updatedBy":
			out[f], err = fmtSysUser(r.UpdatedBy, enc.u)
			if err != nil {
				return err
			}
		case "updatedAt":
			out[f] = fmtTime(r.UpdatedAt)
		case "deletedBy":
			out[f], err = fmtSysUser(r.DeletedBy, enc.u)
			if err != nil {
				return err
			}
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

func (enc structuredEncoder) Record(r *types.Record) (err error) {
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
			out[f.name], err = fmtSysUser(r.OwnedBy, enc.u)
			if err != nil {
				return err
			}
		case "createdBy":
			out[f.name], err = fmtSysUser(r.CreatedBy, enc.u)
			if err != nil {
				return err
			}
		case "createdAt":
			out[f.name] = fmtTime(&r.CreatedAt)
		case "updatedBy":
			out[f.name], err = fmtSysUser(r.UpdatedBy, enc.u)
			if err != nil {
				return err
			}
		case "updatedAt":
			if r.UpdatedAt == nil {
				out[f.name] = nil
			} else {
				out[f.name] = fmtTime(r.UpdatedAt)
			}

		case "deletedBy":
			out[f.name], err = fmtSysUser(r.DeletedBy, enc.u)
			if err != nil {
				return err
			}
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

func (enc *excelizeEncoder) Record(r *types.Record) (err error) {
	enc.row++
	var u string

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
			u, err = fmtSysUser(r.OwnedBy, enc.u)
			if err != nil {
				return err
			}
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), u)
		case "createdBy":
			u, err = fmtSysUser(r.CreatedBy, enc.u)
			if err != nil {
				return err
			}
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), u)
		case "createdAt":
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), fmtTime(&r.CreatedAt))
		case "updatedBy":
			u, err = fmtSysUser(r.UpdatedBy, enc.u)
			if err != nil {
				return err
			}
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), u)
		case "updatedAt":
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), fmtTime(r.UpdatedAt))
		case "deletedBy":
			u, err = fmtSysUser(r.DeletedBy, enc.u)
			if err != nil {
				return err
			}
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), u)
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
