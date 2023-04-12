package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// compose/types/types.yaml

type (

	// AttachmentSet slice of Attachment
	//
	// This type is auto-generated.
	AttachmentSet []*Attachment

	// ChartSet slice of Chart
	//
	// This type is auto-generated.
	ChartSet []*Chart

	// DeDupRuleSet slice of DeDupRule
	//
	// This type is auto-generated.
	DeDupRuleSet []*DeDupRule

	// IconSet slice of Icon
	//
	// This type is auto-generated.
	IconSet []*Icon

	// ModuleSet slice of Module
	//
	// This type is auto-generated.
	ModuleSet []*Module

	// ModuleFieldSet slice of ModuleField
	//
	// This type is auto-generated.
	ModuleFieldSet []*ModuleField

	// NamespaceSet slice of Namespace
	//
	// This type is auto-generated.
	NamespaceSet []*Namespace

	// PageSet slice of Page
	//
	// This type is auto-generated.
	PageSet []*Page

	// PageLayoutSet slice of PageLayout
	//
	// This type is auto-generated.
	PageLayoutSet []*PageLayout

	// PrivacyModuleSet slice of PrivacyModule
	//
	// This type is auto-generated.
	PrivacyModuleSet []*PrivacyModule

	// RecordSet slice of Record
	//
	// This type is auto-generated.
	RecordSet []*Record

	// RecordValueSet slice of RecordValue
	//
	// This type is auto-generated.
	RecordValueSet []*RecordValue
)

// Walk iterates through every slice item and calls w(Attachment) err
//
// This function is auto-generated.
func (set AttachmentSet) Walk(w func(*Attachment) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Attachment) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set AttachmentSet) Filter(f func(*Attachment) (bool, error)) (out AttachmentSet, err error) {
	var ok bool
	out = AttachmentSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set AttachmentSet) FindByID(ID uint64) *Attachment {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set AttachmentSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(Chart) err
//
// This function is auto-generated.
func (set ChartSet) Walk(w func(*Chart) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Chart) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ChartSet) Filter(f func(*Chart) (bool, error)) (out ChartSet, err error) {
	var ok bool
	out = ChartSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set ChartSet) FindByID(ID uint64) *Chart {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set ChartSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(DeDupRule) err
//
// This function is auto-generated.
func (set DeDupRuleSet) Walk(w func(*DeDupRule) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(DeDupRule) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set DeDupRuleSet) Filter(f func(*DeDupRule) (bool, error)) (out DeDupRuleSet, err error) {
	var ok bool
	out = DeDupRuleSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// Walk iterates through every slice item and calls w(Icon) err
//
// This function is auto-generated.
func (set IconSet) Walk(w func(*Icon) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Icon) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set IconSet) Filter(f func(*Icon) (bool, error)) (out IconSet, err error) {
	var ok bool
	out = IconSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// Walk iterates through every slice item and calls w(Module) err
//
// This function is auto-generated.
func (set ModuleSet) Walk(w func(*Module) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Module) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ModuleSet) Filter(f func(*Module) (bool, error)) (out ModuleSet, err error) {
	var ok bool
	out = ModuleSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set ModuleSet) FindByID(ID uint64) *Module {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set ModuleSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(ModuleField) err
//
// This function is auto-generated.
func (set ModuleFieldSet) Walk(w func(*ModuleField) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(ModuleField) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ModuleFieldSet) Filter(f func(*ModuleField) (bool, error)) (out ModuleFieldSet, err error) {
	var ok bool
	out = ModuleFieldSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set ModuleFieldSet) FindByID(ID uint64) *ModuleField {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set ModuleFieldSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(Namespace) err
//
// This function is auto-generated.
func (set NamespaceSet) Walk(w func(*Namespace) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Namespace) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set NamespaceSet) Filter(f func(*Namespace) (bool, error)) (out NamespaceSet, err error) {
	var ok bool
	out = NamespaceSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set NamespaceSet) FindByID(ID uint64) *Namespace {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set NamespaceSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(Page) err
//
// This function is auto-generated.
func (set PageSet) Walk(w func(*Page) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Page) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set PageSet) Filter(f func(*Page) (bool, error)) (out PageSet, err error) {
	var ok bool
	out = PageSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set PageSet) FindByID(ID uint64) *Page {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set PageSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(PageLayout) err
//
// This function is auto-generated.
func (set PageLayoutSet) Walk(w func(*PageLayout) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(PageLayout) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set PageLayoutSet) Filter(f func(*PageLayout) (bool, error)) (out PageLayoutSet, err error) {
	var ok bool
	out = PageLayoutSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set PageLayoutSet) FindByID(ID uint64) *PageLayout {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set PageLayoutSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(PrivacyModule) err
//
// This function is auto-generated.
func (set PrivacyModuleSet) Walk(w func(*PrivacyModule) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(PrivacyModule) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set PrivacyModuleSet) Filter(f func(*PrivacyModule) (bool, error)) (out PrivacyModuleSet, err error) {
	var ok bool
	out = PrivacyModuleSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// Walk iterates through every slice item and calls w(Record) err
//
// This function is auto-generated.
func (set RecordSet) Walk(w func(*Record) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Record) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set RecordSet) Filter(f func(*Record) (bool, error)) (out RecordSet, err error) {
	var ok bool
	out = RecordSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set RecordSet) FindByID(ID uint64) *Record {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set RecordSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(RecordValue) err
//
// This function is auto-generated.
func (set RecordValueSet) Walk(w func(*RecordValue) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(RecordValue) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set RecordValueSet) Filter(f func(*RecordValue) (bool, error)) (out RecordValueSet, err error) {
	var ok bool
	out = RecordValueSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}
