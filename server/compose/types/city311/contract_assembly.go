package city311

import "sort"

const (
	EndpointProvidedByCRM = "provided_by_crm"
	EndpointConsumedByCRM = "consumed_by_crm"
)

func assembleContract(document *ContractDocument) {
	addProvisions(document,
		"2.1.4", "3.3.2", "4.1.1", "4.3.3", "6.4.3", "6.5.4", "6.7.2", "6.7.3",
		"7.1.1", "7.1.3", "7.1.4", "7.1.5", "7.1.6", "7.3.1", "7.3.3", "7.3.4",
		"7.5.1", "7.5.2", "7.5.3", "7.5.4", "7.6.1", "7.6.2", "7.6.3", "7.6.4",
		"7.7.2", "7.7.3", "7.9.2", "7.9.3", "7.10.3", "7.11.1", "7.11.2", "7.11.3", "7.13.2",
		"8.2.1", "8.2.2", "8.2.3", "8.2.4", "9.2.4", "9.3.2", "9.3.5", "9.4.4", "9.5.3",
		"9.6.3", "9.7.3", "9.8.4", "9.8.6", "10.1.3", "10.1.5", "10.2.6",
		"10.3.1", "10.3.3", "10.3.4", "10.3.5", "10.4.1", "10.4.2", "10.4.3",
		"10.5.1", "10.5.6", "10.7.1", "11.3.1", "11.4.5", "11.5.3", "12.1.1", "12.2.1", "12.2.2", "12.2.3", "13.1.2", "13.1.3", "13.2.3", "15.4.1",
	)

	document.Decisions = append(document.Decisions,
		ContractDecision{
			ID:         "CLIENT-SERVER-PROTOCOL",
			Provisions: []string{"4.3.3", "6.4.3", "8.2.1", "8.2.2", "8.2.3", "8.2.4", "10.3.4", "10.3.5", "10.7.1", "10.7.2"},
			Decision:   "all browser APIs share the session, authorization, If-Match, JSON-Pointer validation, list, idempotency, asynchronous-operation, and atomic-bulk conventions in protocol",
		},
		ContractDecision{
			ID:         "APPLICATION-ROLE-VOCABULARY",
			Provisions: []string{"4.1.1", "4.3.3", "11.3.4"},
			Decision:   "application_role contains the eight permission-matrix roles; oidc_actor_type and actor_role remain the narrower identity-provider wire vocabularies; actor_role is bound by identity_configuration.actor_role_mappings; audit_actor_type is independently constrained",
		},
		ContractDecision{
			ID:         "PORTAL-ATTACHMENT-TRANSPORT",
			Provisions: []string{"7.6.1", "8.2.1", "9.4.1", "9.4.2", "9.4.3", "9.4.4"},
			Decision:   "the portal stages each attachment with a retryable upload and submits attachment tokens; the integration-client POST /api/v1/service-requests retains inline base64 for its published Chapter 11 contract",
		},
		ContractDecision{
			ID:         "CIVICWORKS-TERMINAL-REDELIVERY",
			Provisions: []string{"10.1.3", "10.2.4"},
			Decision:   "a repeated COMPLETED event is acknowledged as a no-op when CRM is RESOLVED, CLOSED, or REOPENED; direct completion from ASSIGNED creates two atomic audit and public-history entries, IN_PROGRESS then RESOLVED",
		},
		ContractDecision{
			ID:         "INTERNAL-API-PATHS",
			Provisions: []string{"7.1.1", "7.3.4", "7.6.1", "7.9.2", "7.9.3", "12.1.1", "12.2.1"},
			Decision:   "the specification fixes required behavior but not internal REST paths; the paths in endpoints are the shared implementation boundary for Developer 1 and Developer 2",
		},
		ContractDecision{
			ID:         "ACCOUNT-EXISTENCE-PRIVACY",
			Provisions: []string{"9.1.1", "10.5.2", "12.1.1"},
			Decision:   "local sign-in returns the same UNAUTHENTICATED response for an unknown identifier and an incorrect password; account registration returns the same accepted response for a new identifier and one already associated with a verified account",
		},
		ContractDecision{
			ID:         "IDENTITY-CONFIGURATION-BOUNDARY",
			Provisions: []string{"2.1.4", "11.3.1", "12.2.1", "13.2.3"},
			Decision:   "identity provider endpoints, client identifiers, role mappings, and secrets are supplied by runtime configuration; the administration view displays effective non-secret values and whether the OIDC secret is configured, while PATCH may enable or disable OIDC and SAML only and never reads or writes secret values",
		},
		ContractDecision{
			ID:         "FRONTEND-CAPABILITY-GATING",
			Provisions: []string{"4.1.1", "12.2.3"},
			Decision:   "each session-protected CRM operation publishes its stable operation name as required_capability; current_actor publishes the granted capability, route, and OAuth-scope vocabularies, while record-specific lifecycle permissions remain in available_actions and the server remains authoritative",
		},
	)

	document.Protocol = contractProtocol()
	applyContractEnums(document)
	document.ServiceTypeRules = map[string]ServiceTypeRule{
		"TREE_MAINTENANCE": {Department: "PUBLIC_WORKS", LocationRequired: true, ConfirmedCoordinatesRequired: true},
		"POTHOLE":          {Department: "STREETS", LocationRequired: true, ConfirmedCoordinatesRequired: true},
		"MISSED_TRASH":     {Department: "SANITATION", LocationRequired: true, ConfirmedCoordinatesRequired: true},
		"GENERAL_INQUIRY":  {Department: "GENERAL_SERVICES", LocationRequired: false, ConfirmedCoordinatesRequired: false},
	}
	document.HelpKeys = map[string]string{
		"public.request.submit":     "public_request_submission",
		"public.request.lookup":     "anonymous_status_lookup",
		"staff.request.triage":      "staff_triage",
		"staff.request.reassign":    "reassignment",
		"staff.request.bulk-update": "bulk_update",
		"admin.workflow.author":     "workflow_authoring",
		"staff.report.create":       "report_creation",
		"admin.branding.publish":    "branding_publication",
	}

	applyEndpointDirections(document)
	addContractEndpoints(document)
	applyCapabilityContract(document)
	mergeSchemas(document.Schemas, clientSchemas())
	applySchemaRules(document)
	mergeMocks(document.Mocks, clientMocks())
	linkMocks(document.Mocks)

	document.CivicWorksTransitionPlans["RESOLVED+COMPLETED"] = []string{}
	document.CivicWorksTransitionPlans["CLOSED+COMPLETED"] = []string{}
	document.CivicWorksTransitionPlans["REOPENED+COMPLETED"] = []string{}
}

func applySchemaRules(document *ContractDocument) {
	document.Schemas["service_request_create"]["conditional_validation_ref"] = "service_type_rules"
	auditProperties := document.Schemas["audit_event"]["properties"].(map[string]interface{})
	auditProperties["actor_type"] = map[string]interface{}{"enum_ref": "audit_actor_type"}
}

func contractProtocol() map[string]interface{} {
	return map[string]interface{}{
		"session_and_authorization": map[string]interface{}{
			"current_session_endpoint":       "session_current",
			"local_sign_in_endpoint":         "session_sign_in",
			"session_transport":              "Secure HttpOnly SameSite=Lax cookie",
			"anonymous_role":                 "public_visitor",
			"current_actor_schema":           "current_actor",
			"authentication_mode_enum_ref":   "authentication_mode",
			"actor_class_enum_ref":           "authenticated_actor_class",
			"resource_relationship_enum_ref": "resource_relationship",
			"capability_enum_ref":            "capability",
			"route_enum_ref":                 "route",
			"oauth_scope_enum_ref":           "oauth_scope",
			"optional_cookie_rule":           "an absent, expired, or invalid cookie is discarded and the request proceeds anonymously; optional-session endpoints do not return UNAUTHENTICATED solely because of the cookie",
			"authorization_rule":             "available_actions is computed from application roles plus department and district record scope; the server remains authoritative",
		},
		"optimistic_concurrency": map[string]interface{}{
			"request_header": IfMatchHeader,
			"syntax":         "quoted decimal resource version, for example \"7\"",
			"applies_to":     []string{"service_request", "workflow", "branding", "content", "help", "custom_field"},
			"missing_status": 428,
			"missing_error":  string(ErrorExpectedVersionRequired),
			"stale_status":   409,
			"stale_error":    string(ErrorVersionConflict),
		},
		"validation_errors": map[string]interface{}{
			"field_path_format": "RFC 6901 JSON Pointer rooted at the request body, including zero-based array indexes; examples: /requester/email and /attachments/2/media_type",
			"code_enum_ref":     "validation_code",
			"ordering":          "errors are returned in form order so the first item is the focus target",
			"translation_rule":  "validation code is a stable translation-catalogue key; the client supplies the localized message",
		},
		"lists": map[string]interface{}{
			"response_schema": "list_response",
			"page_size":       map[string]interface{}{"minimum": 1, "maximum": 100, "default": 50},
			"page_token":      "opaque",
			"sort_syntax":     "comma-separated field names; prefix descending fields with -; maximum three fields",
			"filter_echo":     true,
			"empty_total":     0,
		},
		"server_error_policy": map[string]interface{}{
			"declared_statuses":   []int{503},
			"undeclared_statuses": []int{500, 502, 504},
			"rule":                "only declared 5xx statuses have a stable City 311 error mapping; undeclared 5xx responses are terminal transport failures outside this contract",
		},
		"idempotency": map[string]interface{}{
			"request_header":    IdempotencyHeader,
			"required_for":      []string{"portal_service_request_submit", "service_request_create", "staff_request_bulk", "mail_send", "workflow_action_execute", "civicworks_work_order_create"},
			"equivalent_replay": "returns the original body without creating another request, work order, email, or workflow action",
			"different_content": map[string]interface{}{"status": 409, "error": string(ErrorIdempotencyConflict)},
		},
		"asynchronous_operations": map[string]interface{}{
			"accepted_status": 202,
			"response_schema": "operation",
			"poll_endpoint":   "operation_get",
			"terminal_states": []string{"SUCCEEDED", "FAILED", "CANCELLED"},
		},
		"atomic_bulk_operations": map[string]interface{}{
			"request_schema":         "bulk_request",
			"item_schema":            "bulk_request_item",
			"expected_version_field": "request_items[].expected_version",
			"success_schema":         "bulk_result",
			"failure_rule":           "validation or concurrency failure returns the generic error envelope with failing_request_id and changes no selected record",
		},
		"attachment_transport": map[string]interface{}{
			"portal":      "stage with portal_attachment_upload and submit returned attachment_token values",
			"integration": "inline base64 in service_request_create",
			"validation":  "filename, media type, count, and decoded size are validated before request creation",
		},
	}
}

func applyContractEnums(document *ContractDocument) {
	delete(document.Enums, "actor_type")
	document.Enums["oidc_actor_type"] = []string{"constituent"}
	document.Enums["application_role"] = stringsOf(ApplicationRoles)
	document.Enums["audit_actor_type"] = stringsOf(AuditActorTypes)
	document.Enums["validation_code"] = stringsOf(ValidationCodes)
	document.Enums["operation_status"] = []string{"PENDING", "RUNNING", "SUCCEEDED", "FAILED", "CANCELLED"}
	document.Enums["identity_provider"] = []string{"oidc", "saml"}
	document.Enums["public_content_key"] = []string{"HOME", "SERVICE_CATALOGUE", "HELP", "FOOTER", "TERMS"}
	document.Enums["export_entity"] = []string{"constituents", "service-requests", "audit-events", "follow-up-actions"}
	document.Enums["authentication_mode"] = []string{"none", "session_cookie_optional", "session_cookie", "oauth2_bearer", "oauth2_client_credentials", "server_api_token", "request_signature"}
	document.Enums["authenticated_actor_class"] = []string{"any_authenticated_actor", "constituent", "staff"}
	document.Enums["resource_relationship"] = []string{"linked_request_constituent", "resource_owner", "reminder_recipient"}
	document.Enums["request_action"] = []string{"TRIAGE", "ASSIGN", "START_PROGRESS", "RESOLVE", "CLOSE", "REQUEST_REOPEN", "APPROVE_REOPEN"}
	document.Enums["reminder_action"] = []string{"SNOOZE", "COMPLETE", "CANCEL"}
}

func applyEndpointDirections(document *ContractDocument) {
	for name, endpoint := range document.Endpoints {
		endpoint.Direction = EndpointProvidedByCRM
		endpoint.Authentication = authenticationContract("public")
		if endpoint.Scope != "" {
			endpoint.Authentication = authenticationContract("OAuth 2.0 bearer token")
		}
		if name == "workflow_action_execute" {
			endpoint.Direction = EndpointConsumedByCRM
			endpoint.Authentication = authenticationContract("OAuth 2.0 client credentials")
			endpoint.RequiredHeaders = appendUnique(endpoint.RequiredHeaders, IdempotencyHeader)
		}
		document.Endpoints[name] = endpoint
	}
}

func applyCapabilityContract(document *ContractDocument) {
	capabilities := make([]string, 0, len(document.Endpoints))
	routes := make([]string, 0, len(document.Endpoints))
	for name, endpoint := range document.Endpoints {
		if endpoint.Direction != EndpointProvidedByCRM {
			continue
		}
		routes = append(routes, name)
		if endpoint.Authentication.Mode != "session_cookie" {
			continue
		}
		endpoint.RequiredCapability = name
		document.Endpoints[name] = endpoint
		capabilities = append(capabilities, name)
	}
	sort.Strings(capabilities)
	sort.Strings(routes)
	document.Enums["capability"] = capabilities
	document.Enums["route"] = routes
	document.Enums["oauth_scope"] = []string{ScopeCRMExport, ScopeRequestWrite, ScopeWorkflowExecute}
}

func addProvisions(document *ContractDocument, provisions ...string) {
	seen := make(map[string]bool, len(document.Provisions)+len(provisions))
	for _, provision := range append(document.Provisions, provisions...) {
		seen[provision] = true
	}
	document.Provisions = document.Provisions[:0]
	for provision := range seen {
		document.Provisions = append(document.Provisions, provision)
	}
	sort.Strings(document.Provisions)
}

func appendUnique(values []string, additions ...string) []string {
	seen := make(map[string]bool, len(values)+len(additions))
	for _, value := range values {
		seen[value] = true
	}
	for _, value := range additions {
		if !seen[value] {
			values = append(values, value)
			seen[value] = true
		}
	}
	return values
}

func mergeSchemas(target, additions map[string]map[string]interface{}) {
	for name, schema := range additions {
		target[name] = schema
	}
}

func mergeMocks(target, additions map[string]MockContract) {
	for name, response := range additions {
		target[name] = response
	}
}
