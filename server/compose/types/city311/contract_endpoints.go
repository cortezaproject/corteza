package city311

type endpointSpec struct {
	Name           string
	Method         string
	Path           string
	Authentication string
	Request        string
	Response       string
	Status         int
}

func addContractEndpoints(document *ContractDocument) {
	for _, spec := range []endpointSpec{
		{"health", "GET", "/healthz", "public", "", "health", 200},
		{"session_current", "GET", "/api/v1/session", "session optional", "", "session", 200},
		{"session_sign_in", "POST", "/api/v1/session", "public", "local_sign_in", "session", 200},
		{"session_sign_out", "DELETE", "/api/v1/session", "session required", "", "empty_response", 204},
		{"account_register", "POST", "/api/v1/accounts", "public", "account_registration", "session", 201},
		{"federated_sign_in_start", "GET", "/api/v1/auth/{provider}/start", "public", "", "federated_redirect", 200},
		{"federated_sign_in_callback", "GET", "/api/v1/auth/{provider}/callback", "public", "", "session", 200},
		{"operation_get", "GET", "/api/v1/operations/{operation_id}", "session required", "", "operation", 200},
		{"geocode_proxy", "POST", "/api/v1/geocode", "session optional", "geocode_request", "geocode_response", 200},
		{"portal_attachment_upload", "POST", "/api/v1/portal/attachments", "session optional", "portal_attachment_upload", "portal_attachment", 201},
		{"attachment_download", "GET", "/api/v1/attachments/{attachment_id}", "request-view permission", "", "binary_attachment", 200},
		{"portal_service_request_submit", "POST", "/api/v1/portal/service-requests", "session optional", "portal_service_request_submit", "service_request_response", 201},
		{"portal_draft_create", "POST", "/api/v1/portal/service-request-drafts", "constituent session", "portal_draft_write", "service_request_record", 201},
		{"portal_draft_get", "GET", "/api/v1/portal/service-request-drafts/{request_id}", "constituent session", "", "service_request_record", 200},
		{"portal_draft_update", "PATCH", "/api/v1/portal/service-request-drafts/{request_id}", "constituent session", "portal_draft_write", "service_request_record", 200},
		{"portal_draft_delete", "DELETE", "/api/v1/portal/service-request-drafts/{request_id}", "constituent session", "", "empty_response", 204},
		{"portal_draft_submit", "POST", "/api/v1/portal/service-request-drafts/{request_id}/submit", "constituent session", "empty_request", "service_request_response", 200},
		{"portal_my_requests", "GET", "/api/v1/portal/service-requests", "constituent session", "", "list_response", 200},
		{"portal_link_anonymous_request", "POST", "/api/v1/portal/service-requests/link", "constituent session", "anonymous_request_link", "service_request_record", 200},
		{"profile_get", "GET", "/api/v1/account/profile", "constituent session", "", "constituent", 200},
		{"profile_update", "PATCH", "/api/v1/account/profile", "constituent session", "profile_update", "constituent", 200},
		{"password_change", "POST", "/api/v1/account/password", "constituent session", "password_change", "empty_response", 204},
		{"login_identifier_change", "POST", "/api/v1/account/login-identifier", "constituent session", "login_identifier_change", "session", 200},
		{"language_update", "PATCH", "/api/v1/preferences/language", "session optional", "language_preference", "language_preference", 200},
		{"public_branding_get", "GET", "/api/v1/public/branding", "public", "", "branding", 200},
		{"public_content_get", "GET", "/api/v1/public/content/{content_key}", "public", "", "content_object", 200},
		{"public_help_get", "GET", "/api/v1/public/help/{help_key}", "public", "", "help_content", 200},
		{"staff_request_queue", "GET", "/api/v1/staff/service-requests", "staff session and record scope", "", "list_response", 200},
		{"staff_request_detail", "GET", "/api/v1/staff/service-requests/{request_id}", "staff session and record scope", "", "staff_service_request_detail", 200},
		{"staff_service_request_create", "POST", "/api/v1/staff/service-requests", "service_agent", "staff_service_request_create", "staff_service_request_detail", 201},
		{"staff_request_transition", "POST", "/api/v1/staff/service-requests/{request_id}/transitions", "available action required", "request_transition", "staff_service_request_detail", 200},
		{"staff_request_reassign", "POST", "/api/v1/staff/service-requests/{request_id}/assignment", "supervisor or manager", "reassignment", "staff_service_request_detail", 200},
		{"staff_collaborator_add", "PUT", "/api/v1/staff/service-requests/{request_id}/collaborators/{staff_id}", "supervisor or manager", "reason", "staff_service_request_detail", 200},
		{"staff_collaborator_remove", "DELETE", "/api/v1/staff/service-requests/{request_id}/collaborators/{staff_id}", "supervisor or manager", "reason", "staff_service_request_detail", 200},
		{"staff_constituent_link", "POST", "/api/v1/staff/service-requests/{request_id}/constituents", "service_agent", "constituent_link", "staff_service_request_detail", 201},
		{"staff_constituent_unlink", "DELETE", "/api/v1/staff/service-requests/{request_id}/constituents/{constituent_id}", "service_agent", "constituent_unlink", "staff_service_request_detail", 200},
		{"staff_note_create", "POST", "/api/v1/staff/service-requests/{request_id}/notes", "service_agent", "request_note", "request_note", 201},
		{"staff_reminder_create", "POST", "/api/v1/staff/service-requests/{request_id}/reminders", "service_agent", "reminder_write", "reminder", 201},
		{"staff_reminder_action", "POST", "/api/v1/staff/reminders/{reminder_id}/{action}", "recipient staff or supervisor", "reminder_action", "reminder", 200},
		{"staff_request_bulk", "POST", "/api/v1/staff/service-requests/bulk", "supervisor or manager", "bulk_request", "bulk_result", 200},
		{"staff_duplicate_group_confirm", "POST", "/api/v1/staff/service-requests/{request_id}/duplicate-group", "supervisor", "duplicate_group_change", "staff_service_request_detail", 200},
		{"staff_duplicate_group_remove", "DELETE", "/api/v1/staff/service-requests/{request_id}/duplicate-group", "supervisor", "reason", "staff_service_request_detail", 200},
		{"portal_reopen_request", "POST", "/api/v1/portal/service-requests/{request_id}/reopen", "linked constituent", "reason", "reopen_request", 202},
		{"staff_reopen_approve", "POST", "/api/v1/staff/service-requests/{request_id}/reopen/approve", "supervisor or manager", "reason", "staff_service_request_detail", 200},
		{"staff_origin_override", "POST", "/api/v1/staff/service-requests/{request_id}/origin-class", "department_manager or platform_administrator", "origin_override", "staff_service_request_detail", 200},
		{"staff_scope_override", "POST", "/api/v1/staff/service-requests/{request_id}/scope-override", "department_manager or platform_administrator", "scope_override", "staff_service_request_detail", 200},
		{"staff_constituent_search", "GET", "/api/v1/staff/constituents", "staff session and record scope", "", "list_response", 200},
		{"staff_constituent_detail", "GET", "/api/v1/staff/constituents/{constituent_id}", "staff session and record scope", "", "constituent", 200},
		{"admin_branding_get", "GET", "/api/v1/admin/branding", "platform_administrator", "", "branding", 200},
		{"admin_branding_update", "PATCH", "/api/v1/admin/branding", "platform_administrator", "branding_write", "branding", 200},
		{"admin_branding_preview", "POST", "/api/v1/admin/branding/preview", "platform_administrator", "branding_write", "branding", 200},
		{"admin_branding_publish", "POST", "/api/v1/admin/branding/publish", "platform_administrator", "empty_request", "branding", 200},
		{"admin_branding_versions", "GET", "/api/v1/admin/branding/versions", "platform_administrator", "", "list_response", 200},
		{"admin_branding_rollback", "POST", "/api/v1/admin/branding/rollback", "platform_administrator", "rollback", "branding", 200},
		{"admin_content_list", "GET", "/api/v1/admin/content", "platform_administrator", "", "list_response", 200},
		{"admin_content_get", "GET", "/api/v1/admin/content/{content_key}", "platform_administrator", "", "content_object", 200},
		{"admin_content_update", "PATCH", "/api/v1/admin/content/{content_key}", "platform_administrator", "content_write", "content_object", 200},
		{"admin_content_preview", "POST", "/api/v1/admin/content/{content_key}/preview", "platform_administrator", "content_write", "content_object", 200},
		{"admin_content_publish", "POST", "/api/v1/admin/content/{content_key}/publish", "platform_administrator", "empty_request", "content_object", 200},
		{"admin_content_versions", "GET", "/api/v1/admin/content/{content_key}/versions", "platform_administrator", "", "list_response", 200},
		{"admin_content_rollback", "POST", "/api/v1/admin/content/{content_key}/rollback", "platform_administrator", "rollback", "content_object", 200},
		{"admin_help_update", "PATCH", "/api/v1/admin/help/{help_key}", "platform_administrator", "help_write", "help_content", 200},
		{"admin_categories_list", "GET", "/api/v1/admin/contact-categories", "department_manager or platform_administrator", "", "list_response", 200},
		{"admin_categories_create", "POST", "/api/v1/admin/contact-categories", "department_manager or platform_administrator", "category_write", "category", 201},
		{"admin_categories_update", "PATCH", "/api/v1/admin/contact-categories/{category_code}", "department_manager or platform_administrator", "category_write", "category", 200},
		{"admin_custom_fields_list", "GET", "/api/v1/admin/custom-fields", "platform_administrator", "", "list_response", 200},
		{"admin_custom_fields_create", "POST", "/api/v1/admin/custom-fields", "platform_administrator", "custom_field_definition", "custom_field_definition", 201},
		{"admin_custom_fields_update", "PATCH", "/api/v1/admin/custom-fields/{field_key}", "platform_administrator", "custom_field_definition", "custom_field_definition", 200},
		{"workflow_list", "GET", "/api/v1/admin/workflows", "workflow_designer", "", "list_response", 200},
		{"workflow_create", "POST", "/api/v1/admin/workflows", "workflow_designer", "workflow_definition", "workflow_definition", 201},
		{"workflow_get", "GET", "/api/v1/admin/workflows/{workflow_id}", "workflow_designer", "", "workflow_definition", 200},
		{"workflow_update", "PATCH", "/api/v1/admin/workflows/{workflow_id}", "workflow_designer", "workflow_definition", "workflow_definition", 200},
		{"workflow_activate", "POST", "/api/v1/admin/workflows/{workflow_id}/activate", "workflow_designer", "empty_request", "workflow_definition", 200},
		{"workflow_deactivate", "POST", "/api/v1/admin/workflows/{workflow_id}/deactivate", "workflow_designer", "empty_request", "workflow_definition", 200},
		{"workflow_test", "POST", "/api/v1/admin/workflows/{workflow_id}/test", "workflow_designer", "workflow_test", "operation", 202},
		{"workflow_execution_list", "GET", "/api/v1/admin/workflow-executions", "workflow_designer", "", "list_response", 200},
		{"workflow_execution_get", "GET", "/api/v1/admin/workflow-executions/{execution_id}", "workflow_designer", "", "workflow_execution", 200},
		{"identity_configuration_get", "GET", "/api/v1/admin/identity", "platform_administrator", "", "identity_configuration", 200},
		{"identity_configuration_update", "PATCH", "/api/v1/admin/identity", "platform_administrator", "identity_configuration_write", "identity_configuration", 200},
		{"integration_list", "GET", "/api/v1/admin/integrations", "platform_administrator", "", "list_response", 200},
		{"integration_get", "GET", "/api/v1/admin/integrations/{integration_id}", "platform_administrator", "", "integration_connection", 200},
		{"integration_update", "PATCH", "/api/v1/admin/integrations/{integration_id}", "platform_administrator", "integration_connection_write", "integration_connection", 200},
		{"integration_rotate", "POST", "/api/v1/admin/integrations/{integration_id}/rotate", "platform_administrator", "secret_rotation", "integration_connection", 200},
		{"integration_revoke", "POST", "/api/v1/admin/integrations/{integration_id}/revoke", "platform_administrator", "reason", "integration_connection", 200},
		{"report_catalogue", "GET", "/api/v1/staff/reports/catalogue", "staff session", "", "list_response", 200},
		{"report_run", "POST", "/api/v1/staff/reports/run", "staff session and record scope", "report_run", "operation", 202},
		{"saved_report_list", "GET", "/api/v1/staff/reports", "staff session", "", "list_response", 200},
		{"saved_report_create", "POST", "/api/v1/staff/reports", "staff session", "report_definition", "report_definition", 201},
		{"saved_report_update", "PATCH", "/api/v1/staff/reports/{report_id}", "owner", "report_definition", "report_definition", 200},
		{"saved_report_share", "POST", "/api/v1/staff/reports/{report_id}/share", "owner", "report_share", "report_definition", 200},
		{"report_export", "POST", "/api/v1/staff/reports/{report_id}/export", "staff session and record scope", "report_export", "operation", 202},
		{"contact_email_export", "POST", "/api/v1/staff/contact-email-export", "department_manager or platform_administrator", "contact_email_export", "operation", 202},
		{"audit_list", "GET", "/api/v1/staff/audit-events", "department_manager or platform_administrator", "", "list_response", 200},
		{"audit_export", "POST", "/api/v1/staff/audit-events/export", "department_manager or platform_administrator", "audit_export", "operation", 202},
		{"mail_preview", "POST", "/api/v1/staff/mail/preview", "authorised staff", "mail_compose", "mail_preview", 200},
		{"mail_send", "POST", "/api/v1/staff/mail", "authorised staff", "mail_compose", "mail_delivery", 202},
		{"mail_delivery_get", "GET", "/api/v1/staff/mail/{delivery_id}", "authorised staff", "", "mail_delivery", 200},
		{"calendar_import", "POST", "/api/v1/staff/calendar/import", "staff session", "calendar_import", "operation", 202},
		{"calendar_export", "GET", "/api/v1/staff/calendar/export", "staff session", "", "calendar_export", 200},
	} {
		document.Endpoints[spec.Name] = EndpointContract{
			Method:          spec.Method,
			Path:            spec.Path,
			Direction:       EndpointProvidedByCRM,
			Authentication:  authenticationContract(spec.Authentication),
			RequestSchema:   spec.Request,
			ResponseSchema:  spec.Response,
			SuccessStatuses: map[string]int{"success": spec.Status},
			ErrorStatuses:   standardBrowserErrors(),
		}
	}

	addConsumedEndpoints(document)
	configureEndpointMechanics(document)
}

func addConsumedEndpoints(document *ContractDocument) {
	document.Endpoints["mapping_geocode"] = EndpointContract{
		Method: "POST", Path: "/api/v1/geocode", Direction: EndpointConsumedByCRM,
		Authentication: authenticationContract("MAP_API_TOKEN supplied only to server runtime"), RequiredHeaders: []string{"Authorization", "Content-Type"},
		RequestSchema: "geocode_request", ResponseSchema: "geocode_response", SuccessStatuses: map[string]int{"success": 200},
		ErrorStatuses: map[string]int{string(ErrorAddressNotFound): 404, string(ErrorMapUnauthenticated): 401, string(ErrorMapTemporarilyUnavailable): 503},
	}
	document.Endpoints["civicworks_work_order_create"] = EndpointContract{
		Method: "POST", Path: "/api/v1/work-orders", Direction: EndpointConsumedByCRM,
		Authentication:  authenticationContract("CIVICWORKS_API_TOKEN supplied only to server runtime"),
		RequiredHeaders: []string{"Authorization", "X-Benchmark-Run-Id", "Content-Type", IdempotencyHeader},
		RequestSchema:   "civicworks_work_order_create", ResponseSchema: "civicworks_work_order", SuccessStatuses: map[string]int{"created": 201, "equivalent_replay": 200},
		ErrorStatuses: map[string]int{string(ErrorIdempotencyConflict): 409, string(ErrorValidation): 422, string(ErrorTemporarilyUnavailable): 503},
	}
	document.Endpoints["civicworks_event_callback"] = EndpointContract{
		Method: "POST", Path: "/integrations/civicworks/events", Direction: EndpointProvidedByCRM,
		Authentication:  authenticationContract("HMAC-SHA256 signature over exact request body"),
		RequiredHeaders: []string{"Content-Type", "X-CivicWorks-Event-Id", "X-CivicWorks-Signature"},
		RequestSchema:   "civicworks_event", ResponseSchema: "event_acknowledgement", SuccessStatuses: map[string]int{"acknowledged": 204, "duplicate_acknowledged": 204},
		ErrorStatuses: map[string]int{string(ErrorInvalidSignature): 401, string(ErrorValidation): 422},
	}
}

func configureEndpointMechanics(document *ContractDocument) {
	for _, name := range []string{
		"portal_draft_update", "portal_draft_delete", "portal_draft_submit", "profile_update",
		"staff_request_transition", "staff_request_reassign", "staff_collaborator_add", "staff_collaborator_remove",
		"staff_constituent_link", "staff_constituent_unlink", "staff_duplicate_group_confirm", "staff_duplicate_group_remove",
		"staff_reopen_approve", "staff_origin_override", "staff_scope_override", "admin_branding_update", "admin_branding_publish",
		"admin_branding_rollback", "admin_content_update", "admin_content_publish", "admin_content_rollback", "admin_custom_fields_update",
		"admin_help_update",
		"workflow_update", "workflow_activate", "workflow_deactivate", "identity_configuration_update", "integration_update", "integration_rotate",
		"integration_revoke", "saved_report_update", "saved_report_share",
	} {
		endpoint := document.Endpoints[name]
		endpoint.RequiredHeaders = appendUnique(endpoint.RequiredHeaders, IfMatchHeader)
		endpoint.ErrorStatuses[string(ErrorExpectedVersionRequired)] = 428
		endpoint.ErrorStatuses[string(ErrorVersionConflict)] = 409
		document.Endpoints[name] = endpoint
	}

	for _, name := range []string{"portal_service_request_submit", "staff_request_bulk", "mail_send", "workflow_action_execute", "civicworks_work_order_create"} {
		endpoint := document.Endpoints[name]
		endpoint.RequiredHeaders = appendUnique(endpoint.RequiredHeaders, IdempotencyHeader)
		endpoint.ErrorStatuses[string(ErrorIdempotencyConflict)] = 409
		document.Endpoints[name] = endpoint
	}

	bulk := document.Endpoints["staff_request_bulk"]
	bulk.ErrorStatuses[string(ErrorVersionConflict)] = 409
	document.Endpoints["staff_request_bulk"] = bulk

	configurePublicEndpointErrors(document)
	configureOptionalSessionEndpointErrors(document)
	bindPathParameters(document)

	for name, itemSchema := range map[string]string{
		"portal_my_requests": "portal_request_summary", "staff_request_queue": "request_queue_item",
		"staff_constituent_search": "constituent", "admin_branding_versions": "branding",
		"admin_content_list": "content_object", "admin_content_versions": "content_object",
		"admin_categories_list": "category", "admin_custom_fields_list": "custom_field_definition",
		"workflow_list": "workflow_definition", "workflow_execution_list": "workflow_execution",
		"integration_list": "integration_connection", "report_catalogue": "report_catalogue_item",
		"saved_report_list": "report_definition", "audit_list": "audit_event",
	} {
		endpoint := document.Endpoints[name]
		endpoint.EntityResponseSchemas = map[string]string{"items": itemSchema}
		endpoint.QueryParameters = standardListQuery()
		document.Endpoints[name] = endpoint
	}

	queue := document.Endpoints["staff_request_queue"]
	queue.QueryParameters["filters"] = map[string]interface{}{
		"type": "object", "fields": []string{"status", "service_type", "department", "district", "origin_class", "source_channel", "assignee", "collaborator", "category", "created_from", "created_to", "duplicate_group"},
	}
	document.Endpoints["staff_request_queue"] = queue

	for _, name := range []string{"workflow_test", "report_run", "report_export", "contact_email_export", "audit_export", "mail_send", "calendar_import"} {
		endpoint := document.Endpoints[name]
		endpoint.SuccessStatuses = map[string]int{"accepted": 202}
		document.Endpoints[name] = endpoint
	}
}

func configureOptionalSessionEndpointErrors(document *ContractDocument) {
	for name, errors := range map[string]map[string]int{
		"session_current":               {},
		"language_update":               {string(ErrorValidation): 422},
		"portal_attachment_upload":      {string(ErrorValidation): 422},
		"portal_service_request_submit": {string(ErrorIdempotencyConflict): 409, string(ErrorValidation): 422},
		"geocode_proxy": {
			string(ErrorAddressNotFound):           404,
			string(ErrorMapTemporarilyUnavailable): 503,
			string(ErrorValidation):                422,
		},
	} {
		endpoint := document.Endpoints[name]
		endpoint.ErrorStatuses = errors
		document.Endpoints[name] = endpoint
	}
}

func configurePublicEndpointErrors(document *ContractDocument) {
	setErrors := func(name string, errors map[string]int) EndpointContract {
		endpoint := document.Endpoints[name]
		endpoint.ErrorStatuses = errors
		document.Endpoints[name] = endpoint
		return endpoint
	}

	setErrors("health", map[string]int{string(ErrorTemporarilyUnavailable): 503})
	setErrors("public_branding_get", map[string]int{})
	setErrors("public_content_get", map[string]int{string(ErrorNotFound): 404})
	setErrors("public_help_get", map[string]int{string(ErrorNotFound): 404})
	setErrors("federated_sign_in_start", map[string]int{})
	setErrors("federated_sign_in_callback", map[string]int{
		string(ErrorUnauthenticated):        401,
		string(ErrorTemporarilyUnavailable): 503,
	})

	signIn := setErrors("session_sign_in", map[string]int{
		string(ErrorUnauthenticated): 401,
		string(ErrorValidation):      422,
	})
	signIn.PrivacyRule = "an unknown login identifier and an incorrect password return the identical status and generic error body"
	document.Endpoints["session_sign_in"] = signIn

	registration := setErrors("account_register", map[string]int{string(ErrorValidation): 422})
	registration.ResponseSchema = "account_registration_acknowledgement"
	registration.SuccessStatuses = map[string]int{"accepted": 202, "existing_account_indistinguishable": 202}
	registration.PrivacyRule = "a syntactically valid new identifier and one already associated with a verified account return the identical status and acknowledgement body"
	document.Endpoints["account_register"] = registration
}

func authenticationContract(profile string) AuthenticationContract {
	staffRoles := []string{"service_agent", "supervisor", "department_manager", "platform_administrator", "workflow_designer"}
	staff := func(roles ...string) []AuthorizationAlternative {
		if len(roles) == 0 {
			roles = staffRoles
		}
		return []AuthorizationAlternative{{ApplicationRoles: roles}}
	}
	switch profile {
	case "public":
		return AuthenticationContract{Mode: "none"}
	case "session optional":
		return AuthenticationContract{Mode: "session_cookie_optional"}
	case "session required":
		return AuthenticationContract{Mode: "session_cookie", ActorClass: "any_authenticated_actor"}
	case "constituent session":
		return AuthenticationContract{Mode: "session_cookie", ActorClass: "constituent"}
	case "staff session", "authorised staff":
		return AuthenticationContract{Mode: "session_cookie", ActorClass: "staff", Alternatives: staff()}
	case "staff session and record scope":
		return AuthenticationContract{Mode: "session_cookie", ActorClass: "staff", Alternatives: []AuthorizationAlternative{{ApplicationRoles: staffRoles, RecordScopeRequired: true}}}
	case "service_agent":
		return AuthenticationContract{Mode: "session_cookie", ActorClass: "staff", Alternatives: staff("service_agent")}
	case "supervisor":
		return AuthenticationContract{Mode: "session_cookie", ActorClass: "staff", Alternatives: staff("supervisor")}
	case "supervisor or manager":
		return AuthenticationContract{Mode: "session_cookie", ActorClass: "staff", Alternatives: staff("supervisor", "department_manager")}
	case "department_manager or platform_administrator":
		return AuthenticationContract{Mode: "session_cookie", ActorClass: "staff", Alternatives: staff("department_manager", "platform_administrator")}
	case "platform_administrator":
		return AuthenticationContract{Mode: "session_cookie", ActorClass: "staff", Alternatives: staff("platform_administrator")}
	case "workflow_designer":
		return AuthenticationContract{Mode: "session_cookie", ActorClass: "staff", Alternatives: staff("workflow_designer")}
	case "request-view permission":
		return AuthenticationContract{Mode: "session_cookie", ActorClass: "any_authenticated_actor", Alternatives: []AuthorizationAlternative{{Permission: "request_view"}}}
	case "available action required":
		return AuthenticationContract{Mode: "session_cookie", ActorClass: "staff", Alternatives: []AuthorizationAlternative{{ApplicationRoles: staffRoles, AvailableActionRequired: true}}}
	case "linked constituent":
		return AuthenticationContract{Mode: "session_cookie", ActorClass: "constituent", Alternatives: []AuthorizationAlternative{{ResourceRelationship: "linked_request_constituent"}}}
	case "owner":
		return AuthenticationContract{Mode: "session_cookie", ActorClass: "staff", Alternatives: []AuthorizationAlternative{{ResourceRelationship: "resource_owner"}}}
	case "recipient staff or supervisor":
		return AuthenticationContract{Mode: "session_cookie", ActorClass: "staff", Alternatives: []AuthorizationAlternative{{ResourceRelationship: "reminder_recipient"}, {ApplicationRoles: []string{"supervisor"}}}}
	case "OAuth 2.0 bearer token":
		return AuthenticationContract{Mode: "oauth2_bearer"}
	case "OAuth 2.0 client credentials":
		return AuthenticationContract{Mode: "oauth2_client_credentials"}
	case "MAP_API_TOKEN supplied only to server runtime":
		return AuthenticationContract{Mode: "server_api_token", Credential: "MAP_API_TOKEN"}
	case "CIVICWORKS_API_TOKEN supplied only to server runtime":
		return AuthenticationContract{Mode: "server_api_token", Credential: "CIVICWORKS_API_TOKEN"}
	case "HMAC-SHA256 signature over exact request body":
		return AuthenticationContract{Mode: "request_signature", Credential: "X-CivicWorks-Signature", SignatureAlgorithm: "HMAC-SHA256 over exact request body"}
	default:
		panic("unsupported authentication profile: " + profile)
	}
}

func bindPathParameters(document *ContractDocument) {
	for name, endpoint := range document.Endpoints {
		parameters := pathParameterNames(endpoint.Path)
		if len(parameters) == 0 {
			continue
		}
		endpoint.PathParameters = make(map[string]map[string]interface{}, len(parameters))
		for _, parameter := range parameters {
			property := map[string]interface{}{"type": "string"}
			switch parameter {
			case "provider":
				property["enum_ref"] = "identity_provider"
			case "content_key":
				property["enum_ref"] = "public_content_key"
			case "help_key":
				property["enum_source"] = "help_keys"
			case "action":
				property["enum_ref"] = "reminder_action"
			case "entity":
				property["enum_ref"] = "export_entity"
			}
			endpoint.PathParameters[parameter] = property
		}
		document.Endpoints[name] = endpoint
	}
}

func pathParameterNames(path string) []string {
	var parameters []string
	for start := 0; start < len(path); {
		open := indexByteFrom(path, '{', start)
		if open < 0 {
			break
		}
		close := indexByteFrom(path, '}', open+1)
		if close < 0 {
			panic("unclosed path parameter in " + path)
		}
		parameters = append(parameters, path[open+1:close])
		start = close + 1
	}
	return parameters
}

func indexByteFrom(value string, target byte, start int) int {
	for index := start; index < len(value); index++ {
		if value[index] == target {
			return index
		}
	}
	return -1
}

func standardBrowserErrors() map[string]int {
	return map[string]int{
		string(ErrorUnauthenticated): 401,
		string(ErrorForbidden):       403,
		string(ErrorNotFound):        404,
		string(ErrorValidation):      422,
	}
}

func standardListQuery() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"page_size":  {"type": "integer", "minimum": 1, "maximum": 100, "default": 50},
		"page_token": {"type": "string", "opaque": true},
		"sort":       {"type": "string", "maximum_fields": 3},
		"filters":    {"type": "object"},
	}
}
