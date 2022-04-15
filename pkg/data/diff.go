package data

type (
	modelDiffType string

	// ModelDiff defines one identified missmatch between two models
	ModelDiff struct {
		Type     modelDiffType
		Original *Attribute
		Asserted *Attribute
	}

	ModelDiffSet []*ModelDiff
)

const (
	AttributeMissing       modelDiffType = "attributeMissing"
	AttributeTypeMissmatch modelDiffType = "typeMissmatch"
)

// Diff calculates the diff between models a and b where a is used as base
func (a *Model) Diff(b *Model) (out ModelDiffSet) {
	bIndex := make(map[string]struct {
		found bool
		attr  *Attribute
	})
	for _, _attr := range b.Attributes {
		attr := _attr
		bIndex[attr.Ident] = struct {
			found bool
			attr  *Attribute
		}{
			attr: attr,
		}
	}

	for _, _attrA := range a.Attributes {
		attrA := _attrA

		// Missmatches
		attrBAux, ok := bIndex[attrA.Ident]
		if !ok {
			out = append(out, &ModelDiff{
				Type:     AttributeMissing,
				Original: attrA,
			})
			continue
		}

		// Typecheck
		if attrA.Type.Type() != attrBAux.attr.Type.Type() {
			out = append(out, &ModelDiff{
				Type:     AttributeTypeMissmatch,
				Original: attrA,
				Asserted: attrBAux.attr,
			})
		}
	}

	return
}
