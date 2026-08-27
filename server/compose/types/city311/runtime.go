package city311

import "time"

// PortalServiceRequestSubmit is the public-portal submission DTO.
type PortalServiceRequestSubmit struct {
	Summary          string         `json:"summary"`
	Description      string         `json:"description"`
	ServiceType      ServiceType    `json:"service_type"`
	Requester        RequesterInput `json:"requester"`
	Location         *LocationInput `json:"location,omitempty"`
	CustomFields     map[string]any `json:"custom_fields,omitempty"`
	AttachmentTokens []string       `json:"attachment_tokens,omitempty"`
}

// StaffServiceRequestCreate is the contract request used by staff intake.
type StaffServiceRequestCreate struct {
	Request     PortalServiceRequestSubmit `json:"request"`
	Constituent StaffConstituentInput      `json:"constituent"`
}

// StaffConstituentInput represents the contract's constituent-reference XOR
// constituent-create union. Pointers preserve field presence for oneOf checks.
type StaffConstituentInput struct {
	ConstituentID *string `json:"constituent_id,omitempty"`
	DisplayName   *string `json:"display_name,omitempty"`
	Email         *string `json:"email,omitempty"`
}

type RequestTransition struct {
	ToStatus ServiceRequestStatus `json:"to_status"`
	Reason   string               `json:"reason,omitempty"`
}

type AuditEvent struct {
	EntityType    string         `json:"entity_type"`
	EntityID      string         `json:"entity_id"`
	EventType     string         `json:"event_type"`
	ActorType     AuditActorType `json:"actor_type"`
	ActorID       string         `json:"actor_id"`
	OccurredAt    time.Time      `json:"occurred_at"`
	SourceChannel SourceChannel  `json:"source_channel"`
	Before        map[string]any `json:"before"`
	After         map[string]any `json:"after"`
}

type StaffServiceRequestDetail struct {
	Request           ServiceRequest      `json:"request"`
	AvailableActions  []string            `json:"available_actions"`
	PrimaryAssigneeID *string             `json:"primary_assignee_id"`
	CollaboratorIDs   []string            `json:"collaborator_ids"`
	Reminders         []any               `json:"reminders"`
	History           []PublicHistoryItem `json:"history"`
	Audit             []AuditEvent        `json:"audit"`
	ExternalWorkOrder any                 `json:"external_work_order"`
}

type RequestQueueItem struct {
	RequestID         string               `json:"request_id"`
	RequestNumber     string               `json:"request_number"`
	Summary           string               `json:"summary"`
	ServiceType       ServiceType          `json:"service_type"`
	Status            ServiceRequestStatus `json:"status"`
	OwningDepartment  DepartmentCode       `json:"owning_department"`
	CouncilDistrict   DistrictCode         `json:"council_district,omitempty"`
	OriginClass       OriginClass          `json:"origin_class"`
	SourceChannel     SourceChannel        `json:"source_channel"`
	PrimaryAssigneeID *string              `json:"primary_assignee_id"`
	DuplicateGroupID  *string              `json:"duplicate_group_id"`
	Version           uint64               `json:"version"`
	UpdatedAt         time.Time            `json:"updated_at"`
	AvailableActions  []string             `json:"available_actions"`
}

type ListResponse struct {
	Items          []RequestQueueItem `json:"items"`
	NextPageToken  *string            `json:"next_page_token"`
	TotalCount     int                `json:"total_count"`
	AppliedFilters map[string]any     `json:"applied_filters"`
	Sort           []string           `json:"sort"`
}

// Actor carries server-resolved roles and record scope.
type Actor struct {
	ID         uint64
	Roles      []ApplicationRole
	Department DepartmentCode
	Districts  []DistrictCode
}
