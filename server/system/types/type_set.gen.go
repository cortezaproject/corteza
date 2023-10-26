package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/types/types.yaml

type (

	// ApigwFilterSet slice of ApigwFilter
	//
	// This type is auto-generated.
	ApigwFilterSet []*ApigwFilter

	// ApigwProfilerAggregationSet slice of ApigwProfilerAggregation
	//
	// This type is auto-generated.
	ApigwProfilerAggregationSet []*ApigwProfilerAggregation

	// ApigwProfilerHitSet slice of ApigwProfilerHit
	//
	// This type is auto-generated.
	ApigwProfilerHitSet []*ApigwProfilerHit

	// ApigwRouteSet slice of ApigwRoute
	//
	// This type is auto-generated.
	ApigwRouteSet []*ApigwRoute

	// ApplicationSet slice of Application
	//
	// This type is auto-generated.
	ApplicationSet []*Application

	// AttachmentSet slice of Attachment
	//
	// This type is auto-generated.
	AttachmentSet []*Attachment

	// AuthClientSet slice of AuthClient
	//
	// This type is auto-generated.
	AuthClientSet []*AuthClient

	// AuthConfirmedClientSet slice of AuthConfirmedClient
	//
	// This type is auto-generated.
	AuthConfirmedClientSet []*AuthConfirmedClient

	// AuthOa2tokenSet slice of AuthOa2token
	//
	// This type is auto-generated.
	AuthOa2tokenSet []*AuthOa2token

	// AuthSessionSet slice of AuthSession
	//
	// This type is auto-generated.
	AuthSessionSet []*AuthSession

	// CredentialSet slice of Credential
	//
	// This type is auto-generated.
	CredentialSet []*Credential

	// DalConnectionSet slice of DalConnection
	//
	// This type is auto-generated.
	DalConnectionSet []*DalConnection

	// DalSchemaAlterationSet slice of DalSchemaAlteration
	//
	// This type is auto-generated.
	DalSchemaAlterationSet []*DalSchemaAlteration

	// DalSensitivityLevelSet slice of DalSensitivityLevel
	//
	// This type is auto-generated.
	DalSensitivityLevelSet []*DalSensitivityLevel

	// DataPrivacyRequestSet slice of DataPrivacyRequest
	//
	// This type is auto-generated.
	DataPrivacyRequestSet []*DataPrivacyRequest

	// DataPrivacyRequestCommentSet slice of DataPrivacyRequestComment
	//
	// This type is auto-generated.
	DataPrivacyRequestCommentSet []*DataPrivacyRequestComment

	// PrivacyDalConnectionSet slice of PrivacyDalConnection
	//
	// This type is auto-generated.
	PrivacyDalConnectionSet []*PrivacyDalConnection

	// QueueSet slice of Queue
	//
	// This type is auto-generated.
	QueueSet []*Queue

	// QueueMessageSet slice of QueueMessage
	//
	// This type is auto-generated.
	QueueMessageSet []*QueueMessage

	// ReminderSet slice of Reminder
	//
	// This type is auto-generated.
	ReminderSet []*Reminder

	// ReportSet slice of Report
	//
	// This type is auto-generated.
	ReportSet []*Report

	// ResourceTranslationSet slice of ResourceTranslation
	//
	// This type is auto-generated.
	ResourceTranslationSet []*ResourceTranslation

	// RoleSet slice of Role
	//
	// This type is auto-generated.
	RoleSet []*Role

	// RoleMemberSet slice of RoleMember
	//
	// This type is auto-generated.
	RoleMemberSet []*RoleMember

	// SettingValueSet slice of SettingValue
	//
	// This type is auto-generated.
	SettingValueSet []*SettingValue

	// TemplateSet slice of Template
	//
	// This type is auto-generated.
	TemplateSet []*Template

	// UserSet slice of User
	//
	// This type is auto-generated.
	UserSet []*User
)

// Walk iterates through every slice item and calls w(ApigwFilter) err
//
// This function is auto-generated.
func (set ApigwFilterSet) Walk(w func(*ApigwFilter) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(ApigwFilter) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ApigwFilterSet) Filter(f func(*ApigwFilter) (bool, error)) (out ApigwFilterSet, err error) {
	var ok bool
	out = ApigwFilterSet{}
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
func (set ApigwFilterSet) FindByID(ID uint64) *ApigwFilter {
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
func (set ApigwFilterSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(ApigwProfilerAggregation) err
//
// This function is auto-generated.
func (set ApigwProfilerAggregationSet) Walk(w func(*ApigwProfilerAggregation) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(ApigwProfilerAggregation) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ApigwProfilerAggregationSet) Filter(f func(*ApigwProfilerAggregation) (bool, error)) (out ApigwProfilerAggregationSet, err error) {
	var ok bool
	out = ApigwProfilerAggregationSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// Walk iterates through every slice item and calls w(ApigwProfilerHit) err
//
// This function is auto-generated.
func (set ApigwProfilerHitSet) Walk(w func(*ApigwProfilerHit) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(ApigwProfilerHit) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ApigwProfilerHitSet) Filter(f func(*ApigwProfilerHit) (bool, error)) (out ApigwProfilerHitSet, err error) {
	var ok bool
	out = ApigwProfilerHitSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// Walk iterates through every slice item and calls w(ApigwRoute) err
//
// This function is auto-generated.
func (set ApigwRouteSet) Walk(w func(*ApigwRoute) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(ApigwRoute) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ApigwRouteSet) Filter(f func(*ApigwRoute) (bool, error)) (out ApigwRouteSet, err error) {
	var ok bool
	out = ApigwRouteSet{}
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
func (set ApigwRouteSet) FindByID(ID uint64) *ApigwRoute {
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
func (set ApigwRouteSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(Application) err
//
// This function is auto-generated.
func (set ApplicationSet) Walk(w func(*Application) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Application) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ApplicationSet) Filter(f func(*Application) (bool, error)) (out ApplicationSet, err error) {
	var ok bool
	out = ApplicationSet{}
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
func (set ApplicationSet) FindByID(ID uint64) *Application {
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
func (set ApplicationSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

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

// Walk iterates through every slice item and calls w(AuthClient) err
//
// This function is auto-generated.
func (set AuthClientSet) Walk(w func(*AuthClient) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(AuthClient) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set AuthClientSet) Filter(f func(*AuthClient) (bool, error)) (out AuthClientSet, err error) {
	var ok bool
	out = AuthClientSet{}
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
func (set AuthClientSet) FindByID(ID uint64) *AuthClient {
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
func (set AuthClientSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(AuthConfirmedClient) err
//
// This function is auto-generated.
func (set AuthConfirmedClientSet) Walk(w func(*AuthConfirmedClient) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(AuthConfirmedClient) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set AuthConfirmedClientSet) Filter(f func(*AuthConfirmedClient) (bool, error)) (out AuthConfirmedClientSet, err error) {
	var ok bool
	out = AuthConfirmedClientSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// Walk iterates through every slice item and calls w(AuthOa2token) err
//
// This function is auto-generated.
func (set AuthOa2tokenSet) Walk(w func(*AuthOa2token) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(AuthOa2token) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set AuthOa2tokenSet) Filter(f func(*AuthOa2token) (bool, error)) (out AuthOa2tokenSet, err error) {
	var ok bool
	out = AuthOa2tokenSet{}
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
func (set AuthOa2tokenSet) FindByID(ID uint64) *AuthOa2token {
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
func (set AuthOa2tokenSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(AuthSession) err
//
// This function is auto-generated.
func (set AuthSessionSet) Walk(w func(*AuthSession) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(AuthSession) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set AuthSessionSet) Filter(f func(*AuthSession) (bool, error)) (out AuthSessionSet, err error) {
	var ok bool
	out = AuthSessionSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// Walk iterates through every slice item and calls w(Credential) err
//
// This function is auto-generated.
func (set CredentialSet) Walk(w func(*Credential) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Credential) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set CredentialSet) Filter(f func(*Credential) (bool, error)) (out CredentialSet, err error) {
	var ok bool
	out = CredentialSet{}
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
func (set CredentialSet) FindByID(ID uint64) *Credential {
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
func (set CredentialSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(DalConnection) err
//
// This function is auto-generated.
func (set DalConnectionSet) Walk(w func(*DalConnection) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(DalConnection) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set DalConnectionSet) Filter(f func(*DalConnection) (bool, error)) (out DalConnectionSet, err error) {
	var ok bool
	out = DalConnectionSet{}
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
func (set DalConnectionSet) FindByID(ID uint64) *DalConnection {
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
func (set DalConnectionSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(DalSchemaAlteration) err
//
// This function is auto-generated.
func (set DalSchemaAlterationSet) Walk(w func(*DalSchemaAlteration) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(DalSchemaAlteration) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set DalSchemaAlterationSet) Filter(f func(*DalSchemaAlteration) (bool, error)) (out DalSchemaAlterationSet, err error) {
	var ok bool
	out = DalSchemaAlterationSet{}
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
func (set DalSchemaAlterationSet) FindByID(ID uint64) *DalSchemaAlteration {
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
func (set DalSchemaAlterationSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(DalSensitivityLevel) err
//
// This function is auto-generated.
func (set DalSensitivityLevelSet) Walk(w func(*DalSensitivityLevel) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(DalSensitivityLevel) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set DalSensitivityLevelSet) Filter(f func(*DalSensitivityLevel) (bool, error)) (out DalSensitivityLevelSet, err error) {
	var ok bool
	out = DalSensitivityLevelSet{}
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
func (set DalSensitivityLevelSet) FindByID(ID uint64) *DalSensitivityLevel {
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
func (set DalSensitivityLevelSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(DataPrivacyRequest) err
//
// This function is auto-generated.
func (set DataPrivacyRequestSet) Walk(w func(*DataPrivacyRequest) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(DataPrivacyRequest) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set DataPrivacyRequestSet) Filter(f func(*DataPrivacyRequest) (bool, error)) (out DataPrivacyRequestSet, err error) {
	var ok bool
	out = DataPrivacyRequestSet{}
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
func (set DataPrivacyRequestSet) FindByID(ID uint64) *DataPrivacyRequest {
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
func (set DataPrivacyRequestSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(DataPrivacyRequestComment) err
//
// This function is auto-generated.
func (set DataPrivacyRequestCommentSet) Walk(w func(*DataPrivacyRequestComment) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(DataPrivacyRequestComment) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set DataPrivacyRequestCommentSet) Filter(f func(*DataPrivacyRequestComment) (bool, error)) (out DataPrivacyRequestCommentSet, err error) {
	var ok bool
	out = DataPrivacyRequestCommentSet{}
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
func (set DataPrivacyRequestCommentSet) FindByID(ID uint64) *DataPrivacyRequestComment {
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
func (set DataPrivacyRequestCommentSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(PrivacyDalConnection) err
//
// This function is auto-generated.
func (set PrivacyDalConnectionSet) Walk(w func(*PrivacyDalConnection) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(PrivacyDalConnection) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set PrivacyDalConnectionSet) Filter(f func(*PrivacyDalConnection) (bool, error)) (out PrivacyDalConnectionSet, err error) {
	var ok bool
	out = PrivacyDalConnectionSet{}
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
func (set PrivacyDalConnectionSet) FindByID(ID uint64) *PrivacyDalConnection {
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
func (set PrivacyDalConnectionSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(Queue) err
//
// This function is auto-generated.
func (set QueueSet) Walk(w func(*Queue) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Queue) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set QueueSet) Filter(f func(*Queue) (bool, error)) (out QueueSet, err error) {
	var ok bool
	out = QueueSet{}
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
func (set QueueSet) FindByID(ID uint64) *Queue {
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
func (set QueueSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(QueueMessage) err
//
// This function is auto-generated.
func (set QueueMessageSet) Walk(w func(*QueueMessage) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(QueueMessage) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set QueueMessageSet) Filter(f func(*QueueMessage) (bool, error)) (out QueueMessageSet, err error) {
	var ok bool
	out = QueueMessageSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// Walk iterates through every slice item and calls w(Reminder) err
//
// This function is auto-generated.
func (set ReminderSet) Walk(w func(*Reminder) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Reminder) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ReminderSet) Filter(f func(*Reminder) (bool, error)) (out ReminderSet, err error) {
	var ok bool
	out = ReminderSet{}
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
func (set ReminderSet) FindByID(ID uint64) *Reminder {
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
func (set ReminderSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(Report) err
//
// This function is auto-generated.
func (set ReportSet) Walk(w func(*Report) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Report) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ReportSet) Filter(f func(*Report) (bool, error)) (out ReportSet, err error) {
	var ok bool
	out = ReportSet{}
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
func (set ReportSet) FindByID(ID uint64) *Report {
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
func (set ReportSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(ResourceTranslation) err
//
// This function is auto-generated.
func (set ResourceTranslationSet) Walk(w func(*ResourceTranslation) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(ResourceTranslation) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ResourceTranslationSet) Filter(f func(*ResourceTranslation) (bool, error)) (out ResourceTranslationSet, err error) {
	var ok bool
	out = ResourceTranslationSet{}
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
func (set ResourceTranslationSet) FindByID(ID uint64) *ResourceTranslation {
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
func (set ResourceTranslationSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(Role) err
//
// This function is auto-generated.
func (set RoleSet) Walk(w func(*Role) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Role) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set RoleSet) Filter(f func(*Role) (bool, error)) (out RoleSet, err error) {
	var ok bool
	out = RoleSet{}
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
func (set RoleSet) FindByID(ID uint64) *Role {
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
func (set RoleSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(RoleMember) err
//
// This function is auto-generated.
func (set RoleMemberSet) Walk(w func(*RoleMember) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(RoleMember) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set RoleMemberSet) Filter(f func(*RoleMember) (bool, error)) (out RoleMemberSet, err error) {
	var ok bool
	out = RoleMemberSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// Walk iterates through every slice item and calls w(SettingValue) err
//
// This function is auto-generated.
func (set SettingValueSet) Walk(w func(*SettingValue) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(SettingValue) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set SettingValueSet) Filter(f func(*SettingValue) (bool, error)) (out SettingValueSet, err error) {
	var ok bool
	out = SettingValueSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// Walk iterates through every slice item and calls w(Template) err
//
// This function is auto-generated.
func (set TemplateSet) Walk(w func(*Template) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Template) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set TemplateSet) Filter(f func(*Template) (bool, error)) (out TemplateSet, err error) {
	var ok bool
	out = TemplateSet{}
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
func (set TemplateSet) FindByID(ID uint64) *Template {
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
func (set TemplateSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(User) err
//
// This function is auto-generated.
func (set UserSet) Walk(w func(*User) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(User) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set UserSet) Filter(f func(*User) (bool, error)) (out UserSet, err error) {
	var ok bool
	out = UserSet{}
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
func (set UserSet) FindByID(ID uint64) *User {
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
func (set UserSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
