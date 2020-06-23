package encoder

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	parsedTime struct {
		field string
		value string
	}
)

// Time formatter
//
// Takes ptr to time.Time so we can conver both cases (value + ptr)
// The function also generates additional fields with included timezone.
func fmtTime(field string, tp *time.Time, tz string) (pt []*parsedTime, err error) {
	if tp == nil {
		return pt, nil
	}

	pt = append(pt, &parsedTime{
		field: field,
		value: tp.UTC().Format(time.RFC3339),
	})
	if tz == "" || tz == "UTC" {
		return
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return pt, err
	}
	tt := tp.In(loc)
	pt = append(pt,
		&parsedTime{
			field: field + "_date",
			value: tt.Format("2006-01-02"),
		}, &parsedTime{
			field: field + "_time",
			value: tt.Format("15:04:05"),
		})
	return
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

	switch true {
	case su.Name != "":
		return su.Name, nil
	case su.Handle != "":
		return su.Handle, nil
	case su.Email != "":
		return su.Email, nil
	case su.Username != "":
		return su.Username, nil
	}

	return fmt.Sprintf("%d", su.ID), nil
}

func (enc flatWriter) Record(r *types.Record) (err error) {
	var out = make([]string, len(enc.ff))

	procTime := func(d []string, pts []*parsedTime, base int) {
		for i, p := range pts {
			d[base+i] = p.value
		}
	}

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
			tt, err := fmtTime("createdAt", &r.CreatedAt, enc.tz)
			if err != nil {
				return err
			}
			procTime(out, tt, f)
		case "updatedBy":
			out[f], err = fmtSysUser(r.UpdatedBy, enc.u)
			if err != nil {
				return err
			}
		case "updatedAt":
			tt, err := fmtTime("updatedAt", r.UpdatedAt, enc.tz)
			if err != nil {
				return err
			}
			procTime(out, tt, f)
		case "deletedBy":
			out[f], err = fmtSysUser(r.DeletedBy, enc.u)
			if err != nil {
				return err
			}
		case "deletedAt":
			tt, err := fmtTime("deletedAt", r.DeletedAt, enc.tz)
			if err != nil {
				return err
			}
			procTime(out, tt, f)
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

	procTime := func(d map[string]interface{}, pts []*parsedTime) {
		for _, p := range pts {
			d[p.field] = p.value
		}
	}

	for _, f := range enc.ff {
		switch f.name {
		case "recordID", "ID":
			// encode all numbers as string (to prevent overflow of uint64 values)
			out[f.name] = fmtUint64(r.ID)
		case "moduleID":
			// encode all numbers as string (to prevent overflow of uint64 values)
			out[f.name] = fmtUint64(r.ModuleID)
		case "namespaceID":
			// encode all numbers as string (to prevent overflow of uint64 values)
			out[f.name] = fmtUint64(r.NamespaceID)
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
			tt, err := fmtTime("createdAt", &r.CreatedAt, enc.tz)
			if err != nil {
				return err
			}
			procTime(out, tt)
		case "updatedBy":
			out[f.name], err = fmtSysUser(r.UpdatedBy, enc.u)
			if err != nil {
				return err
			}
		case "updatedAt":
			if r.UpdatedAt == nil {
				out[f.name] = nil
			} else {
				tt, err := fmtTime("updatedAt", r.UpdatedAt, enc.tz)
				if err != nil {
					return err
				}
				procTime(out, tt)
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
				tt, err := fmtTime("deletedAt", r.DeletedAt, enc.tz)
				if err != nil {
					return err
				}
				procTime(out, tt)
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

	procTime := func(pts []*parsedTime, base int) {
		for i, p := range pts {
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(base+i), p.value)
		}
	}

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
			tt, err := fmtTime("createdAt", &r.CreatedAt, enc.tz)
			if err != nil {
				return err
			}
			procTime(tt, p)
		case "updatedBy":
			u, err = fmtSysUser(r.UpdatedBy, enc.u)
			if err != nil {
				return err
			}
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), u)
		case "updatedAt":
			tt, err := fmtTime("updatedAt", r.UpdatedAt, enc.tz)
			if err != nil {
				return err
			}
			procTime(tt, p)
		case "deletedBy":
			u, err = fmtSysUser(r.DeletedBy, enc.u)
			if err != nil {
				return err
			}
			_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), u)
		case "deletedAt":
			tt, err := fmtTime("deletedAt", r.DeletedAt, enc.tz)
			if err != nil {
				return err
			}
			procTime(tt, p)
		default:
			vv := r.Values.FilterByName(f.name)
			if len(vv) > 0 {
				_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p), vv[0].Value)
			}
		}
	}

	return nil
}
