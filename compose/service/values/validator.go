package values

import (
	"errors"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/slice"
	"math/big"
	"net/mail"
	"net/url"
	"strings"
	"time"
)

// Validator package provides tooling to validate
// record and it's values against field constraints
//
// Structures and basic logic is similar to what we offer on the frontend
// (see corteza-js/validator package) but with less features as there
// is no need for such level of interaction and dynamic we require on the frontend

type (
	UniqueChecker    func(*types.RecordValue, *types.ModuleField, *types.Module) (uint64, error)
	ReferenceChecker func(*types.RecordValue, *types.ModuleField, *types.Module) (bool, error)

	validator struct {
		uniqueCheckerFn    UniqueChecker
		recordRefCheckerFn ReferenceChecker
		userRefCheckerFn   ReferenceChecker
		fileRefCheckerFn   ReferenceChecker

		now func() time.Time
	}
)

func makeInternalErr(field *types.ModuleField, err error) types.RecordValueError {
	return types.RecordValueError{Kind: "internal", Message: err.Error(), Meta: map[string]interface{}{"field": field.Name}}
}

func makeEmptyErr(field *types.ModuleField) types.RecordValueError {
	return types.RecordValueError{Kind: "empty", Meta: map[string]interface{}{"field": field.Name}}
}

func makeInvalidValueErr(field *types.ModuleField, value string) types.RecordValueError {
	return types.RecordValueError{Kind: "invalidValue", Meta: map[string]interface{}{"field": field.Name, "value": value}}
}

func makeInvalidRefErr(field *types.ModuleField, ref uint64) types.RecordValueError {
	return types.RecordValueError{Kind: "invalidRef", Meta: map[string]interface{}{"field": field.Name, "ref": ref}}
}

func makeDuplicateValueInSetErr(field *types.ModuleField, value string) types.RecordValueError {
	return types.RecordValueError{Kind: "duplicateValueInSet", Meta: map[string]interface{}{"field": field.Name, "value": value}}
}

func makeDuplicateValueErr(field *types.ModuleField, recordID uint64) types.RecordValueError {
	return types.RecordValueError{Kind: "duplicateValue", Meta: map[string]interface{}{"field": field.Name, "recordID": recordID}}
}

// Simple wrapper for easier error returning from validation functions
func e2s(ee ...types.RecordValueError) []types.RecordValueError {
	return ee
}

func Validator() *validator {
	return &validator{
		now: func() time.Time { return time.Now() },
	}
}

func (vldtr *validator) UniqueChecker(fn UniqueChecker) {
	vldtr.uniqueCheckerFn = fn
}

func (vldtr *validator) RecordRefChecker(fn ReferenceChecker) {
	vldtr.recordRefCheckerFn = fn
}

func (vldtr *validator) UserRefChecker(fn ReferenceChecker) {
	vldtr.userRefCheckerFn = fn
}

func (vldtr *validator) FileRefChecker(fn ReferenceChecker) {
	vldtr.fileRefCheckerFn = fn
}

// Run validates record and it's values against module & module fields options
//
//
// Validation is done in phases for optimal resource usage:
//   - check if required values are present
//   - check for unique-multi-value in multi value fields
//   - field-kind specific validation on all values
//   - unique check on all all values
func (vldtr validator) Run(m *types.Module, r *types.Record) (out *types.RecordValueErrorSet) {
	var (
		f *types.ModuleField
	)

	out = &types.RecordValueErrorSet{}

fields:
	for _, f := range m.Fields {
		vv := r.Values.FilterByName(f.Name)

		if f.Required {
			if len(vv) == 0 {
				out.Push(makeEmptyErr(f))
				continue fields
			}

			for _, v := range vv {
				if len(v.Value) == 0 || (f.IsRef() && v.Ref == 0) {
					out.Push(makeEmptyErr(f))
					continue fields
				}
			}
		}

		if f.Multi && f.Options.IsUniqueMultiValue() {
			flipped := make(map[string]bool)
			for _, v := range vv {
				if flipped[v.Value] {
					out.Push(makeDuplicateValueInSetErr(f, v.Value))
					continue fields
				}

				flipped[v.Value] = true
			}
		}
	}

	for _, v := range r.Values {
		if !v.IsUpdated() || v.IsUpdated() {
			// We'll validate only updated (and non-deleted) values
			continue
		}

		if f = m.Fields.FindByName(v.Name); f == nil {
			continue
		}

		if v.Value == "" {
			// Nothing to do with empty value
			return nil
		}

		// Per field type validators
		switch strings.ToLower(f.Kind) {
		case "bool":
			out.Push(vldtr.vBool(v, f, r, m)...)
		case "datetime":
			out.Push(vldtr.vDatetime(v, f, r, m)...)
		case "email":
			out.Push(vldtr.vEmail(v, f, r, m)...)
		case "file":
			out.Push(vldtr.vFile(v, f, r, m)...)
		case "number":
			out.Push(vldtr.vNumber(v, f, r, m)...)
		case "record":
			out.Push(vldtr.vRecord(v, f, r, m)...)
		case "select":
			out.Push(vldtr.vSelect(v, f, r, m)...)
		//case "string":
		//	out.Push(vldtr.vString(v, f, r, m)...)
		case "url":
			out.Push(vldtr.vUrl(v, f, r, m)...)
		case "user":
			out.Push(vldtr.vUser(v, f, r, m)...)
		}
	}

	// This is the most resource-heavy operation
	// we'll do in at the end
	for _, v := range r.Values {
		if f = m.Fields.FindByName(v.Name); f == nil {
			continue
		}

		if !f.Options.IsUnique() {
			// Only interested in unique fields
			continue
		}

		duplicateRecordID, err := vldtr.uniqueCheckerFn(v, f, m)
		if err != nil {
			out.Push(makeInternalErr(f, err))
		} else if duplicateRecordID > 0 && duplicateRecordID != r.ID {
			out.Push(makeDuplicateValueErr(f, duplicateRecordID))
		}
	}

	if out.IsValid() {
		return nil
	}

	return out
}

func (vldtr validator) vBool(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	if v.Value == "" {
		return nil
	}

	if v.Value != strBoolTrue && v.Value != strBoolFalse {
		return e2s(makeInvalidValueErr(f, v.Value))
	}

	return nil
}

func (vldtr validator) vDatetime(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	var (
		inputFormat string
		t           time.Time
		err         error

		// We'll validate against this value
		refTime = vldtr.now()
	)

	if f.Options.Bool(fieldOpt_Datetime_onlyDate) {
		inputFormat = datetimeInternalFormatDate

		// Round down ref time to midnight
		refTime = time.Date(refTime.Year(), refTime.Month(), refTime.Day(), 0, 0, 0, 0, refTime.Location())
	} else if f.Options.Bool(fieldOpt_Datetime_onlyTime) {
		inputFormat = datetimeIntenralFormatTime

		// Round down ref time to day one
		refTime = time.Date(0, 1, 1, refTime.Hour(), refTime.Minute(), refTime.Second(), refTime.Nanosecond(), refTime.Location())
	} else {
		inputFormat = datetimeInternalFormatFull
	}

	t, err = time.Parse(inputFormat, v.Value)
	if err != nil {
		return e2s(makeInvalidValueErr(f, v.Value))
	}

	if f.Options.Bool(fieldOpt_Datetime_onlyFutureValues) {
		if !t.After(refTime) {
			return e2s(makeInvalidValueErr(f, v.Value))
		}
	} else if f.Options.Bool(fieldOpt_Datetime_onlyPastValues) {
		if !t.Before(refTime) {
			return e2s(makeInvalidValueErr(f, v.Value))
		}
	}

	// @todo check past/future

	return nil
}

func (vldtr validator) vEmail(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	if _, err := mail.ParseAddress(v.Value); err != nil {
		return e2s(makeInvalidValueErr(f, v.Value))
	}

	return nil
}

func (vldtr validator) vFile(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	if ok, err := vldtr.fileRefCheckerFn(v, f, m); err != nil {
		return e2s(makeInternalErr(f, err))
	} else if !ok {
		return e2s(makeInvalidRefErr(f, v.Ref))
	}

	return nil
}

func (vldtr validator) vNumber(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	var (
		precision = uint(f.Options.Int64Def(fieldOpt_Number_precision, 2))
	)

	if _, _, err := big.ParseFloat(v.Value, 0, precision, big.ToNearestEven); err != nil {
		return e2s(makeInvalidValueErr(f, v.Value))
	}

	return nil
}

func (vldtr validator) vRecord(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	if ok, err := vldtr.recordRefCheckerFn(v, f, m); err != nil {
		return e2s(makeInternalErr(f, err))
	} else if !ok {
		return e2s(makeInvalidRefErr(f, v.Ref))
	}

	return nil
}

func (vldtr validator) vSelect(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	// @todo validate v.Value against
	var (
		options = f.Options["options"]
		sbm     = make(map[string]bool)
	)

	if oo, ok := options.([]string); ok {
		sbm = slice.ToStringBoolMap(oo)
	} else if oo, ok := options.([]map[string]string); ok {
		for _, o := range oo {
			sbm[o["value"]] = true
		}
	} else if options == nil {
		return nil
	} else {
		return e2s(makeInternalErr(f, errors.New("invalid select options definition")))
	}

	if !sbm[v.Value] {
		return e2s(makeInvalidValueErr(f, v.Value))

	}

	return nil
}

//func (vldtr validator) vString(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
//	return nil
//}

func (vldtr validator) vUrl(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	if p, err := url.Parse(v.Value); err != nil {
		return e2s(makeInvalidValueErr(f, v.Value))
	} else if p.Scheme == "" || p.Host == "" {
		return e2s(makeInvalidValueErr(f, v.Value))
	} else if f.Options.Bool(fieldOpt_Url_onlySecure) && p.Scheme != "https" {
		return e2s(makeInvalidValueErr(f, v.Value))
	}

	return nil
}

func (vldtr validator) vUser(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	if ok, err := vldtr.userRefCheckerFn(v, f, m); err != nil {
		return e2s(makeInternalErr(f, err))
	} else if !ok {
		return e2s(makeInvalidRefErr(f, v.Ref))
	}

	return nil
}
