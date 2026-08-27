package city311

import "time"

// ContractVersion identifies the frozen City 311 client/server contract.
const ContractVersion = "1.0.0"

type (
	ServiceRequestStatus string
	SourceChannel        string
	OriginClass          string
	RelationshipType     string
	ReminderStatus       string
	ReminderChannel      string
	Language             string
	CustomFieldType      string
	ContactCategory      string
	ServiceType          string
	DepartmentCode       string
	DistrictCode         string
	CivicWorksStatus     string
	ActorRole            string
	ApplicationRole      string
	AuditActorType       string
	PhoneLabel           string
	ValidationCode       string
	ErrorCode            string
)

const (
	ServiceRequestStatusDraft      ServiceRequestStatus = "DRAFT"
	ServiceRequestStatusSubmitted  ServiceRequestStatus = "SUBMITTED"
	ServiceRequestStatusTriaged    ServiceRequestStatus = "TRIAGED"
	ServiceRequestStatusAssigned   ServiceRequestStatus = "ASSIGNED"
	ServiceRequestStatusInProgress ServiceRequestStatus = "IN_PROGRESS"
	ServiceRequestStatusResolved   ServiceRequestStatus = "RESOLVED"
	ServiceRequestStatusClosed     ServiceRequestStatus = "CLOSED"
	ServiceRequestStatusReopened   ServiceRequestStatus = "REOPENED"
)

var ServiceRequestStatuses = []ServiceRequestStatus{
	ServiceRequestStatusDraft,
	ServiceRequestStatusSubmitted,
	ServiceRequestStatusTriaged,
	ServiceRequestStatusAssigned,
	ServiceRequestStatusInProgress,
	ServiceRequestStatusResolved,
	ServiceRequestStatusClosed,
	ServiceRequestStatusReopened,
}

var ServiceRequestTransitions = map[ServiceRequestStatus][]ServiceRequestStatus{
	ServiceRequestStatusDraft:      {ServiceRequestStatusSubmitted},
	ServiceRequestStatusSubmitted:  {ServiceRequestStatusTriaged},
	ServiceRequestStatusTriaged:    {ServiceRequestStatusAssigned},
	ServiceRequestStatusAssigned:   {ServiceRequestStatusInProgress},
	ServiceRequestStatusInProgress: {ServiceRequestStatusResolved},
	ServiceRequestStatusResolved:   {ServiceRequestStatusClosed, ServiceRequestStatusReopened},
	ServiceRequestStatusClosed:     {ServiceRequestStatusReopened},
	ServiceRequestStatusReopened:   {ServiceRequestStatusAssigned, ServiceRequestStatusInProgress},
}

const (
	SourceChannelPortalAnonymous     SourceChannel = "PORTAL_ANONYMOUS"
	SourceChannelPortalAuthenticated SourceChannel = "PORTAL_AUTHENTICATED"
	SourceChannelStaffInPerson       SourceChannel = "STAFF_IN_PERSON"
	SourceChannelAPI                 SourceChannel = "API"
)

var SourceChannels = []SourceChannel{
	SourceChannelPortalAnonymous,
	SourceChannelPortalAuthenticated,
	SourceChannelStaffInPerson,
	SourceChannelAPI,
}

const (
	OriginClassExternal OriginClass = "EXTERNAL"
	OriginClassInternal OriginClass = "INTERNAL"
)

var OriginClasses = []OriginClass{OriginClassExternal, OriginClassInternal}

const (
	RelationshipPrimaryRequester    RelationshipType = "PRIMARY_REQUESTER"
	RelationshipAffectedResident    RelationshipType = "AFFECTED_RESIDENT"
	RelationshipPropertyOwner       RelationshipType = "PROPERTY_OWNER"
	RelationshipReporter            RelationshipType = "REPORTER"
	RelationshipOrganisationContact RelationshipType = "ORGANISATION_CONTACT"
)

var RelationshipTypes = []RelationshipType{
	RelationshipPrimaryRequester,
	RelationshipAffectedResident,
	RelationshipPropertyOwner,
	RelationshipReporter,
	RelationshipOrganisationContact,
}

const (
	ReminderStatusScheduled ReminderStatus  = "SCHEDULED"
	ReminderStatusSnoozed   ReminderStatus  = "SNOOZED"
	ReminderStatusCompleted ReminderStatus  = "COMPLETED"
	ReminderStatusCancelled ReminderStatus  = "CANCELLED"
	ReminderChannelEmail    ReminderChannel = "EMAIL"
	ReminderChannelInApp    ReminderChannel = "IN_APP"
)

var (
	ReminderStatuses = []ReminderStatus{ReminderStatusScheduled, ReminderStatusSnoozed, ReminderStatusCompleted, ReminderStatusCancelled}
	ReminderChannels = []ReminderChannel{ReminderChannelEmail, ReminderChannelInApp}
)

const (
	LanguageEN Language = "EN"
	LanguageES Language = "ES"
	LanguageVI Language = "VI"
)

var Languages = []Language{LanguageEN, LanguageES, LanguageVI}

const (
	PhoneLabelMobile PhoneLabel = "MOBILE"
	PhoneLabelHome   PhoneLabel = "HOME"
	PhoneLabelWork   PhoneLabel = "WORK"
)

var PhoneLabels = []PhoneLabel{PhoneLabelMobile, PhoneLabelHome, PhoneLabelWork}

const (
	CustomFieldTypeText         CustomFieldType = "TEXT"
	CustomFieldTypeInteger      CustomFieldType = "INTEGER"
	CustomFieldTypeDecimal      CustomFieldType = "DECIMAL"
	CustomFieldTypeDate         CustomFieldType = "DATE"
	CustomFieldTypeDateTime     CustomFieldType = "DATETIME"
	CustomFieldTypeBoolean      CustomFieldType = "BOOLEAN"
	CustomFieldTypeSingleChoice CustomFieldType = "SINGLE_CHOICE"
	CustomFieldTypeMultiChoice  CustomFieldType = "MULTI_CHOICE"
)

var CustomFieldTypes = []CustomFieldType{
	CustomFieldTypeText,
	CustomFieldTypeInteger,
	CustomFieldTypeDecimal,
	CustomFieldTypeDate,
	CustomFieldTypeDateTime,
	CustomFieldTypeBoolean,
	CustomFieldTypeSingleChoice,
	CustomFieldTypeMultiChoice,
}

const (
	ContactCategoryResident                ContactCategory = "RESIDENT"
	ContactCategoryBusiness                ContactCategory = "BUSINESS"
	ContactCategoryBusinessOwner           ContactCategory = "BUSINESS_OWNER"
	ContactCategoryVeteran                 ContactCategory = "VETERAN"
	ContactCategoryNeighborhoodAssociation ContactCategory = "NEIGHBORHOOD_ASSOCIATION"
	ContactCategoryGovernment              ContactCategory = "GOVERNMENT"
	ContactCategoryOther                   ContactCategory = "OTHER"
)

var ContactCategories = []ContactCategory{
	ContactCategoryResident,
	ContactCategoryBusiness,
	ContactCategoryBusinessOwner,
	ContactCategoryVeteran,
	ContactCategoryNeighborhoodAssociation,
	ContactCategoryGovernment,
	ContactCategoryOther,
}

const (
	ServiceTypeTreeMaintenance ServiceType    = "TREE_MAINTENANCE"
	ServiceTypePothole         ServiceType    = "POTHOLE"
	ServiceTypeMissedTrash     ServiceType    = "MISSED_TRASH"
	ServiceTypeGeneralInquiry  ServiceType    = "GENERAL_INQUIRY"
	DepartmentPublicWorks      DepartmentCode = "PUBLIC_WORKS"
	DepartmentStreets          DepartmentCode = "STREETS"
	DepartmentSanitation       DepartmentCode = "SANITATION"
	DepartmentGeneralServices  DepartmentCode = "GENERAL_SERVICES"
	DistrictNorth              DistrictCode   = "NORTH"
	DistrictCentral            DistrictCode   = "CENTRAL"
	DistrictSouth              DistrictCode   = "SOUTH"
)

var (
	ServiceTypes    = []ServiceType{ServiceTypeTreeMaintenance, ServiceTypePothole, ServiceTypeMissedTrash, ServiceTypeGeneralInquiry}
	DepartmentCodes = []DepartmentCode{DepartmentPublicWorks, DepartmentStreets, DepartmentSanitation, DepartmentGeneralServices}
	DistrictCodes   = []DistrictCode{DistrictNorth, DistrictCentral, DistrictSouth}
)

const (
	CivicWorksStatusAssigned           CivicWorksStatus = "ASSIGNED"
	CivicWorksStatusInProgress         CivicWorksStatus = "IN_PROGRESS"
	CivicWorksStatusPartiallyCompleted CivicWorksStatus = "PARTIALLY_COMPLETED"
	CivicWorksStatusCompleted          CivicWorksStatus = "COMPLETED"
)

var CivicWorksStatuses = []CivicWorksStatus{
	CivicWorksStatusAssigned,
	CivicWorksStatusInProgress,
	CivicWorksStatusPartiallyCompleted,
	CivicWorksStatusCompleted,
}

const (
	RoleServiceAgent          ActorRole = "service_agent"
	RoleSupervisor            ActorRole = "supervisor"
	RoleDepartmentManager     ActorRole = "department_manager"
	RolePlatformAdministrator ActorRole = "platform_administrator"
	RoleWorkflowDesigner      ActorRole = "workflow_designer"
)

var ActorRoles = []ActorRole{
	RoleServiceAgent,
	RoleSupervisor,
	RoleDepartmentManager,
	RolePlatformAdministrator,
	RoleWorkflowDesigner,
}

const (
	ApplicationRolePublicVisitor         ApplicationRole = "public_visitor"
	ApplicationRoleConstituent           ApplicationRole = "constituent"
	ApplicationRoleServiceAgent          ApplicationRole = "service_agent"
	ApplicationRoleSupervisor            ApplicationRole = "supervisor"
	ApplicationRoleDepartmentManager     ApplicationRole = "department_manager"
	ApplicationRolePlatformAdministrator ApplicationRole = "platform_administrator"
	ApplicationRoleWorkflowDesigner      ApplicationRole = "workflow_designer"
	ApplicationRoleIntegrationClient     ApplicationRole = "integration_client"
)

var ApplicationRoles = []ApplicationRole{
	ApplicationRolePublicVisitor,
	ApplicationRoleConstituent,
	ApplicationRoleServiceAgent,
	ApplicationRoleSupervisor,
	ApplicationRoleDepartmentManager,
	ApplicationRolePlatformAdministrator,
	ApplicationRoleWorkflowDesigner,
	ApplicationRoleIntegrationClient,
}

const (
	AuditActorConstituent       AuditActorType = "constituent"
	AuditActorStaff             AuditActorType = "staff"
	AuditActorIntegrationClient AuditActorType = "integration_client"
	AuditActorSystem            AuditActorType = "system"
)

var AuditActorTypes = []AuditActorType{
	AuditActorConstituent,
	AuditActorStaff,
	AuditActorIntegrationClient,
	AuditActorSystem,
}

const (
	ValidationRequired            ValidationCode = "REQUIRED"
	ValidationInvalidFormat       ValidationCode = "INVALID_FORMAT"
	ValidationTooShort            ValidationCode = "TOO_SHORT"
	ValidationTooLong             ValidationCode = "TOO_LONG"
	ValidationOutOfRange          ValidationCode = "OUT_OF_RANGE"
	ValidationInvalidValue        ValidationCode = "INVALID_VALUE"
	ValidationInactiveValue       ValidationCode = "INACTIVE_VALUE"
	ValidationTooManyItems        ValidationCode = "TOO_MANY_ITEMS"
	ValidationDuplicate           ValidationCode = "DUPLICATE"
	ValidationConflict            ValidationCode = "CONFLICT"
	ValidationLocationRequired    ValidationCode = "LOCATION_REQUIRED"
	ValidationCoordinatesRequired ValidationCode = "COORDINATES_REQUIRED"
)

var ValidationCodes = []ValidationCode{
	ValidationRequired,
	ValidationInvalidFormat,
	ValidationTooShort,
	ValidationTooLong,
	ValidationOutOfRange,
	ValidationInvalidValue,
	ValidationInactiveValue,
	ValidationTooManyItems,
	ValidationDuplicate,
	ValidationConflict,
	ValidationLocationRequired,
	ValidationCoordinatesRequired,
}

const (
	ErrorUnauthenticated           ErrorCode = "UNAUTHENTICATED"
	ErrorForbidden                 ErrorCode = "FORBIDDEN"
	ErrorValidation                ErrorCode = "VALIDATION_ERROR"
	ErrorInvalidStatusTransition   ErrorCode = "INVALID_STATUS_TRANSITION"
	ErrorIdempotencyConflict       ErrorCode = "IDEMPOTENCY_CONFLICT"
	ErrorVersionConflict           ErrorCode = "VERSION_CONFLICT"
	ErrorInvalidFilter             ErrorCode = "INVALID_FILTER"
	ErrorInvalidPageToken          ErrorCode = "INVALID_PAGE_TOKEN"
	ErrorRateLimited               ErrorCode = "RATE_LIMITED"
	ErrorAddressNotFound           ErrorCode = "ADDRESS_NOT_FOUND"
	ErrorMapUnauthenticated        ErrorCode = "MAP_UNAUTHENTICATED"
	ErrorMapTemporarilyUnavailable ErrorCode = "MAP_TEMPORARILY_UNAVAILABLE"
	ErrorInvalidResetToken         ErrorCode = "INVALID_RESET_TOKEN"
	ErrorExpiredResetToken         ErrorCode = "EXPIRED_RESET_TOKEN"
	ErrorInsufficientScope         ErrorCode = "INSUFFICIENT_SCOPE"
	ErrorInvalidClient             ErrorCode = "INVALID_CLIENT"
	ErrorInvalidToken              ErrorCode = "INVALID_TOKEN"
	ErrorTemporarilyUnavailable    ErrorCode = "TEMPORARILY_UNAVAILABLE"
	ErrorNotFound                  ErrorCode = "NOT_FOUND"
	ErrorExpectedVersionRequired   ErrorCode = "EXPECTED_VERSION_REQUIRED"
	ErrorInvalidSignature          ErrorCode = "INVALID_SIGNATURE"
	ErrorOperationFailed           ErrorCode = "OPERATION_FAILED"
)

var ErrorCodes = []ErrorCode{
	ErrorUnauthenticated,
	ErrorForbidden,
	ErrorValidation,
	ErrorInvalidStatusTransition,
	ErrorIdempotencyConflict,
	ErrorVersionConflict,
	ErrorInvalidFilter,
	ErrorInvalidPageToken,
	ErrorRateLimited,
	ErrorAddressNotFound,
	ErrorMapUnauthenticated,
	ErrorMapTemporarilyUnavailable,
	ErrorInvalidResetToken,
	ErrorExpiredResetToken,
	ErrorInsufficientScope,
	ErrorInvalidClient,
	ErrorInvalidToken,
	ErrorTemporarilyUnavailable,
	ErrorNotFound,
	ErrorExpectedVersionRequired,
	ErrorInvalidSignature,
	ErrorOperationFailed,
}

const (
	ServiceRequestsPath       = "/api/v1/service-requests"
	ExportPathTemplate        = "/api/v1/export/{entity}"
	AnonymousStatusLookupPath = "/api/v1/public/service-request-status"
	PasswordResetRequestPath  = "/api/v1/auth/password-reset/request"
	PasswordResetConfirmPath  = "/api/v1/auth/password-reset/confirm"
	WorkflowActionsPath       = "/api/v1/actions"
	IdempotencyHeader         = "Idempotency-Key"
	IfMatchHeader             = "If-Match"
	RetryAfterHeader          = "Retry-After"
	ScopeRequestWrite         = "service_requests.write"
	ScopeCRMExport            = "crm.export"
	ScopeWorkflowExecute      = "workflow.execute"
)

type (
	PhoneNumber struct {
		Label PhoneLabel `json:"label"`
		Value string     `json:"value"`
	}

	Address struct {
		Line1      string `json:"line1"`
		Line2      string `json:"line2,omitempty"`
		City       string `json:"city"`
		Region     string `json:"region"`
		PostalCode string `json:"postal_code"`
		Country    string `json:"country"`
		Primary    bool   `json:"primary,omitempty"`
	}

	Constituent struct {
		ConstituentID     string          `json:"constituent_id"`
		DisplayName       string          `json:"display_name"`
		LoginIdentifier   string          `json:"login_identifier,omitempty"`
		Emails            []string        `json:"emails"`
		PhoneNumbers      []PhoneNumber   `json:"phone_numbers"`
		Addresses         []Address       `json:"addresses"`
		PrimaryCategory   ContactCategory `json:"primary_category"`
		PreferredLanguage Language        `json:"preferred_language"`
		EmailOptOut       bool            `json:"email_opt_out"`
	}

	RequesterInput struct {
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
		Phone       string `json:"phone,omitempty"`
	}

	LocationInput struct {
		Address   string   `json:"address"`
		Latitude  *float64 `json:"latitude,omitempty"`
		Longitude *float64 `json:"longitude,omitempty"`
	}

	ServiceRequestLocation struct {
		Address   Address  `json:"address"`
		Latitude  *float64 `json:"latitude,omitempty"`
		Longitude *float64 `json:"longitude,omitempty"`
	}

	AttachmentInput struct {
		Filename      string `json:"filename"`
		MediaType     string `json:"media_type"`
		ContentBase64 string `json:"content_base64"`
	}

	ServiceRequestCreate struct {
		Summary      string                 `json:"summary"`
		Description  string                 `json:"description"`
		ServiceType  ServiceType            `json:"service_type"`
		Requester    RequesterInput         `json:"requester"`
		Location     *LocationInput         `json:"location,omitempty"`
		Attachments  []AttachmentInput      `json:"attachments,omitempty"`
		CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
	}

	ServiceRequest struct {
		RequestID        string                  `json:"request_id"`
		RequestNumber    string                  `json:"request_number,omitempty"`
		Summary          string                  `json:"summary"`
		Description      string                  `json:"description"`
		ServiceType      ServiceType             `json:"service_type"`
		OwningDepartment DepartmentCode          `json:"owning_department"`
		CouncilDistrict  *DistrictCode           `json:"council_district,omitempty"`
		SourceChannel    SourceChannel           `json:"source_channel"`
		OriginClass      OriginClass             `json:"origin_class"`
		Status           ServiceRequestStatus    `json:"status"`
		PrimaryRequester Constituent             `json:"primary_requester"`
		Location         *ServiceRequestLocation `json:"location,omitempty"`
		CustomFields     map[string]interface{}  `json:"custom_fields"`
		DuplicateGroupID string                  `json:"duplicate_group_id,omitempty"`
		Version          uint64                  `json:"version"`
		CreatedAt        time.Time               `json:"created_at"`
		UpdatedAt        time.Time               `json:"updated_at"`
	}

	ResourceLinks struct {
		Self string `json:"self"`
	}

	ServiceRequestResponse struct {
		RequestID     string               `json:"request_id"`
		RequestNumber string               `json:"request_number"`
		Status        ServiceRequestStatus `json:"status"`
		Version       uint64               `json:"version"`
		CreatedAt     time.Time            `json:"created_at"`
		Links         ResourceLinks        `json:"links"`
	}

	FieldError struct {
		Field string         `json:"field"`
		Code  ValidationCode `json:"code"`
	}

	APIError struct {
		Error            ErrorCode    `json:"error"`
		Message          string       `json:"message"`
		Retryable        bool         `json:"retryable"`
		Errors           []FieldError `json:"errors,omitempty"`
		CurrentVersion   *uint64      `json:"current_version,omitempty"`
		FailingRequestID string       `json:"failing_request_id,omitempty"`
		OperationID      string       `json:"operation_id,omitempty"`
	}

	ExportResponse struct {
		Items         []map[string]interface{} `json:"items"`
		NextPageToken *string                  `json:"next_page_token"`
		GeneratedAt   time.Time                `json:"generated_at"`
	}

	Reminder struct {
		ReminderID       string          `json:"reminder_id"`
		RequestID        string          `json:"request_id"`
		Title            string          `json:"title"`
		DueAt            time.Time       `json:"due_at"`
		Timezone         string          `json:"timezone"`
		RecipientStaffID string          `json:"recipient_staff_id"`
		Channel          ReminderChannel `json:"channel"`
		Status           ReminderStatus  `json:"status"`
		CompletedAt      *time.Time      `json:"completed_at,omitempty"`
		CompletedBy      string          `json:"completed_by,omitempty"`
	}

	PublicHistoryItem struct {
		Action                string    `json:"action"`
		OccurredAt            time.Time `json:"occurred_at"`
		ResponsibleDepartment string    `json:"responsible_department"`
	}

	PublicServiceRequestDetail struct {
		RequestNumber    string               `json:"request_number"`
		Summary          string               `json:"summary"`
		ServiceType      ServiceType          `json:"service_type"`
		Status           ServiceRequestStatus `json:"status"`
		OwningDepartment DepartmentCode       `json:"owning_department"`
		UpdatedAt        time.Time            `json:"updated_at"`
		History          []PublicHistoryItem  `json:"history"`
	}

	AnonymousStatusLookupRequest struct {
		RequestNumber string `json:"request_number"`
		Email         string `json:"email"`
	}

	AnonymousStatusLookupResponse struct {
		RequestDetail *PublicServiceRequestDetail `json:"request_detail"`
	}

	PasswordResetRequest struct {
		Email string `json:"email"`
	}

	PasswordResetConfirm struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	PasswordResetResponse struct {
		Message string `json:"message"`
	}

	WorkflowExecution struct {
		ExecutionID      string    `json:"execution_id"`
		WorkflowVersion  uint64    `json:"workflow_version"`
		Trigger          string    `json:"trigger"`
		Outcome          string    `json:"outcome"`
		ActionsAttempted []string  `json:"actions_attempted"`
		Succeeded        bool      `json:"succeeded"`
		OccurredAt       time.Time `json:"occurred_at"`
		ResponseStatus   *int      `json:"response_status,omitempty"`
		Error            *APIError `json:"error,omitempty"`
	}

	GeocodeResponse struct {
		Address         string  `json:"address"`
		Latitude        float64 `json:"latitude"`
		Longitude       float64 `json:"longitude"`
		PrecisionDigits int     `json:"precision_digits"`
		Provider        string  `json:"provider"`
	}

	CivicWorksWorkOrder struct {
		WorkOrderID          string           `json:"work_order_id"`
		SourceCaseID         string           `json:"source_case_id"`
		ServiceRequestNumber string           `json:"service_request_number"`
		Status               CivicWorksStatus `json:"status"`
		ExternalStatusURL    string           `json:"external_status_url"`
		Version              uint64           `json:"version"`
		CreatedAt            time.Time        `json:"created_at"`
		UpdatedAt            time.Time        `json:"updated_at"`
	}

	CivicWorksEvent struct {
		EventID        string           `json:"event_id"`
		EventType      string           `json:"event_type"`
		WorkOrderID    string           `json:"work_order_id"`
		SourceCaseID   string           `json:"source_case_id"`
		PreviousStatus CivicWorksStatus `json:"previous_status"`
		Status         CivicWorksStatus `json:"status"`
		Version        uint64           `json:"version"`
		OccurredAt     time.Time        `json:"occurred_at"`
	}
)

func CanTransition(from, to ServiceRequestStatus) bool {
	for _, candidate := range ServiceRequestTransitions[from] {
		if candidate == to {
			return true
		}
	}
	return false
}

func MapCivicWorksStatus(status CivicWorksStatus) (ServiceRequestStatus, bool) {
	switch status {
	case CivicWorksStatusAssigned:
		return ServiceRequestStatusAssigned, true
	case CivicWorksStatusInProgress, CivicWorksStatusPartiallyCompleted:
		return ServiceRequestStatusInProgress, true
	case CivicWorksStatusCompleted:
		return ServiceRequestStatusResolved, true
	default:
		return "", false
	}
}

// PlanCivicWorksTransition returns the CRM lifecycle steps needed to absorb one
// CivicWorks status event without adding transitions forbidden by provision
// 10.1.2. A work order may legally complete directly from ASSIGNED, so that
// external edge is normalised atomically through IN_PROGRESS before RESOLVED.
func PlanCivicWorksTransition(current ServiceRequestStatus, status CivicWorksStatus) ([]ServiceRequestStatus, bool) {
	target, ok := MapCivicWorksStatus(status)
	if !ok {
		return nil, false
	}
	if current == target {
		return []ServiceRequestStatus{}, true
	}
	if status == CivicWorksStatusCompleted && (current == ServiceRequestStatusClosed || current == ServiceRequestStatusReopened) {
		return []ServiceRequestStatus{}, true
	}
	if current == ServiceRequestStatusAssigned && target == ServiceRequestStatusResolved {
		return []ServiceRequestStatus{ServiceRequestStatusInProgress, ServiceRequestStatusResolved}, true
	}
	if CanTransition(current, target) {
		return []ServiceRequestStatus{target}, true
	}
	return nil, false
}
