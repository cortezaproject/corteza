package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

type (
	AttachmentSet    []*Attachment
	ChartSet         []*Chart
	ModuleSet        []*Module
	ModuleFieldSet   []*ModuleField
	NamespaceSet     []*Namespace
	PageSet          []*Page
	PageLayoutSet    []*PageLayout
	RecordSet        []*Record
	IconSet          []*Icon
	RecordValueSet   []*RecordValue
	PrivacyModuleSet []*PrivacyModule
	DeDupRuleSet     []*DeDupRule
)

func (set AttachmentSet) Walk(w func(*Attachment) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set AttachmentSet) FindByID(ID uint64) *Attachment {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set AttachmentSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ChartSet) Walk(w func(*Chart) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set ChartSet) FindByID(ID uint64) *Chart {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ChartSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ModuleSet) Walk(w func(*Module) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set ModuleSet) FindByID(ID uint64) *Module {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ModuleSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ModuleFieldSet) Walk(w func(*ModuleField) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set ModuleFieldSet) FindByID(ID uint64) *ModuleField {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ModuleFieldSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set NamespaceSet) Walk(w func(*Namespace) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set NamespaceSet) FindByID(ID uint64) *Namespace {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set NamespaceSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set PageSet) Walk(w func(*Page) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set PageSet) FindByID(ID uint64) *Page {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set PageSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set PageLayoutSet) Walk(w func(*PageLayout) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set PageLayoutSet) FindByID(ID uint64) *PageLayout {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set PageLayoutSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set RecordSet) Walk(w func(*Record) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set RecordSet) FindByID(ID uint64) *Record {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set RecordSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set IconSet) Walk(w func(*Icon) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set RecordValueSet) Walk(w func(*RecordValue) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set PrivacyModuleSet) Walk(w func(*PrivacyModule) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set DeDupRuleSet) Walk(w func(*DeDupRule) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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
