package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

type (
	AttachmentSet                []*Attachment
	ApplicationSet               []*Application
	ApigwRouteSet                []*ApigwRoute
	ApigwFilterSet               []*ApigwFilter
	AuthClientSet                []*AuthClient
	AuthConfirmedClientSet       []*AuthConfirmedClient
	AuthSessionSet               []*AuthSession
	AuthOa2tokenSet              []*AuthOa2token
	CredentialSet                []*Credential
	DataPrivacyRequestSet        []*DataPrivacyRequest
	DataPrivacyRequestCommentSet []*DataPrivacyRequestComment
	QueueSet                     []*Queue
	QueueMessageSet              []*QueueMessage
	ReminderSet                  []*Reminder
	NotificationSet              []*Notification
	ReportSet                    []*Report
	ResourceTranslationSet       []*ResourceTranslation
	RoleSet                      []*Role
	RoleMemberSet                []*RoleMember
	UserGroupSet                 []*UserGroup
	SettingValueSet              []*SettingValue
	TemplateSet                  []*Template
	UserSet                      []*User
	DalConnectionSet             []*DalConnection
	DalSensitivityLevelSet       []*DalSensitivityLevel
	DalSchemaAlterationSet       []*DalSchemaAlteration
	ConnectionSet                []*Connection
	ConfiguredConnectionSet      []*ConfiguredConnection
	LlmProviderSet               []*LlmProvider
	AgentSet                     []*Agent
	AiConversationSet            []*AiConversation
	KnowledgeBaseSet             []*KnowledgeBase
	ChatbotSet                   []*Chatbot
	ChatbotSessionSet            []*ChatbotSession
	ChatbotSessionStepSet        []*ChatbotSessionStep
	ChatbotSessionHandoffSet     []*ChatbotSessionHandoff
	TenantSet                    []*Tenant
	TenantMembershipSet          []*TenantMembership
	ProjectSet                   []*Project
	ProjectMemberSet             []*ProjectMember
	ProjectGroupSet              []*ProjectGroup
	ProjectGroupEntrySet         []*ProjectGroupEntry
	ProjectIncidentSet           []*ProjectIncident
	ProjectFeatureSet            []*ProjectFeature
	ProjectPrivacySet            []*ProjectPrivacy
	ProjectTaskSet               []*ProjectTask
	ProjectReviewSet             []*ProjectReview
	DmlConnectionSet             []*DmlConnection
	DmlMappingSet                []*DmlMapping
	DmlImportRunSet              []*DmlImportRun
	ApigwProfilerHitSet          []*ApigwProfilerHit
	ApigwProfilerAggregationSet  []*ApigwProfilerAggregation
	PrivacyDalConnectionSet      []*PrivacyDalConnection
	DmlColumnMapSet              []*DmlColumnMap
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

func (set ApplicationSet) Walk(w func(*Application) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set ApplicationSet) FindByID(ID uint64) *Application {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ApplicationSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ApigwRouteSet) Walk(w func(*ApigwRoute) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set ApigwRouteSet) FindByID(ID uint64) *ApigwRoute {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ApigwRouteSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ApigwFilterSet) Walk(w func(*ApigwFilter) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set ApigwFilterSet) FindByID(ID uint64) *ApigwFilter {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ApigwFilterSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set AuthClientSet) Walk(w func(*AuthClient) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set AuthClientSet) FindByID(ID uint64) *AuthClient {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set AuthClientSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set AuthConfirmedClientSet) Walk(w func(*AuthConfirmedClient) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set AuthSessionSet) Walk(w func(*AuthSession) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set AuthOa2tokenSet) Walk(w func(*AuthOa2token) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set AuthOa2tokenSet) FindByID(ID uint64) *AuthOa2token {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set AuthOa2tokenSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set CredentialSet) Walk(w func(*Credential) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set CredentialSet) FindByID(ID uint64) *Credential {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set CredentialSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set DataPrivacyRequestSet) Walk(w func(*DataPrivacyRequest) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set DataPrivacyRequestSet) FindByID(ID uint64) *DataPrivacyRequest {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set DataPrivacyRequestSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set DataPrivacyRequestCommentSet) Walk(w func(*DataPrivacyRequestComment) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set DataPrivacyRequestCommentSet) FindByID(ID uint64) *DataPrivacyRequestComment {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set DataPrivacyRequestCommentSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set QueueSet) Walk(w func(*Queue) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set QueueSet) FindByID(ID uint64) *Queue {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set QueueSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set QueueMessageSet) Walk(w func(*QueueMessage) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set ReminderSet) Walk(w func(*Reminder) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set ReminderSet) FindByID(ID uint64) *Reminder {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ReminderSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set NotificationSet) Walk(w func(*Notification) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set NotificationSet) Filter(f func(*Notification) (bool, error)) (out NotificationSet, err error) {
	var ok bool
	out = NotificationSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set NotificationSet) FindByID(ID uint64) *Notification {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set NotificationSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ReportSet) Walk(w func(*Report) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set ReportSet) FindByID(ID uint64) *Report {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ReportSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ResourceTranslationSet) Walk(w func(*ResourceTranslation) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set ResourceTranslationSet) FindByID(ID uint64) *ResourceTranslation {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ResourceTranslationSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set RoleSet) Walk(w func(*Role) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set RoleSet) FindByID(ID uint64) *Role {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set RoleSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set RoleMemberSet) Walk(w func(*RoleMember) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set UserGroupSet) Walk(w func(*UserGroup) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set UserGroupSet) Filter(f func(*UserGroup) (bool, error)) (out UserGroupSet, err error) {
	var ok bool
	out = UserGroupSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set UserGroupSet) FindByID(ID uint64) *UserGroup {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set UserGroupSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set SettingValueSet) Walk(w func(*SettingValue) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set TemplateSet) Walk(w func(*Template) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set TemplateSet) FindByID(ID uint64) *Template {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set TemplateSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set UserSet) Walk(w func(*User) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set UserSet) FindByID(ID uint64) *User {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set UserSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set DalConnectionSet) Walk(w func(*DalConnection) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set DalConnectionSet) FindByID(ID uint64) *DalConnection {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set DalConnectionSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set DalSensitivityLevelSet) Walk(w func(*DalSensitivityLevel) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set DalSensitivityLevelSet) FindByID(ID uint64) *DalSensitivityLevel {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set DalSensitivityLevelSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set DalSchemaAlterationSet) Walk(w func(*DalSchemaAlteration) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set DalSchemaAlterationSet) FindByID(ID uint64) *DalSchemaAlteration {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set DalSchemaAlterationSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ConnectionSet) Walk(w func(*Connection) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ConnectionSet) Filter(f func(*Connection) (bool, error)) (out ConnectionSet, err error) {
	var ok bool
	out = ConnectionSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set ConnectionSet) FindByID(ID uint64) *Connection {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ConnectionSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ConfiguredConnectionSet) Walk(w func(*ConfiguredConnection) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ConfiguredConnectionSet) Filter(f func(*ConfiguredConnection) (bool, error)) (out ConfiguredConnectionSet, err error) {
	var ok bool
	out = ConfiguredConnectionSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set ConfiguredConnectionSet) FindByID(ID uint64) *ConfiguredConnection {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ConfiguredConnectionSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set LlmProviderSet) Walk(w func(*LlmProvider) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set LlmProviderSet) Filter(f func(*LlmProvider) (bool, error)) (out LlmProviderSet, err error) {
	var ok bool
	out = LlmProviderSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set LlmProviderSet) FindByID(ID uint64) *LlmProvider {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set LlmProviderSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set AgentSet) Walk(w func(*Agent) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set AgentSet) Filter(f func(*Agent) (bool, error)) (out AgentSet, err error) {
	var ok bool
	out = AgentSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set AgentSet) FindByID(ID uint64) *Agent {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set AgentSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set AiConversationSet) Walk(w func(*AiConversation) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set AiConversationSet) Filter(f func(*AiConversation) (bool, error)) (out AiConversationSet, err error) {
	var ok bool
	out = AiConversationSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set AiConversationSet) FindByID(ID uint64) *AiConversation {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set AiConversationSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set KnowledgeBaseSet) Walk(w func(*KnowledgeBase) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set KnowledgeBaseSet) Filter(f func(*KnowledgeBase) (bool, error)) (out KnowledgeBaseSet, err error) {
	var ok bool
	out = KnowledgeBaseSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set KnowledgeBaseSet) FindByID(ID uint64) *KnowledgeBase {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set KnowledgeBaseSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ChatbotSet) Walk(w func(*Chatbot) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ChatbotSet) Filter(f func(*Chatbot) (bool, error)) (out ChatbotSet, err error) {
	var ok bool
	out = ChatbotSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set ChatbotSet) FindByID(ID uint64) *Chatbot {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ChatbotSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ChatbotSessionSet) Walk(w func(*ChatbotSession) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ChatbotSessionSet) Filter(f func(*ChatbotSession) (bool, error)) (out ChatbotSessionSet, err error) {
	var ok bool
	out = ChatbotSessionSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set ChatbotSessionSet) FindByID(ID uint64) *ChatbotSession {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ChatbotSessionSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ChatbotSessionStepSet) Walk(w func(*ChatbotSessionStep) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ChatbotSessionStepSet) Filter(f func(*ChatbotSessionStep) (bool, error)) (out ChatbotSessionStepSet, err error) {
	var ok bool
	out = ChatbotSessionStepSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set ChatbotSessionStepSet) FindByID(ID uint64) *ChatbotSessionStep {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ChatbotSessionStepSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ChatbotSessionHandoffSet) Walk(w func(*ChatbotSessionHandoff) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ChatbotSessionHandoffSet) Filter(f func(*ChatbotSessionHandoff) (bool, error)) (out ChatbotSessionHandoffSet, err error) {
	var ok bool
	out = ChatbotSessionHandoffSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set ChatbotSessionHandoffSet) FindByID(ID uint64) *ChatbotSessionHandoff {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ChatbotSessionHandoffSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set TenantSet) Walk(w func(*Tenant) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set TenantSet) Filter(f func(*Tenant) (bool, error)) (out TenantSet, err error) {
	var ok bool
	out = TenantSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set TenantSet) FindByID(ID uint64) *Tenant {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set TenantSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set TenantMembershipSet) Walk(w func(*TenantMembership) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set TenantMembershipSet) Filter(f func(*TenantMembership) (bool, error)) (out TenantMembershipSet, err error) {
	var ok bool
	out = TenantMembershipSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set TenantMembershipSet) FindByID(ID uint64) *TenantMembership {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set TenantMembershipSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ProjectSet) Walk(w func(*Project) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ProjectSet) Filter(f func(*Project) (bool, error)) (out ProjectSet, err error) {
	var ok bool
	out = ProjectSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set ProjectSet) FindByID(ID uint64) *Project {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ProjectSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ProjectMemberSet) Walk(w func(*ProjectMember) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ProjectMemberSet) Filter(f func(*ProjectMember) (bool, error)) (out ProjectMemberSet, err error) {
	var ok bool
	out = ProjectMemberSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set ProjectMemberSet) FindByID(ID uint64) *ProjectMember {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ProjectMemberSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ProjectGroupSet) Walk(w func(*ProjectGroup) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ProjectGroupSet) Filter(f func(*ProjectGroup) (bool, error)) (out ProjectGroupSet, err error) {
	var ok bool
	out = ProjectGroupSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set ProjectGroupSet) FindByID(ID uint64) *ProjectGroup {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ProjectGroupSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ProjectGroupEntrySet) Walk(w func(*ProjectGroupEntry) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ProjectGroupEntrySet) Filter(f func(*ProjectGroupEntry) (bool, error)) (out ProjectGroupEntrySet, err error) {
	var ok bool
	out = ProjectGroupEntrySet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set ProjectGroupEntrySet) FindByID(ID uint64) *ProjectGroupEntry {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ProjectGroupEntrySet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ProjectIncidentSet) Walk(w func(*ProjectIncident) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ProjectIncidentSet) Filter(f func(*ProjectIncident) (bool, error)) (out ProjectIncidentSet, err error) {
	var ok bool
	out = ProjectIncidentSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set ProjectIncidentSet) FindByID(ID uint64) *ProjectIncident {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ProjectIncidentSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ProjectFeatureSet) Walk(w func(*ProjectFeature) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ProjectFeatureSet) Filter(f func(*ProjectFeature) (bool, error)) (out ProjectFeatureSet, err error) {
	var ok bool
	out = ProjectFeatureSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set ProjectFeatureSet) FindByID(ID uint64) *ProjectFeature {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ProjectFeatureSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ProjectPrivacySet) Walk(w func(*ProjectPrivacy) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ProjectPrivacySet) Filter(f func(*ProjectPrivacy) (bool, error)) (out ProjectPrivacySet, err error) {
	var ok bool
	out = ProjectPrivacySet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set ProjectPrivacySet) FindByID(ID uint64) *ProjectPrivacy {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ProjectPrivacySet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ProjectTaskSet) Walk(w func(*ProjectTask) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ProjectTaskSet) Filter(f func(*ProjectTask) (bool, error)) (out ProjectTaskSet, err error) {
	var ok bool
	out = ProjectTaskSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set ProjectTaskSet) FindByID(ID uint64) *ProjectTask {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ProjectTaskSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ProjectReviewSet) Walk(w func(*ProjectReview) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ProjectReviewSet) Filter(f func(*ProjectReview) (bool, error)) (out ProjectReviewSet, err error) {
	var ok bool
	out = ProjectReviewSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set ProjectReviewSet) FindByID(ID uint64) *ProjectReview {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ProjectReviewSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set DmlConnectionSet) Walk(w func(*DmlConnection) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set DmlConnectionSet) Filter(f func(*DmlConnection) (bool, error)) (out DmlConnectionSet, err error) {
	var ok bool
	out = DmlConnectionSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set DmlConnectionSet) FindByID(ID uint64) *DmlConnection {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set DmlConnectionSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set DmlMappingSet) Walk(w func(*DmlMapping) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set DmlMappingSet) Filter(f func(*DmlMapping) (bool, error)) (out DmlMappingSet, err error) {
	var ok bool
	out = DmlMappingSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set DmlMappingSet) FindByID(ID uint64) *DmlMapping {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set DmlMappingSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set DmlImportRunSet) Walk(w func(*DmlImportRun) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set DmlImportRunSet) Filter(f func(*DmlImportRun) (bool, error)) (out DmlImportRunSet, err error) {
	var ok bool
	out = DmlImportRunSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set DmlImportRunSet) FindByID(ID uint64) *DmlImportRun {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set DmlImportRunSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ApigwProfilerHitSet) Walk(w func(*ApigwProfilerHit) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set ApigwProfilerAggregationSet) Walk(w func(*ApigwProfilerAggregation) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set PrivacyDalConnectionSet) Walk(w func(*PrivacyDalConnection) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set PrivacyDalConnectionSet) FindByID(ID uint64) *PrivacyDalConnection {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set PrivacyDalConnectionSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set DmlColumnMapSet) Walk(w func(*DmlColumnMap) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set DmlColumnMapSet) Filter(f func(*DmlColumnMap) (bool, error)) (out DmlColumnMapSet, err error) {
	var ok bool
	out = DmlColumnMapSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}
