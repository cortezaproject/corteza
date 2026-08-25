package city311

// ContractDocument is the generated, language-neutral handoff consumed by the
// frontend developer. contract.json is checked semantically against this entire
// value, so additions, removals, and changed values drift in either direction.
type ContractDocument struct {
	ContractVersion           string                            `json:"contract_version"`
	Versioning                ContractVersioning                `json:"versioning"`
	Maintainer                string                            `json:"maintainer"`
	Provisions                []string                          `json:"provisions"`
	Decisions                 []ContractDecision                `json:"decisions"`
	Protocol                  map[string]interface{}            `json:"protocol"`
	Endpoints                 map[string]EndpointContract       `json:"endpoints"`
	Enums                     map[string][]string               `json:"enums"`
	ServiceTypeRules          map[string]ServiceTypeRule        `json:"service_type_rules"`
	HelpKeys                  map[string]string                 `json:"help_keys"`
	StatusTransitions         map[string][]string               `json:"status_transitions"`
	CivicWorksStatusMapping   map[string]string                 `json:"civicworks_status_mapping"`
	CivicWorksTransitionPlans map[string][]string               `json:"civicworks_transition_plans"`
	Schemas                   map[string]map[string]interface{} `json:"schemas"`
	Mocks                     map[string]MockContract           `json:"mocks"`
}

type ContractVersioning struct {
	Scheme                string `json:"scheme"`
	Stability             string `json:"stability"`
	EffectiveAt           string `json:"effective_at"`
	FirstPublishedVersion string `json:"first_published_version"`
	SupportedMajor        int    `json:"supported_major"`
	MajorVersionRule      string `json:"major_version_rule"`
}

type ContractDecision struct {
	ID         string   `json:"id"`
	Provisions []string `json:"provisions"`
	Decision   string   `json:"decision"`
}

type EndpointContract struct {
	Method                      string                            `json:"method"`
	Path                        string                            `json:"path"`
	Direction                   string                            `json:"direction"`
	Authentication              AuthenticationContract            `json:"authentication"`
	RequiredCapability          string                            `json:"required_capability,omitempty"`
	Scope                       string                            `json:"scope,omitempty"`
	RequiredHeaders             []string                          `json:"required_headers,omitempty"`
	PathParameters              map[string]map[string]interface{} `json:"path_parameters,omitempty"`
	RequestSchema               string                            `json:"request_schema,omitempty"`
	ResponseSchema              string                            `json:"response_schema,omitempty"`
	EntityResponseSchemas       map[string]string                 `json:"entity_response_schemas,omitempty"`
	SuccessStatuses             map[string]int                    `json:"success_statuses"`
	ErrorStatuses               map[string]int                    `json:"error_statuses,omitempty"`
	QueryParameters             map[string]map[string]interface{} `json:"query_parameters,omitempty"`
	Ordering                    string                            `json:"ordering,omitempty"`
	RateLimitPerClientPerMinute int                               `json:"rate_limit_per_client_per_minute,omitempty"`
	ResponseHeaders             map[string]string                 `json:"response_headers,omitempty"`
	PrivacyRule                 string                            `json:"privacy_rule,omitempty"`
}

type AuthenticationContract struct {
	Mode               string                     `json:"mode"`
	ActorClass         string                     `json:"actor_class,omitempty"`
	Credential         string                     `json:"credential,omitempty"`
	SignatureAlgorithm string                     `json:"signature_algorithm,omitempty"`
	Alternatives       []AuthorizationAlternative `json:"alternatives,omitempty"`
}

type AuthorizationAlternative struct {
	ApplicationRoles        []string `json:"application_roles,omitempty"`
	RecordScopeRequired     bool     `json:"record_scope_required,omitempty"`
	ResourceRelationship    string   `json:"resource_relationship,omitempty"`
	Permission              string   `json:"permission,omitempty"`
	AvailableActionRequired bool     `json:"available_action_required,omitempty"`
}

type ServiceTypeRule struct {
	Department                   string `json:"department"`
	LocationRequired             bool   `json:"location_required"`
	ConfirmedCoordinatesRequired bool   `json:"confirmed_coordinates_required"`
}

type MockContract struct {
	Endpoint   string            `json:"endpoint"`
	Role       string            `json:"role"`
	HTTPStatus int               `json:"http_status,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       interface{}       `json:"body"`
}

func NewContractDocument() ContractDocument {
	document := ContractDocument{
		ContractVersion: ContractVersion,
		Versioning: ContractVersioning{
			Scheme:                "semantic_versioning",
			Stability:             "frozen",
			EffectiveAt:           "merge_to_2024.9.x",
			FirstPublishedVersion: ContractVersion,
			SupportedMajor:        1,
			MajorVersionRule:      "after first publication, any consumer-incompatible contract revision increments the major version",
		},
		Maintainer: "Developer 1 - backend, integrations, and runtime",
		Provisions: []string{
			"4.3.1", "4.3.2", "6.6.3", "7.2.3", "7.3.5", "7.6.5", "7.6.6", "8.3.3",
			"9.1.1", "9.2.1", "9.3.1", "9.3.6", "9.3.7", "9.4.1", "9.4.2", "9.4.3", "9.5.1", "9.6.1", "9.6.2", "9.7.2",
			"10.1.1", "10.1.2", "10.1.4", "10.2.2", "10.2.3", "10.2.4", "10.5.2", "10.5.4", "10.7.2",
			"11.2.2", "11.2.3", "11.2.4", "11.3.4", "11.4.3", "11.4.4",
			"11.6.1", "11.6.2", "11.6.3", "11.6.4", "11.6.5", "11.6.6",
			"11.7.1", "11.7.2", "11.7.3", "11.7.4", "12.1.3", "12.3.1", "15.1.1", "15.2.1", "15.2.3",
		},
		Decisions: []ContractDecision{
			{
				ID:         "CW-COMPLETION-NORMALISATION",
				Provisions: []string{"10.1.2", "10.2.2", "10.2.3", "10.2.4"},
				Decision:   "CivicWorks ASSIGNED to COMPLETED is applied atomically as CRM ASSIGNED to IN_PROGRESS to RESOLVED; the CRM lifecycle is not widened",
			},
			{
				ID:         "PUBLIC-LOOKUP-PROJECTION",
				Provisions: []string{"8.3.3", "12.1.3"},
				Decision:   "POST /api/v1/public/service-request-status returns only the public request projection and returns the identical 404 body with request_detail=null for either credential mismatch",
			},
			{
				ID:         "ACCOUNT-API-PATHS",
				Provisions: []string{"10.5.2", "10.5.4"},
				Decision:   "the implementation paths are POST /api/v1/auth/password-reset/request and POST /api/v1/auth/password-reset/confirm because the specification fixes behavior but does not prescribe paths",
			},
			{
				ID:         "GENERIC-NOT-FOUND-CODE",
				Provisions: []string{"7.2.3", "12.3.1"},
				Decision:   "authenticated APIs use NOT_FOUND to distinguish HTTP 404 from FORBIDDEN; anonymous status lookup deliberately uses its privacy-safe empty projection instead of this error envelope",
			},
			{
				ID:         "LOCATION-REPRESENTATIONS",
				Provisions: []string{"9.2.1", "9.5.1", "11.6.5"},
				Decision:   "submission input uses a canonical address string with optional coordinates; the persisted record uses a structured address and stores coordinates to four decimal places",
			},
			{
				ID:         "IDENTITY-WIRE-VALUES",
				Provisions: []string{"11.3.4"},
				Decision:   "actor_role contains only identity-provider asserted staff roles and is bound to application_role through the published actor_role_mappings; constituent is actor_type, public visitors are unauthenticated, and integration clients are represented by OAuth scopes",
			},
		},
		Endpoints: map[string]EndpointContract{
			"service_request_create": {
				Method:          "POST",
				Path:            ServiceRequestsPath,
				Scope:           ScopeRequestWrite,
				RequiredHeaders: []string{"Authorization", "Content-Type", IdempotencyHeader},
				RequestSchema:   "service_request_create",
				ResponseSchema:  "service_request_response",
				SuccessStatuses: map[string]int{"created": 201, "equivalent_replay": 200},
				ErrorStatuses: map[string]int{
					string(ErrorUnauthenticated): 401, string(ErrorForbidden): 403, string(ErrorIdempotencyConflict): 409,
					string(ErrorVersionConflict): 409, string(ErrorValidation): 422,
				},
			},
			"data_export": {
				Method:         "GET",
				Path:           ExportPathTemplate,
				Scope:          ScopeCRMExport,
				ResponseSchema: "export_response",
				EntityResponseSchemas: map[string]string{
					"constituents": "constituent", "service-requests": "service_request_record",
					"audit-events": "audit_event", "follow-up-actions": "follow_up_action",
				},
				SuccessStatuses: map[string]int{"success": 200},
				ErrorStatuses: map[string]int{
					string(ErrorUnauthenticated): 401, string(ErrorForbidden): 403, string(ErrorInvalidFilter): 422,
					string(ErrorInvalidPageToken): 400, string(ErrorRateLimited): 429,
				},
				QueryParameters: map[string]map[string]interface{}{
					"page_size":     {"type": "integer", "minimum": 1, "maximum": 100, "default": 50},
					"page_token":    {"type": "string", "opaque": true},
					"updated_since": {"type": "string", "format": "date-time"},
					"filters":       {"type": "object", "entity_relevant": true},
				},
				Ordering:                    "stable ascending identifier order within each page",
				RateLimitPerClientPerMinute: 60,
				ResponseHeaders:             map[string]string{RetryAfterHeader: "required on HTTP 429; whole seconds until retry"},
			},
			"anonymous_status_lookup": {
				Method:          "POST",
				Path:            AnonymousStatusLookupPath,
				RequestSchema:   "anonymous_status_lookup_request",
				ResponseSchema:  "anonymous_status_lookup_response",
				SuccessStatuses: map[string]int{"found": 200, "generic_not_found": 404},
				PrivacyRule:     "request-number mismatch and email mismatch return the same status and body with request_detail=null",
			},
			"password_reset_request": {
				Method:          "POST",
				Path:            PasswordResetRequestPath,
				RequestSchema:   "password_reset_request",
				ResponseSchema:  "password_reset_response",
				SuccessStatuses: map[string]int{"accepted_for_every_email": 202},
				PrivacyRule:     "every submitted email returns the same public response",
			},
			"password_reset_confirm": {
				Method:          "POST",
				Path:            PasswordResetConfirmPath,
				RequestSchema:   "password_reset_confirm",
				ResponseSchema:  "password_reset_response",
				SuccessStatuses: map[string]int{"reset": 200},
				ErrorStatuses: map[string]int{
					string(ErrorInvalidResetToken): 422, string(ErrorExpiredResetToken): 422, string(ErrorValidation): 422,
				},
			},
			"workflow_action_execute": {
				Method:          "POST",
				Path:            WorkflowActionsPath,
				Scope:           ScopeWorkflowExecute,
				RequiredHeaders: []string{"Authorization", "Content-Type"},
				RequestSchema:   "workflow_action_request",
				ResponseSchema:  "workflow_action_accepted",
				SuccessStatuses: map[string]int{"accepted": 202},
				ErrorStatuses: map[string]int{
					string(ErrorInsufficientScope): 403, string(ErrorInvalidClient): 401,
					string(ErrorInvalidToken): 401, string(ErrorTemporarilyUnavailable): 503,
				},
			},
		},
		Enums: map[string][]string{
			"service_request_status": stringsOf(ServiceRequestStatuses),
			"source_channel":         stringsOf(SourceChannels),
			"origin_class":           stringsOf(OriginClasses),
			"relationship_type":      stringsOf(RelationshipTypes),
			"reminder_status":        stringsOf(ReminderStatuses),
			"reminder_channel":       stringsOf(ReminderChannels),
			"language":               stringsOf(Languages),
			"phone_label":            stringsOf(PhoneLabels),
			"custom_field_type":      stringsOf(CustomFieldTypes),
			"contact_category":       stringsOf(ContactCategories),
			"service_type":           stringsOf(ServiceTypes),
			"department_code":        stringsOf(DepartmentCodes),
			"district_code":          stringsOf(DistrictCodes),
			"civicworks_status":      stringsOf(CivicWorksStatuses),
			"actor_role":             stringsOf(ActorRoles),
			"actor_type":             {"constituent"},
			"error_code":             stringsOf(ErrorCodes),
		},
		StatusTransitions: map[string][]string{
			"DRAFT": {"SUBMITTED"}, "SUBMITTED": {"TRIAGED"}, "TRIAGED": {"ASSIGNED"},
			"ASSIGNED": {"IN_PROGRESS"}, "IN_PROGRESS": {"RESOLVED"},
			"RESOLVED": {"CLOSED", "REOPENED"}, "CLOSED": {"REOPENED"},
			"REOPENED": {"ASSIGNED", "IN_PROGRESS"},
		},
		CivicWorksStatusMapping: map[string]string{
			"ASSIGNED": "ASSIGNED", "IN_PROGRESS": "IN_PROGRESS", "PARTIALLY_COMPLETED": "IN_PROGRESS", "COMPLETED": "RESOLVED",
		},
		CivicWorksTransitionPlans: map[string][]string{
			"ASSIGNED+ASSIGNED": {}, "ASSIGNED+IN_PROGRESS": {"IN_PROGRESS"},
			"ASSIGNED+PARTIALLY_COMPLETED": {"IN_PROGRESS"}, "ASSIGNED+COMPLETED": {"IN_PROGRESS", "RESOLVED"},
			"IN_PROGRESS+IN_PROGRESS": {}, "IN_PROGRESS+PARTIALLY_COMPLETED": {}, "IN_PROGRESS+COMPLETED": {"RESOLVED"},
		},
		Schemas: contractSchemas(),
		Mocks:   contractMocks(),
	}
	assembleContract(&document)
	return document
}

func contractSchemas() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"structured_address": object([]string{"line1", "city", "region", "postal_code", "country"}, map[string]interface{}{
			"line1": stringProperty(1, 200), "line2": stringProperty(0, 200), "city": stringProperty(1, 120),
			"region": stringProperty(1, 120), "postal_code": stringProperty(1, 32), "country": stringProperty(2, 2),
			"primary": map[string]interface{}{"type": "boolean"},
		}),
		"phone_number": object([]string{"label", "value"}, map[string]interface{}{
			"label": map[string]interface{}{"enum_ref": "phone_label"}, "value": map[string]interface{}{"type": "string", "format": "e164"},
		}),
		"constituent": object([]string{"constituent_id", "display_name", "emails", "phone_numbers", "addresses", "primary_category", "preferred_language", "email_opt_out"}, map[string]interface{}{
			"constituent_id":     map[string]interface{}{"type": "string", "immutable": true},
			"display_name":       stringProperty(1, 120),
			"login_identifier":   map[string]interface{}{"type": "string", "min_length": 3, "max_length": 64, "pattern": "^[a-z0-9._-]+$", "case_insensitive_unique": true},
			"emails":             map[string]interface{}{"type": "array", "min_items_for_portal_account": 1, "items": map[string]interface{}{"type": "string", "format": "email"}, "verified_login_case_insensitive_unique": true},
			"phone_numbers":      map[string]interface{}{"type": "array", "max_items": 3, "items_ref": "phone_number"},
			"addresses":          map[string]interface{}{"type": "array", "max_items": 5, "items_ref": "structured_address", "maximum_primary_items": 1},
			"primary_category":   map[string]interface{}{"enum_ref": "contact_category", "must_be_active": true},
			"preferred_language": map[string]interface{}{"enum_ref": "language", "default": "EN"},
			"email_opt_out":      map[string]interface{}{"type": "boolean", "default": false},
			"custom_fields":      map[string]interface{}{"type": "object", "additional_properties": true},
		}),
		"requester": object([]string{"display_name", "email"}, map[string]interface{}{
			"display_name": stringProperty(1, 120), "email": map[string]interface{}{"type": "string", "format": "email"},
			"phone": map[string]interface{}{"type": "string", "format": "e164"},
		}),
		"location_input": object([]string{"address"}, map[string]interface{}{
			"address": map[string]interface{}{"type": "string"}, "latitude": coordinateProperty(-90, 90), "longitude": coordinateProperty(-180, 180),
		}),
		"service_request_location": object([]string{"address"}, map[string]interface{}{
			"address":  map[string]interface{}{"schema_ref": "structured_address"},
			"latitude": coordinateProperty(-90, 90), "longitude": coordinateProperty(-180, 180),
		}),
		"attachment": object([]string{"filename", "media_type", "content_base64"}, map[string]interface{}{
			"filename":       stringProperty(1, 120),
			"media_type":     map[string]interface{}{"enum": []string{"image/jpeg", "image/png", "application/pdf", "text/plain", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"}},
			"content_base64": map[string]interface{}{"type": "string", "content_encoding": "base64", "maximum_decoded_bytes": 10485760},
		}),
		"service_request_create": object([]string{"summary", "description", "service_type", "requester"}, map[string]interface{}{
			"summary": stringProperty(5, 160), "description": map[string]interface{}{"type": "string", "min_length": 10, "max_length": 5000, "content": "plain_text"},
			"service_type": map[string]interface{}{"enum_ref": "service_type", "must_be_active": true}, "requester": map[string]interface{}{"schema_ref": "requester"},
			"location": map[string]interface{}{"schema_ref": "location_input"}, "attachments": map[string]interface{}{"type": "array", "max_items": 5, "items_ref": "attachment", "validate_before_request_creation": true},
			"custom_fields": map[string]interface{}{"type": "object", "additional_properties": true, "keys_must_be_active": true},
		}),
		"service_request_response": object([]string{"request_id", "request_number", "status", "version", "created_at", "links"}, map[string]interface{}{
			"request_id": map[string]interface{}{"type": "string"}, "request_number": requestNumberProperty(), "status": map[string]interface{}{"const": "SUBMITTED"},
			"version": map[string]interface{}{"type": "integer", "minimum": 1}, "created_at": timestampProperty(),
			"links": object([]string{"self"}, map[string]interface{}{"self": map[string]interface{}{"type": "string"}}),
		}),
		"service_request_record": object([]string{"request_id", "summary", "description", "service_type", "owning_department", "source_channel", "origin_class", "status", "primary_requester", "version", "created_at", "updated_at"}, map[string]interface{}{
			"request_id": map[string]interface{}{"type": "string", "immutable": true}, "request_number": requestNumberProperty(),
			"summary": stringProperty(5, 160), "description": map[string]interface{}{"type": "string", "min_length": 10, "max_length": 5000, "content": "plain_text"},
			"service_type": map[string]interface{}{"enum_ref": "service_type"}, "owning_department": map[string]interface{}{"enum_ref": "department_code"},
			"council_district": map[string]interface{}{"enum_ref": "district_code"}, "source_channel": map[string]interface{}{"enum_ref": "source_channel"},
			"origin_class": map[string]interface{}{"enum_ref": "origin_class"}, "status": map[string]interface{}{"enum_ref": "service_request_status"},
			"primary_requester": map[string]interface{}{"schema_ref": "constituent"}, "location": map[string]interface{}{"schema_ref": "service_request_location"},
			"duplicate_group_id": map[string]interface{}{"type": "string"}, "version": map[string]interface{}{"type": "integer", "minimum": 1},
			"created_at": timestampProperty(), "updated_at": timestampProperty(), "custom_fields": map[string]interface{}{"type": "object", "additional_properties": true},
		}),
		"reminder": object([]string{"reminder_id", "request_id", "title", "due_at", "timezone", "recipient_staff_id", "channel", "status"}, map[string]interface{}{
			"reminder_id": map[string]interface{}{"type": "string"}, "request_id": map[string]interface{}{"type": "string"}, "title": map[string]interface{}{"type": "string"},
			"due_at": timestampProperty(), "timezone": map[string]interface{}{"type": "string"}, "recipient_staff_id": map[string]interface{}{"type": "string"},
			"channel": map[string]interface{}{"enum_ref": "reminder_channel"}, "status": map[string]interface{}{"enum_ref": "reminder_status"},
			"completed_at": timestampProperty(), "completed_by": map[string]interface{}{"type": "string"},
		}),
		"audit_event": object([]string{"entity_type", "entity_id", "event_type", "actor_type", "actor_id", "occurred_at", "source_channel", "before", "after"}, map[string]interface{}{
			"entity_type": map[string]interface{}{"type": "string"}, "entity_id": map[string]interface{}{"type": "string"}, "event_type": map[string]interface{}{"type": "string"},
			"actor_type": map[string]interface{}{"type": "string"}, "actor_id": map[string]interface{}{"type": "string"}, "occurred_at": timestampProperty(),
			"source_channel": map[string]interface{}{"enum_ref": "source_channel"}, "before": map[string]interface{}{"type": "object"}, "after": map[string]interface{}{"type": "object"},
		}),
		"follow_up_action": object([]string{"action_type", "actor", "occurred_at", "local_display_time", "request_id", "visibility", "payload"}, map[string]interface{}{
			"action_type": map[string]interface{}{"type": "string"}, "actor": map[string]interface{}{"type": "string"}, "occurred_at": timestampProperty(),
			"local_display_time": map[string]interface{}{"type": "string"}, "request_id": map[string]interface{}{"type": "string"},
			"visibility": map[string]interface{}{"type": "string"}, "payload": map[string]interface{}{"type": "object", "reportable": true},
		}),
		"public_history_item": object([]string{"action", "occurred_at", "responsible_department"}, map[string]interface{}{
			"action": map[string]interface{}{"type": "string"}, "occurred_at": timestampProperty(), "responsible_department": map[string]interface{}{"enum_ref": "department_code"},
		}),
		"public_service_request_detail": object([]string{"request_number", "summary", "service_type", "status", "owning_department", "updated_at", "history"}, map[string]interface{}{
			"request_number": requestNumberProperty(), "summary": stringProperty(5, 160), "service_type": map[string]interface{}{"enum_ref": "service_type"},
			"status": map[string]interface{}{"enum_ref": "service_request_status"}, "owning_department": map[string]interface{}{"enum_ref": "department_code"},
			"updated_at": timestampProperty(), "history": map[string]interface{}{"type": "array", "items_ref": "public_history_item"},
		}),
		"anonymous_status_lookup_request": object([]string{"request_number", "email"}, map[string]interface{}{
			"request_number": requestNumberProperty(), "email": map[string]interface{}{"type": "string", "format": "email"},
		}),
		"anonymous_status_lookup_response": object([]string{"request_detail"}, map[string]interface{}{
			"request_detail": map[string]interface{}{"schema_ref": "public_service_request_detail", "nullable": true},
		}),
		"password_reset_request": object([]string{"email"}, map[string]interface{}{"email": map[string]interface{}{"type": "string", "format": "email"}}),
		"password_reset_confirm": object([]string{"token", "password"}, map[string]interface{}{
			"token":    map[string]interface{}{"type": "string"},
			"password": map[string]interface{}{"type": "string", "min_length": 12, "max_length": 128, "minimum_character_classes": 3, "character_classes": []string{"uppercase", "lowercase", "digit", "symbol"}},
		}),
		"password_reset_response": object([]string{"message"}, map[string]interface{}{"message": map[string]interface{}{"type": "string"}}),
		"workflow_action_request": object([]string{"action", "request_id", "payload"}, map[string]interface{}{
			"action": map[string]interface{}{"type": "string"}, "request_id": map[string]interface{}{"type": "string"}, "payload": map[string]interface{}{"type": "object", "additional_properties": true},
		}),
		"workflow_action_accepted": object([]string{"execution_id", "accepted_at"}, map[string]interface{}{
			"execution_id": map[string]interface{}{"type": "string"}, "accepted_at": timestampProperty(),
		}),
		"workflow_execution": object([]string{"execution_id", "workflow_version", "trigger", "outcome", "actions_attempted", "succeeded", "occurred_at"}, map[string]interface{}{
			"execution_id": map[string]interface{}{"type": "string"}, "workflow_version": map[string]interface{}{"type": "integer", "minimum": 1}, "trigger": map[string]interface{}{"type": "string"},
			"outcome": map[string]interface{}{"type": "string"}, "actions_attempted": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
			"succeeded": map[string]interface{}{"type": "boolean"}, "occurred_at": timestampProperty(), "response_status": map[string]interface{}{"type": "integer", "minimum": 100, "maximum": 599},
			"error": map[string]interface{}{"schema_ref": "error"},
		}),
		"export_response": object([]string{"items", "next_page_token", "generated_at"}, map[string]interface{}{
			"items":           map[string]interface{}{"type": "array", "entity_schema_selected_by_path": true},
			"next_page_token": map[string]interface{}{"type": "string", "nullable": true, "opaque": true}, "generated_at": timestampProperty(),
		}),
		"error": object([]string{"error", "message", "retryable"}, map[string]interface{}{
			"error": map[string]interface{}{"enum_ref": "error_code"}, "message": map[string]interface{}{"type": "string"}, "retryable": map[string]interface{}{"type": "boolean"},
			"errors":          map[string]interface{}{"type": "array", "items": object([]string{"field", "code"}, map[string]interface{}{"field": map[string]interface{}{"type": "string"}, "code": map[string]interface{}{"type": "string"}})},
			"current_version": map[string]interface{}{"type": "integer", "minimum": 1},
		}),
	}
}

func contractMocks() map[string]MockContract {
	return map[string]MockContract{
		"service_request_created":     mock(201, MockCreatedServiceRequest()),
		"service_request_replay":      mock(200, MockCreatedServiceRequest()),
		"validation_error":            mock(422, MockValidationError()),
		"idempotency_conflict":        mock(409, MockIdempotencyConflict()),
		"unauthenticated":             mock(401, MockUnauthenticated()),
		"forbidden":                   mock(403, MockForbidden()),
		"not_found":                   mock(404, MockNotFound()),
		"version_conflict":            mock(409, MockVersionConflict(2)),
		"rate_limited":                {HTTPStatus: 429, Headers: map[string]string{RetryAfterHeader: "60"}, Body: MockRateLimited()},
		"invalid_reset_token":         mock(422, MockInvalidResetToken()),
		"expired_reset_token":         mock(422, MockExpiredResetToken()),
		"workflow_insufficient_scope": mock(403, MockWorkflowFailure(ErrorInsufficientScope, false)),
		"workflow_invalid_client":     mock(401, MockWorkflowFailure(ErrorInvalidClient, false)),
		"workflow_invalid_token":      mock(401, MockWorkflowFailure(ErrorInvalidToken, false)),
		"workflow_unavailable":        mock(503, MockWorkflowFailure(ErrorTemporarilyUnavailable, true)),
		"geocode_success":             mock(200, MockGeocodeSuccess()),
		"geocode_not_found":           mock(404, MockGeocodeNotFound()),
		"geocode_unavailable":         mock(503, MockGeocodeUnavailable()),
		"civicworks_created":          mock(201, MockCivicWorksCreated()),
		"civicworks_completed_event":  mock(200, MockCivicWorksCompletedEvent()),
		"anonymous_status_found":      mock(200, MockAnonymousStatusFound()),
		"anonymous_status_not_found":  mock(404, MockAnonymousStatusNotFound()),
		"password_reset_requested":    mock(202, MockPasswordResetRequested()),
	}
}

func mock(status int, body interface{}) MockContract {
	return MockContract{HTTPStatus: status, Body: body}
}

func object(required []string, properties map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{"type": "object", "required": required, "properties": properties}
}

func stringProperty(minimum, maximum int) map[string]interface{} {
	p := map[string]interface{}{"type": "string", "max_length": maximum}
	if minimum > 0 {
		p["min_length"] = minimum
	}
	return p
}

func coordinateProperty(minimum, maximum int) map[string]interface{} {
	return map[string]interface{}{"type": "number", "minimum": minimum, "maximum": maximum, "stored_decimal_places": 4}
}

func timestampProperty() map[string]interface{} {
	return map[string]interface{}{"type": "string", "format": "date-time", "explicit_utc_offset": true}
}

func requestNumberProperty() map[string]interface{} {
	return map[string]interface{}{"type": "string", "pattern": "^SR-[0-9]{4}-[0-9]{5}$"}
}

func stringsOf[T ~string](values []T) []string {
	out := make([]string, len(values))
	for i := range values {
		out[i] = string(values[i])
	}
	return out
}
