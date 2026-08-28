package city311

func clientMocks() map[string]MockContract {
	return map[string]MockContract{
		"session_anonymous": mock(200, map[string]interface{}{
			"authenticated": false, "actor": nil, "preferred_language": "EN", "expires_at": nil,
		}),
		"session_staff": mock(200, map[string]interface{}{
			"authenticated": true,
			"actor": map[string]interface{}{
				"actor_id": "staff-0007", "display_name": "Sam Rivera", "application_roles": []string{"service_agent"},
				"department_codes": []string{"STREETS"}, "district_codes": []string{"NORTH", "CENTRAL"}, "scopes": []string{},
				"capabilities":     []string{"staff_request_detail", "staff_request_queue", "staff_request_transition"},
				"available_routes": []string{"staff_request_queue", "staff_request_detail"},
			},
			"preferred_language": "EN", "expires_at": "2026-08-25T20:00:00Z",
		}),
		"local_sign_in_rejected": mock(401, APIError{
			Error: ErrorUnauthenticated, Message: "The login identifier or password is incorrect.", Retryable: false,
		}),
		"account_registration_accepted": mock(202, map[string]interface{}{"accepted": true}),
		"health_unavailable": mock(503, APIError{
			Error: ErrorTemporarilyUnavailable, Message: "A required database connection or migration is unavailable.", Retryable: true,
		}),
		"operation_pending": mock(200, map[string]interface{}{
			"operation_id": "op-00041", "kind": "REPORT_EXPORT", "status": "RUNNING", "progress": 40,
			"result": nil, "error": nil, "created_at": "2026-08-25T12:00:00Z", "updated_at": "2026-08-25T12:00:04Z", "completed_at": nil,
		}),
		"operation_succeeded": mock(200, map[string]interface{}{
			"operation_id": "op-00041", "kind": "REPORT_EXPORT", "status": "SUCCEEDED", "progress": 100,
			"result": map[string]interface{}{"download_url": "/api/v1/operations/op-00041/result"}, "error": nil,
			"created_at": "2026-08-25T12:00:00Z", "updated_at": "2026-08-25T12:00:06Z", "completed_at": "2026-08-25T12:00:06Z",
		}),
		"empty_list": mock(200, map[string]interface{}{
			"items": []interface{}{}, "next_page_token": nil, "total_count": 0,
			"applied_filters": map[string]interface{}{"status": []string{"CLOSED"}}, "sort": []string{"-updated_at"},
		}),
		"nested_validation_error": mock(422, APIError{
			Error: ErrorValidation, Message: "The request contains invalid fields.", Retryable: false,
			Errors: []FieldError{{Field: "/attachments/2/media_type", Code: ValidationInvalidValue}},
		}),
		"bulk_validation_failure": mock(422, APIError{
			Error: ErrorValidation, Message: "One selected request is not eligible for this bulk operation.", Retryable: false,
			FailingRequestID: "case-7c58d2", Errors: []FieldError{{Field: "/changes/status", Code: ValidationInvalidValue}},
		}),
		"bulk_concurrency_failure": mock(409, func() APIError {
			response := MockVersionConflict(4)
			response.FailingRequestID = "case-7c58d2"
			return response
		}()),
		"expected_version_required": mock(428, APIError{
			Error: ErrorExpectedVersionRequired, Message: "If-Match is required for this update.", Retryable: false,
		}),
		"portal_attachment_staged": mock(201, map[string]interface{}{
			"attachment_token": "upload-00031", "filename": "pothole.jpg", "media_type": "image/jpeg", "size": 248031,
			"expires_at": "2026-08-25T12:15:00Z",
		}),
		"portal_service_request_created": mock(201, MockCreatedServiceRequest()),
		"public_branding": mock(200, map[string]interface{}{
			"organisation_name": "City 311", "logo_url": "/assets/city-logo.svg", "favicon_url": "/assets/favicon.ico",
			"portal_wallpaper_url": "/assets/portal-wallpaper.jpg", "login_header": "City services, one place", "public_header": "City 311",
			"public_footer": "City services", "primary_colour": "#005EA8", "accent_colour": "#FFB81C", "font_family": "Arial",
			"published": true, "version": 3, "updated_at": "2026-08-25T10:00:00Z",
		}),
		"public_content_home": mock(200, map[string]interface{}{
			"content_key": "HOME", "body": "<p>Welcome to City 311.</p>", "state": "PUBLISHED", "published": true,
			"version": 2, "updated_at": "2026-08-25T10:00:00Z",
		}),
		"public_help_submit": mock(200, map[string]interface{}{
			"help_key": "public.request.submit", "language": "EN",
			"body":    "<p>Describe the issue, choose its type, and provide the location where City service is needed.</p>",
			"version": 1, "updated_at": "2026-08-25T10:00:00Z",
		}),
		"identity_configuration_effective": mock(200, map[string]interface{}{
			"oidc_enabled": true, "saml_enabled": true,
			"oidc_issuer_url": "https://identity.example.test", "oidc_staff_client_id": "city311-staff",
			"oidc_public_client_id": "city311-public", "oidc_client_secret_configured": true,
			"saml_metadata_url": "https://identity.example.test/saml/metadata", "saml_sp_entity_id": "https://city311.example.test/saml",
			"actor_role_mappings": []map[string]interface{}{
				{"asserted_role": "service_agent", "application_role": "service_agent"},
				{"asserted_role": "supervisor", "application_role": "supervisor"},
				{"asserted_role": "department_manager", "application_role": "department_manager"},
				{"asserted_role": "platform_administrator", "application_role": "platform_administrator"},
				{"asserted_role": "workflow_designer", "application_role": "workflow_designer"},
			},
			"version": 1, "updated_at": "2026-08-25T10:00:00Z",
		}),
		"staff_queue": mock(200, map[string]interface{}{
			"items": []interface{}{map[string]interface{}{
				"request_id": "case-7c58d2", "request_number": "SR-2026-00041", "summary": "Pothole blocking the eastbound lane",
				"service_type": "POTHOLE", "status": "IN_PROGRESS", "owning_department": "STREETS", "council_district": "CENTRAL",
				"origin_class": "EXTERNAL", "source_channel": "PORTAL_ANONYMOUS", "primary_assignee_id": "staff-0007",
				"duplicate_group_id": nil, "version": 3, "updated_at": "2026-08-25T12:00:00Z", "available_actions": []string{"RESOLVE"},
			}},
			"next_page_token": nil, "total_count": 1, "applied_filters": map[string]interface{}{"status": []string{"IN_PROGRESS"}}, "sort": []string{"-updated_at"},
		}),
		"civicworks_event_acknowledged":     {HTTPStatus: 204, Body: nil},
		"civicworks_duplicate_acknowledged": {HTTPStatus: 204, Body: nil},
		"civicworks_invalid_signature": mock(401, APIError{
			Error: ErrorInvalidSignature, Message: "The CivicWorks event signature is invalid.", Retryable: false,
		}),
	}
}

func linkMocks(mocks map[string]MockContract) {
	endpointByMock := map[string]string{
		"service_request_created":           "service_request_create",
		"service_request_replay":            "service_request_create",
		"validation_error":                  "service_request_create",
		"idempotency_conflict":              "service_request_create",
		"unauthenticated":                   "service_request_create",
		"forbidden":                         "service_request_create",
		"not_found":                         "staff_request_detail",
		"version_conflict":                  "staff_request_transition",
		"rate_limited":                      "data_export",
		"invalid_reset_token":               "password_reset_confirm",
		"expired_reset_token":               "password_reset_confirm",
		"workflow_insufficient_scope":       "workflow_action_execute",
		"workflow_invalid_client":           "workflow_action_execute",
		"workflow_invalid_token":            "workflow_action_execute",
		"workflow_unavailable":              "workflow_action_execute",
		"geocode_success":                   "geocode_proxy",
		"geocode_not_found":                 "geocode_proxy",
		"geocode_unavailable":               "geocode_proxy",
		"civicworks_created":                "civicworks_work_order_create",
		"civicworks_completed_event":        "civicworks_event_callback",
		"anonymous_status_found":            "anonymous_status_lookup",
		"anonymous_status_not_found":        "anonymous_status_lookup",
		"password_reset_requested":          "password_reset_request",
		"session_anonymous":                 "session_current",
		"session_staff":                     "session_current",
		"local_sign_in_rejected":            "session_sign_in",
		"account_registration_accepted":     "account_register",
		"health_unavailable":                "health",
		"operation_pending":                 "operation_get",
		"operation_succeeded":               "operation_get",
		"empty_list":                        "staff_request_queue",
		"nested_validation_error":           "portal_service_request_submit",
		"bulk_validation_failure":           "staff_request_bulk",
		"bulk_concurrency_failure":          "staff_request_bulk",
		"expected_version_required":         "staff_request_transition",
		"portal_attachment_staged":          "portal_attachment_upload",
		"portal_service_request_created":    "portal_service_request_submit",
		"public_branding":                   "public_branding_get",
		"public_content_home":               "public_content_get",
		"public_help_submit":                "public_help_get",
		"identity_configuration_effective":  "identity_configuration_get",
		"staff_queue":                       "staff_request_queue",
		"civicworks_event_acknowledged":     "civicworks_event_callback",
		"civicworks_duplicate_acknowledged": "civicworks_event_callback",
		"civicworks_invalid_signature":      "civicworks_event_callback",
	}
	for name, item := range mocks {
		endpoint, present := endpointByMock[name]
		if !present {
			panic("mock has no endpoint linkage: " + name)
		}
		item.Endpoint = endpoint
		item.Role = "response"
		if name == "civicworks_completed_event" {
			item.Role = "request"
			item.HTTPStatus = 0
		}
		mocks[name] = item
	}
}
