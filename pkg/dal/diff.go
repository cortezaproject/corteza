package dal

type (
	modelDiffType string

	// ModelDiff defines one identified missmatch between two models
	ModelDiff struct {
		Type modelDiffType
		// Original will be nil when a new attribute is being added
		Original *Attribute
		// Asserted will be nil wen an existing attribute is being removed
		Asserted *Attribute
	}

	ModelDiffSet []*ModelDiff
)

const (
	AttributeMissing              modelDiffType = "attributeMissing"
	AttributeTypeMissmatch        modelDiffType = "typeMissmatch"
	AttributeSensitivityMissmatch modelDiffType = "sensitivityMissmatch"
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

	aIndex := make(map[string]struct {
		found bool
		attr  *Attribute
	})
	for _, _attr := range a.Attributes {
		attr := _attr
		aIndex[attr.Ident] = struct {
			found bool
			attr  *Attribute
		}{
			attr: attr,
		}
	}

	// Deleted and update ones
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

		// Other stuff
		// @todo improve; for now it'll do
		if attrA.SensitivityLevel != attrBAux.attr.SensitivityLevel {
			out = append(out, &ModelDiff{
				Type:     AttributeSensitivityMissmatch,
				Original: attrA,
				Asserted: attrBAux.attr,
			})
		}
	}

	// New
	for _, _attrB := range b.Attributes {
		attrB := _attrB

		// Missmatches
		_, ok := aIndex[attrB.Ident]
		if !ok {
			out = append(out, &ModelDiff{
				Type:     AttributeMissing,
				Original: nil,
				Asserted: attrB,
			})
			continue
		}
	}

	return
}
