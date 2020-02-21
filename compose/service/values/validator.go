package values

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"strings"
)

// Validator package provides tooling to validate
// record and it's values against field constraints
//
// Structures and basic logic is similar to what we offer on the frontend
// (see corteza-js/validator package) but with less features as there
// is no need for such level of interaction and dynamic we require on the frontend

type (
	UniqueChecker    func(*types.Module, *types.ModuleField, *types.RecordValue) (uint64, error)
	ReferenceChecker func(*types.Module, *types.ModuleField, *types.RecordValue) (bool, error)

	validator struct {
		uniqueCheckerFn    UniqueChecker
		recordRefCheckerFn ReferenceChecker
		userRefCheckerFn   ReferenceChecker
	}
)

func makeInternalErr(field string, err error) types.RecordValueError {
	return types.RecordValueError{Kind: "internal", Message: err.Error(), Meta: map[string]interface{}{"field": field}}
}
func makeEmptyErr(field string) types.RecordValueError {
	return types.RecordValueError{Kind: "empty", Meta: map[string]interface{}{"field": field}}
}
func makeDuplicateValueInSetErr(field string) types.RecordValueError {
	return types.RecordValueError{Kind: "duplicateValueInSet", Meta: map[string]interface{}{"field": field}}
}
func makeDuplicateValueErr(field string, recordID uint64) types.RecordValueError {
	return types.RecordValueError{Kind: "duplicateValue", Meta: map[string]interface{}{"recordID": recordID, "field": field}}
}

func Validator() *validator {
	return &validator{}
}

func (vldtr validator) UniqueChecker(fn UniqueChecker) {
	vldtr.uniqueCheckerFn = fn
}

func (vldtr validator) RecordRefChecker(fn ReferenceChecker) {
	vldtr.recordRefCheckerFn = fn
}

func (vldtr validator) UserRefChecker(fn ReferenceChecker) {
	vldtr.userRefCheckerFn = fn
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
				out.Push(makeEmptyErr(f.Name))
				continue fields
			}

			for _, v := range vv {
				if len(v.Value) == 0 || (f.IsRef() && v.Ref == 0) {
					out.Push(makeEmptyErr(f.Name))
					continue fields
				}
			}
		}

		if f.Multi && f.Options.IsUniqueMultiValue() {
			flipped := make(map[string]bool)
			for _, v := range vv {
				if flipped[v.Value] {
					out.Push(makeDuplicateValueInSetErr(f.Name))
					continue fields
				}

				flipped[v.Value] = true
			}
		}
	}

	for _, v := range r.Values {
		if f = m.Fields.FindByName(v.Name); f == nil {
			continue
		}

		// Per field type validators
		switch strings.ToLower(f.Kind) {
		case "bool":
			out.Push(vldtr.validateBoolFieldKind(v, f, r, m)...)
		case "datetime":
			out.Push(vldtr.validateDatetimeFieldKind(v, f, r, m)...)
		case "email":
			out.Push(vldtr.validateEmailFieldKind(v, f, r, m)...)
		case "file":
			out.Push(vldtr.validateFileFieldKind(v, f, r, m)...)
		case "number":
			out.Push(vldtr.validateNumberFieldKind(v, f, r, m)...)
		case "record":
			out.Push(vldtr.validateRecordFieldKind(v, f, r, m)...)
		case "select":
			out.Push(vldtr.validateSelectFieldKind(v, f, r, m)...)
		case "string":
			out.Push(vldtr.validateStringFieldKind(v, f, r, m)...)
		case "url":
			out.Push(vldtr.validateUrlFieldKind(v, f, r, m)...)
		case "user":
			out.Push(vldtr.validateUserFieldKind(v, f, r, m)...)
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

		duplicateRecordID, err := vldtr.uniqueCheckerFn(m, f, v)
		if err != nil {
			out.Push(makeInternalErr(f.Name, err))
		} else if duplicateRecordID > 0 && duplicateRecordID != r.ID {
			out.Push(makeDuplicateValueErr(f.Name, duplicateRecordID))
		}
	}

	if out.IsValid() {
		return nil
	}

	return out
}

func (vldtr validator) validateBoolFieldKind(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	// @todo must be 1 or 0
	return nil
}

func (vldtr validator) validateDatetimeFieldKind(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	// @todo must be datetime, UTC!
	// @todo how do we check past/future (no info about previous value)
	return nil
}

func (vldtr validator) validateEmailFieldKind(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	// @todo must be an email
	return nil
}

func (vldtr validator) validateFileFieldKind(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	// @todo file is a Ref!
	// @todo we can check for uniquenes (as we do for other refs)
	return nil
}

func (vldtr validator) validateNumberFieldKind(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	// @todo make sure it is a number and cut decimals (precision)
	return nil
}

func (vldtr validator) validateRecordFieldKind(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	// @todo record is a ref
	return nil
}

func (vldtr validator) validateSelectFieldKind(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	// @todo validate v.Value against f.Options.Strings("options")
	return nil
}

func (vldtr validator) validateStringFieldKind(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	return nil
}

func (vldtr validator) validateUrlFieldKind(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	// @todo must be URL
	return nil
}

func (vldtr validator) validateUserFieldKind(v *types.RecordValue, f *types.ModuleField, r *types.Record, m *types.Module) []types.RecordValueError {
	// @todo user is a ref
	return nil
}
