package decoder

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	RecordCreator func(mod *types.Record) error
)

func fmtTime(tp string) (time.Time, error) {
	return time.Parse(time.RFC3339, tp)
}
func fmtTimePtr(tp string) (*time.Time, error) {
	t, err := fmtTime(tp)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func mapify(header, values []string) map[string]string {
	if len(header) != len(values) {
		return nil
	}

	rtr := make(map[string]string)
	for i, v := range values {
		rtr[header[i]] = v
	}

	return rtr
}

func setSystemField(r *types.Record, name, value string) (is bool, err error) {
	switch name {
	case "recordID", "ID":
		r.ID, err = strconv.ParseUint(value, 10, 64)
	case "moduleID":
		r.ModuleID, err = strconv.ParseUint(value, 10, 64)
	case "namespaceID":
		r.NamespaceID, err = strconv.ParseUint(value, 10, 64)
	case "ownedBy":
		r.OwnedBy, err = strconv.ParseUint(value, 10, 64)
	case "createdBy":
		r.CreatedBy, err = strconv.ParseUint(value, 10, 64)
	case "createdAt":
		r.CreatedAt, err = fmtTime(value)
	case "updatedBy":
		r.UpdatedBy, err = strconv.ParseUint(value, 10, 64)
	case "updatedAt":
		r.UpdatedAt, err = fmtTimePtr(value)
	case "deletedBy":
		r.DeletedBy, err = strconv.ParseUint(value, 10, 64)
	case "deletedAt":
		r.DeletedAt, err = fmtTimePtr(value)
	default:
		return false, err
	}
	return true, err
}

func (dec flatReader) Records(fields map[string]string, Create RecordCreator) error {
	header := dec.Header()

	err := dec.walk(func(row []string) error {
		mapped := mapify(header, row)
		r := types.Record{}
		rvs := types.RecordValueSet{}

		i := 0
		for imp, rec := range fields {
			if rec == "" {
				return errors.New("Can not import record: Record field not defined")
			}

			val := mapped[imp]
			if system, err := setSystemField(&r, rec, val); err != nil {
				return err
			} else if !system {
				rv := types.RecordValue{
					Name:  rec,
					Value: val,
					Place: uint(i),
				}
				i++

				rvs = append(rvs, &rv)
			}
		}

		r.Values = rvs
		return Create(&r)
	})

	return err
}

func (dec structuredDecoder) Records(fields map[string]string, Create RecordCreator) error {
	err := dec.walk(func(entry map[string]interface{}) error {
		r := types.Record{}
		rvs := types.RecordValueSet{}

		i := 0
		for imp, rec := range fields {
			if rec == "" {
				return errors.New("Can not import record: Record field not defined")
			}

			val := fmt.Sprintf("%v", entry[imp])
			if system, err := setSystemField(&r, rec, val); err != nil {
				return err
			} else if !system {
				rv := types.RecordValue{
					Name:  rec,
					Value: val,
					Place: uint(i),
				}
				i++

				rvs = append(rvs, &rv)
			}
		}

		r.Values = rvs
		return Create(&r)
	})

	return err
}
