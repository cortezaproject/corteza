package city311

func clientSchemas() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"empty_request":  object([]string{}, map[string]interface{}{}),
		"empty_response": object([]string{}, map[string]interface{}{}),
		"health": object([]string{"status"}, map[string]interface{}{
			"status": map[string]interface{}{"const": "ok"},
		}),
		"current_actor": object([]string{"actor_id", "display_name", "application_roles", "department_codes", "district_codes", "capabilities", "scopes", "available_routes"}, map[string]interface{}{
			"actor_id":          map[string]interface{}{"type": "string"},
			"display_name":      map[string]interface{}{"type": "string"},
			"oidc_actor_type":   map[string]interface{}{"enum_ref": "oidc_actor_type", "nullable": true},
			"application_roles": arrayEnum("application_role"),
			"department_codes":  arrayEnum("department_code"),
			"district_codes":    arrayEnum("district_code"),
			"capabilities":      arrayEnum("capability"),
			"scopes":            arrayEnum("oauth_scope"),
			"available_routes":  arrayEnum("route"),
		}),
		"session": object([]string{"authenticated", "preferred_language"}, map[string]interface{}{
			"authenticated":      map[string]interface{}{"type": "boolean"},
			"actor":              map[string]interface{}{"schema_ref": "current_actor", "nullable": true},
			"preferred_language": map[string]interface{}{"enum_ref": "language"},
			"expires_at":         map[string]interface{}{"schema_ref": "timestamp", "nullable": true},
		}),
		"local_sign_in": object([]string{"login_identifier", "password"}, map[string]interface{}{
			"login_identifier": map[string]interface{}{
				"one_of": []map[string]interface{}{
					{"type": "string", "format": "email", "max_length": 254},
					{"type": "string", "min_length": 3, "max_length": 64, "pattern": "^[a-z0-9._-]+$"},
				},
				"case_insensitive": true,
				"accepted_forms":   []string{"verified_login_email", "local_login_identifier"},
			},
			"password": map[string]interface{}{"type": "string", "min_length": 1, "write_only": true},
		}),
		"account_registration_acknowledgement": object([]string{"accepted"}, map[string]interface{}{
			"accepted": map[string]interface{}{"const": true},
		}),
		"account_registration": object([]string{"display_name", "email", "login_identifier", "password", "preferred_language"}, map[string]interface{}{
			"display_name":       stringProperty(1, 120),
			"email":              map[string]interface{}{"type": "string", "format": "email"},
			"login_identifier":   map[string]interface{}{"type": "string", "min_length": 3, "max_length": 64, "pattern": "^[a-z0-9._-]+$"},
			"password":           passwordProperty(),
			"preferred_language": map[string]interface{}{"enum_ref": "language"},
		}),
		"federated_redirect": object([]string{"authorization_url"}, map[string]interface{}{
			"authorization_url":          map[string]interface{}{"type": "string", "format": "uri"},
			"link_confirmation_required": map[string]interface{}{"type": "boolean", "default": false},
		}),
		"operation": object([]string{"operation_id", "kind", "status", "created_at", "updated_at"}, map[string]interface{}{
			"operation_id": map[string]interface{}{"type": "string"},
			"kind":         map[string]interface{}{"type": "string"},
			"status":       map[string]interface{}{"enum_ref": "operation_status"},
			"progress":     map[string]interface{}{"type": "integer", "minimum": 0, "maximum": 100},
			"result":       map[string]interface{}{"type": "object", "nullable": true},
			"error":        map[string]interface{}{"schema_ref": "error", "nullable": true},
			"created_at":   timestampProperty(),
			"updated_at":   timestampProperty(),
			"completed_at": map[string]interface{}{"schema_ref": "timestamp", "nullable": true},
		}),
		"list_response": object([]string{"items", "next_page_token", "total_count", "applied_filters", "sort"}, map[string]interface{}{
			"items":           map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "object"}},
			"next_page_token": map[string]interface{}{"type": "string", "nullable": true, "opaque": true},
			"total_count":     map[string]interface{}{"type": "integer", "minimum": 0},
			"applied_filters": map[string]interface{}{"type": "object"},
			"sort":            stringArray(),
		}),
		"geocode_request": object([]string{"address"}, map[string]interface{}{
			"address": map[string]interface{}{"type": "string", "min_length": 1},
		}),
		"geocode_response": object([]string{"address", "latitude", "longitude", "precision_digits", "provider"}, map[string]interface{}{
			"address":          map[string]interface{}{"type": "string"},
			"latitude":         coordinateProperty(-90, 90),
			"longitude":        coordinateProperty(-180, 180),
			"precision_digits": map[string]interface{}{"const": 4},
			"provider":         map[string]interface{}{"const": "BENCHMARK_MAP"},
		}),
		"portal_attachment_upload": object([]string{"file", "filename", "media_type"}, map[string]interface{}{
			"file":       map[string]interface{}{"type": "string", "format": "binary", "transport": "multipart/form-data", "maximum_bytes": 10485760},
			"filename":   stringProperty(1, 120),
			"media_type": attachmentMediaTypeProperty(),
		}),
		"portal_attachment": object([]string{"attachment_token", "filename", "media_type", "size", "expires_at"}, map[string]interface{}{
			"attachment_token": map[string]interface{}{"type": "string", "single_use": true},
			"filename":         stringProperty(1, 120),
			"media_type":       attachmentMediaTypeProperty(),
			"size":             map[string]interface{}{"type": "integer", "minimum": 1, "maximum": 10485760},
			"expires_at":       timestampProperty(),
		}),
		"binary_attachment": object([]string{"content_type", "content_disposition", "body"}, map[string]interface{}{
			"content_type":        map[string]interface{}{"type": "string"},
			"content_disposition": map[string]interface{}{"type": "string"},
			"body":                map[string]interface{}{"type": "string", "format": "binary"},
		}),
		"portal_service_request_submit": portalRequestWriteSchema(false),
		"portal_draft_write":            portalRequestWriteSchema(true),
		"portal_request_summary": object([]string{"request_id", "request_number", "summary", "service_type", "status", "owning_department", "updated_at"}, map[string]interface{}{
			"request_id":        map[string]interface{}{"type": "string"},
			"request_number":    requestNumberProperty(),
			"summary":           stringProperty(5, 160),
			"service_type":      map[string]interface{}{"enum_ref": "service_type"},
			"status":            map[string]interface{}{"enum_ref": "service_request_status"},
			"owning_department": map[string]interface{}{"enum_ref": "department_code"},
			"updated_at":        timestampProperty(),
		}),
		"anonymous_request_link": object([]string{"request_number", "email"}, map[string]interface{}{
			"request_number": requestNumberProperty(),
			"email":          map[string]interface{}{"type": "string", "format": "email"},
		}),
		"profile_update": object([]string{}, map[string]interface{}{
			"display_name":       stringProperty(1, 120),
			"phone_numbers":      map[string]interface{}{"type": "array", "max_items": 3, "items_ref": "phone_number"},
			"addresses":          map[string]interface{}{"type": "array", "max_items": 5, "items_ref": "structured_address", "maximum_primary_items": 1},
			"primary_category":   map[string]interface{}{"enum_ref": "contact_category"},
			"preferred_language": map[string]interface{}{"enum_ref": "language"},
		}),
		"password_change": object([]string{"current_password", "new_password"}, map[string]interface{}{
			"current_password": map[string]interface{}{"type": "string"},
			"new_password":     passwordProperty(),
		}),
		"login_identifier_change": object([]string{"current_password", "login_identifier"}, map[string]interface{}{
			"current_password": map[string]interface{}{"type": "string"},
			"login_identifier": map[string]interface{}{"type": "string", "min_length": 3, "max_length": 64, "pattern": "^[a-z0-9._-]+$"},
		}),
		"language_preference": object([]string{"language"}, map[string]interface{}{
			"language": map[string]interface{}{"enum_ref": "language"},
		}),
		"branding": versionedObject([]string{"organisation_name", "primary_colour", "accent_colour", "font_family", "published"}, map[string]interface{}{
			"organisation_name":    map[string]interface{}{"type": "string"},
			"logo_url":             map[string]interface{}{"type": "string", "format": "uri", "nullable": true},
			"favicon_url":          map[string]interface{}{"type": "string", "format": "uri", "nullable": true},
			"portal_wallpaper_url": map[string]interface{}{"type": "string", "format": "uri", "nullable": true},
			"login_header":         map[string]interface{}{"type": "string"},
			"public_header":        map[string]interface{}{"type": "string"},
			"public_footer":        map[string]interface{}{"type": "string"},
			"primary_colour":       map[string]interface{}{"type": "string", "format": "css-color"},
			"accent_colour":        map[string]interface{}{"type": "string", "format": "css-color"},
			"font_family":          map[string]interface{}{"type": "string"},
			"published":            map[string]interface{}{"type": "boolean"},
		}),
		"branding_write": object([]string{}, map[string]interface{}{
			"organisation_name": map[string]interface{}{"type": "string"}, "logo_url": map[string]interface{}{"type": "string", "format": "uri"},
			"favicon_url": map[string]interface{}{"type": "string", "format": "uri"}, "portal_wallpaper_url": map[string]interface{}{"type": "string", "format": "uri"},
			"login_header": map[string]interface{}{"type": "string"}, "public_header": map[string]interface{}{"type": "string"}, "public_footer": map[string]interface{}{"type": "string"},
			"primary_colour": map[string]interface{}{"type": "string", "format": "css-color"}, "accent_colour": map[string]interface{}{"type": "string", "format": "css-color"},
			"font_family": map[string]interface{}{"type": "string"},
		}),
		"content_object": versionedObject([]string{"content_key", "body", "state", "published"}, map[string]interface{}{
			"content_key": map[string]interface{}{"enum_ref": "public_content_key"},
			"body":        map[string]interface{}{"type": "string", "format": "sanitized_html"},
			"state":       map[string]interface{}{"enum": []string{"DRAFT", "PUBLISHED"}},
			"published":   map[string]interface{}{"type": "boolean"},
		}),
		"content_write": object([]string{"body"}, map[string]interface{}{"body": map[string]interface{}{"type": "string", "format": "sanitized_html"}}),
		"help_content": versionedObject([]string{"help_key", "language", "body"}, map[string]interface{}{
			"help_key": map[string]interface{}{"type": "string", "enum_source": "help_keys"},
			"language": map[string]interface{}{"enum_ref": "language"},
			"body":     map[string]interface{}{"type": "string", "format": "sanitized_html"},
		}),
		"help_write": object([]string{"language", "body"}, map[string]interface{}{
			"language": map[string]interface{}{"enum_ref": "language"}, "body": map[string]interface{}{"type": "string", "format": "sanitized_html"},
		}),
		"request_queue_item": object([]string{"request_id", "request_number", "summary", "service_type", "status", "owning_department", "origin_class", "source_channel", "version", "updated_at", "available_actions"}, map[string]interface{}{
			"request_id": map[string]interface{}{"type": "string"}, "request_number": requestNumberProperty(), "summary": stringProperty(5, 160),
			"service_type": map[string]interface{}{"enum_ref": "service_type"}, "status": map[string]interface{}{"enum_ref": "service_request_status"},
			"owning_department": map[string]interface{}{"enum_ref": "department_code"}, "council_district": map[string]interface{}{"enum_ref": "district_code"},
			"origin_class": map[string]interface{}{"enum_ref": "origin_class"}, "source_channel": map[string]interface{}{"enum_ref": "source_channel"},
			"primary_assignee_id": map[string]interface{}{"type": "string", "nullable": true}, "duplicate_group_id": map[string]interface{}{"type": "string", "nullable": true},
			"version": map[string]interface{}{"type": "integer", "minimum": 1}, "updated_at": timestampProperty(), "available_actions": arrayEnum("request_action"),
		}),
		"staff_service_request_detail": object([]string{"request", "available_actions", "primary_assignee_id", "collaborator_ids", "reminders", "history", "audit", "external_work_order"}, map[string]interface{}{
			"request":             map[string]interface{}{"schema_ref": "service_request_record"},
			"available_actions":   arrayEnum("request_action"),
			"primary_assignee_id": map[string]interface{}{"type": "string", "nullable": true},
			"collaborator_ids":    stringArray(),
			"reminders":           map[string]interface{}{"type": "array", "items_ref": "reminder"},
			"history":             map[string]interface{}{"type": "array", "items_ref": "public_history_item"},
			"audit":               map[string]interface{}{"type": "array", "items_ref": "audit_event"},
			"external_work_order": map[string]interface{}{"schema_ref": "civicworks_work_order", "nullable": true},
		}),
		"staff_service_request_create": object([]string{"request", "constituent"}, map[string]interface{}{
			"request":     map[string]interface{}{"schema_ref": "portal_service_request_submit"},
			"constituent": map[string]interface{}{"one_of": []map[string]interface{}{{"schema_ref": "constituent_reference"}, {"schema_ref": "constituent_create"}}},
		}),
		"request_transition": object([]string{"to_status"}, map[string]interface{}{
			"to_status": map[string]interface{}{"enum_ref": "service_request_status"}, "reason": map[string]interface{}{"type": "string"},
		}),
		"reassignment": object([]string{"assignee_id", "reason"}, map[string]interface{}{
			"assignee_id": map[string]interface{}{"type": "string"}, "reason": map[string]interface{}{"type": "string", "min_length": 1},
		}),
		"reason": object([]string{"reason"}, map[string]interface{}{"reason": map[string]interface{}{"type": "string", "min_length": 1}}),
		"constituent_link": object([]string{"constituent_id", "relationship_type", "portal_visible", "notify_status"}, map[string]interface{}{
			"constituent_id": map[string]interface{}{"type": "string"}, "relationship_type": map[string]interface{}{"enum_ref": "relationship_type"},
			"portal_visible": map[string]interface{}{"type": "boolean"}, "notify_status": map[string]interface{}{"type": "boolean"},
		}),
		"constituent_unlink":    object([]string{"reason"}, map[string]interface{}{"reason": map[string]interface{}{"type": "string"}}),
		"constituent_reference": object([]string{"constituent_id"}, map[string]interface{}{"constituent_id": map[string]interface{}{"type": "string"}}),
		"constituent_create":    object([]string{"display_name", "email"}, map[string]interface{}{"display_name": stringProperty(1, 120), "email": map[string]interface{}{"type": "string", "format": "email"}}),
		"request_note": object([]string{"body", "portal_visible"}, map[string]interface{}{
			"note_id": map[string]interface{}{"type": "string"}, "request_id": map[string]interface{}{"type": "string"}, "author_constituent_id": map[string]interface{}{"type": "string"},
			"body": stringProperty(1, 2000), "created_at": timestampProperty(), "portal_visible": map[string]interface{}{"type": "boolean"},
		}),
		"reminder_write": object([]string{"title", "due_at", "timezone", "recipient_staff_id", "channel"}, map[string]interface{}{
			"title": map[string]interface{}{"type": "string"}, "due_at": timestampProperty(), "timezone": map[string]interface{}{"type": "string"},
			"recipient_staff_id": map[string]interface{}{"type": "string"}, "channel": map[string]interface{}{"enum_ref": "reminder_channel"},
		}),
		"reminder_action": object([]string{}, map[string]interface{}{
			"due_at": map[string]interface{}{"schema_ref": "timestamp", "conditional_required_for_path_action": "SNOOZE", "constraint": "later than the existing due_at"},
		}),
		"bulk_request_item": object([]string{"request_id", "expected_version"}, map[string]interface{}{
			"request_id":       map[string]interface{}{"type": "string"},
			"expected_version": map[string]interface{}{"type": "integer", "minimum": 1},
		}),
		"bulk_request": object([]string{"request_items", "action", "changes"}, map[string]interface{}{
			"request_items": map[string]interface{}{"type": "array", "min_items": 1, "unique_by": "request_id", "items_ref": "bulk_request_item"},
			"action":        map[string]interface{}{"enum": []string{"UPDATE", "CLOSE"}},
			"changes":       map[string]interface{}{"type": "object", "allowed_fields": []string{"primary_assignee_id", "priority", "status", "staff_note"}},
		}),
		"bulk_result": object([]string{"updated_request_ids", "updated_count"}, map[string]interface{}{
			"updated_request_ids": stringArray(), "updated_count": map[string]interface{}{"type": "integer", "minimum": 0},
		}),
		"duplicate_group_change": object([]string{"duplicate_group_id", "reason"}, map[string]interface{}{
			"duplicate_group_id": map[string]interface{}{"type": "string"}, "reason": map[string]interface{}{"type": "string"},
		}),
		"reopen_request": object([]string{"request_id", "status"}, map[string]interface{}{
			"request_id": map[string]interface{}{"type": "string"}, "status": map[string]interface{}{"const": "PENDING_APPROVAL"},
		}),
		"origin_override": object([]string{"origin_class", "reason"}, map[string]interface{}{
			"origin_class": map[string]interface{}{"enum_ref": "origin_class"}, "reason": map[string]interface{}{"type": "string", "min_length": 1},
		}),
		"scope_override": object([]string{"department_code", "district_codes", "reason"}, map[string]interface{}{
			"department_code": map[string]interface{}{"enum_ref": "department_code"}, "district_codes": arrayEnum("district_code"), "reason": map[string]interface{}{"type": "string", "min_length": 1},
		}),
		"rollback": object([]string{"target_version"}, map[string]interface{}{"target_version": map[string]interface{}{"type": "integer", "minimum": 1}}),
		"category": versionedObject([]string{"code", "active"}, map[string]interface{}{
			"code": map[string]interface{}{"type": "string"}, "active": map[string]interface{}{"type": "boolean"}, "labels": localizedStringsProperty(),
		}),
		"category_write": object([]string{"code", "active", "labels"}, map[string]interface{}{
			"code": map[string]interface{}{"type": "string"}, "active": map[string]interface{}{"type": "boolean"}, "labels": localizedStringsProperty(),
		}),
		"custom_field_definition": versionedObject([]string{"key", "labels", "entity", "field_type", "required", "active"}, map[string]interface{}{
			"key": map[string]interface{}{"type": "string", "stable": true}, "labels": localizedStringsProperty(),
			"entity": map[string]interface{}{"enum": []string{"constituent", "service_request"}}, "field_type": map[string]interface{}{"enum_ref": "custom_field_type"},
			"required": map[string]interface{}{"type": "boolean"}, "default": map[string]interface{}{}, "active": map[string]interface{}{"type": "boolean"},
			"validation": map[string]interface{}{"type": "object"}, "choice_values": map[string]interface{}{"type": "array", "ordered": true},
		}),
		"workflow_definition": versionedObject([]string{"workflow_id", "name", "trigger", "active", "conditions", "actions"}, map[string]interface{}{
			"workflow_id": map[string]interface{}{"type": "string"}, "name": map[string]interface{}{"type": "string"},
			"trigger": map[string]interface{}{"enum": []string{"SERVICE_REQUEST_CREATED", "SERVICE_REQUEST_STATUS_CHANGED"}}, "active": map[string]interface{}{"type": "boolean"},
			"conditions": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "object"}}, "actions": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "object"}},
		}),
		"workflow_test": object([]string{"request_id"}, map[string]interface{}{"request_id": map[string]interface{}{"type": "string"}}),
		"actor_role_mapping": object([]string{"asserted_role", "application_role"}, map[string]interface{}{
			"asserted_role":    map[string]interface{}{"enum_ref": "actor_role"},
			"application_role": map[string]interface{}{"enum_ref": "application_role"},
		}),
		"identity_configuration": versionedObject([]string{
			"oidc_enabled", "saml_enabled", "oidc_issuer_url", "oidc_staff_client_id", "oidc_public_client_id",
			"oidc_client_secret_configured", "saml_metadata_url", "saml_sp_entity_id", "actor_role_mappings",
		}, map[string]interface{}{
			"oidc_enabled":                  map[string]interface{}{"type": "boolean"},
			"saml_enabled":                  map[string]interface{}{"type": "boolean"},
			"oidc_issuer_url":               map[string]interface{}{"type": "string", "format": "uri", "read_only": true, "runtime_key": "OIDC_ISSUER_URL"},
			"oidc_staff_client_id":          map[string]interface{}{"type": "string", "read_only": true, "runtime_key": "OIDC_STAFF_CLIENT_ID"},
			"oidc_public_client_id":         map[string]interface{}{"type": "string", "read_only": true, "runtime_key": "OIDC_PUBLIC_CLIENT_ID"},
			"oidc_client_secret_configured": map[string]interface{}{"type": "boolean", "read_only": true, "runtime_key": "OIDC_CLIENT_SECRET", "secret_value_returned": false},
			"saml_metadata_url":             map[string]interface{}{"type": "string", "format": "uri", "read_only": true, "runtime_key": "SAML_METADATA_URL"},
			"saml_sp_entity_id":             map[string]interface{}{"type": "string", "format": "uri", "read_only": true, "runtime_key": "SAML_SP_ENTITY_ID"},
			"actor_role_mappings":           map[string]interface{}{"type": "array", "items_ref": "actor_role_mapping", "read_only": true},
		}),
		"identity_configuration_write": object([]string{}, map[string]interface{}{
			"oidc_enabled": map[string]interface{}{"type": "boolean"},
			"saml_enabled": map[string]interface{}{"type": "boolean"},
		}),
		"integration_connection": versionedObject([]string{"integration_id", "kind", "active", "secret_configured"}, map[string]interface{}{
			"integration_id": map[string]interface{}{"type": "string"}, "kind": map[string]interface{}{"enum": []string{"CIVICWORKS", "MAPPING", "WORKFLOW_OAUTH", "MAIL", "IDENTITY"}},
			"active": map[string]interface{}{"type": "boolean"}, "secret_configured": map[string]interface{}{"type": "boolean"}, "secret": map[string]interface{}{"type": "string", "write_only": true, "never_returned": true},
		}),
		"integration_connection_write": object([]string{"active"}, map[string]interface{}{
			"active": map[string]interface{}{"type": "boolean"}, "configuration": map[string]interface{}{"type": "object"}, "secret": map[string]interface{}{"type": "string", "write_only": true},
		}),
		"secret_rotation": object([]string{"new_secret"}, map[string]interface{}{"new_secret": map[string]interface{}{"type": "string", "write_only": true}}),
		"report_catalogue_item": object([]string{"report_key", "name", "supported_filters", "supported_grouping", "supported_sort"}, map[string]interface{}{
			"report_key": map[string]interface{}{"type": "string"}, "name": map[string]interface{}{"type": "string"},
			"supported_filters": stringArray(), "supported_grouping": stringArray(), "supported_sort": stringArray(),
		}),
		"report_definition": versionedObject([]string{"report_id", "name", "entity", "columns", "filters", "sort"}, map[string]interface{}{
			"report_id": map[string]interface{}{"type": "string"}, "name": map[string]interface{}{"type": "string"},
			"entity":  map[string]interface{}{"enum": []string{"service_requests", "constituents", "follow_up_actions"}},
			"columns": map[string]interface{}{"type": "array", "min_items": 1, "max_items": 20, "items": map[string]interface{}{"type": "string"}},
			"filters": map[string]interface{}{"type": "object"}, "grouping": map[string]interface{}{"type": "string", "nullable": true},
			"sort": map[string]interface{}{"type": "array", "max_items": 3, "items": map[string]interface{}{"type": "string"}},
		}),
		"report_run":           object([]string{"definition"}, map[string]interface{}{"definition": map[string]interface{}{"schema_ref": "report_definition"}}),
		"report_share":         object([]string{"roles"}, map[string]interface{}{"roles": arrayEnum("application_role")}),
		"report_export":        object([]string{"format"}, map[string]interface{}{"format": map[string]interface{}{"const": "CSV"}}),
		"contact_email_export": object([]string{"filters"}, map[string]interface{}{"filters": map[string]interface{}{"type": "object"}}),
		"audit_export":         object([]string{"filters"}, map[string]interface{}{"filters": map[string]interface{}{"type": "object"}}),
		"mail_compose": object([]string{"to", "subject", "text"}, map[string]interface{}{
			"template_id": map[string]interface{}{"type": "string", "nullable": true}, "to": map[string]interface{}{"type": "array", "min_items": 1, "items": map[string]interface{}{"type": "string", "format": "email"}},
			"subject": map[string]interface{}{"type": "string", "min_length": 1}, "text": map[string]interface{}{"type": "string", "min_length": 1},
			"html": map[string]interface{}{"type": "string", "format": "sanitized_html"}, "attachments": map[string]interface{}{"type": "array", "max_items": 3, "items_ref": "portal_attachment"},
		}),
		"mail_preview": object([]string{"subject", "text", "html"}, map[string]interface{}{
			"subject": map[string]interface{}{"type": "string"}, "text": map[string]interface{}{"type": "string"}, "html": map[string]interface{}{"type": "string"},
		}),
		"mail_delivery": object([]string{"delivery_id", "status", "attempts", "updated_at"}, map[string]interface{}{
			"delivery_id": map[string]interface{}{"type": "string"}, "status": map[string]interface{}{"enum": []string{"PENDING", "DELIVERED", "TERMINAL_FAILURE"}},
			"attempts": map[string]interface{}{"type": "integer", "minimum": 0, "maximum": 3}, "updated_at": timestampProperty(), "error": map[string]interface{}{"schema_ref": "error", "nullable": true},
		}),
		"calendar_import": object([]string{"ics"}, map[string]interface{}{"ics": map[string]interface{}{"type": "string", "format": "text/calendar"}}),
		"calendar_export": object([]string{"content_type", "body"}, map[string]interface{}{
			"content_type": map[string]interface{}{"const": "text/calendar"}, "body": map[string]interface{}{"type": "string"},
		}),
		"civicworks_work_order_create": object([]string{"source_case_id", "service_request_number", "service_type", "summary", "department_code", "callback_url"}, map[string]interface{}{
			"source_case_id": map[string]interface{}{"type": "string"}, "service_request_number": requestNumberProperty(), "service_type": map[string]interface{}{"enum_ref": "service_type"},
			"summary": stringProperty(5, 160), "department_code": map[string]interface{}{"enum_ref": "department_code"}, "location": map[string]interface{}{"schema_ref": "location_input"},
			"callback_url": map[string]interface{}{"const": "/integrations/civicworks/events"},
		}),
		"civicworks_work_order": object([]string{"work_order_id", "source_case_id", "service_request_number", "status", "external_status_url", "version", "created_at", "updated_at"}, map[string]interface{}{
			"work_order_id": map[string]interface{}{"type": "string"}, "source_case_id": map[string]interface{}{"type": "string"}, "service_request_number": requestNumberProperty(),
			"status": map[string]interface{}{"enum_ref": "civicworks_status"}, "external_status_url": map[string]interface{}{"type": "string", "format": "uri"},
			"version": map[string]interface{}{"type": "integer", "minimum": 1}, "created_at": timestampProperty(), "updated_at": timestampProperty(),
		}),
		"civicworks_event": object([]string{"event_id", "event_type", "work_order_id", "source_case_id", "previous_status", "status", "version", "occurred_at"}, map[string]interface{}{
			"event_id": map[string]interface{}{"type": "string"}, "event_type": map[string]interface{}{"const": "work_order.status_changed"},
			"work_order_id": map[string]interface{}{"type": "string"}, "source_case_id": map[string]interface{}{"type": "string"},
			"previous_status": map[string]interface{}{"enum_ref": "civicworks_status"}, "status": map[string]interface{}{"enum_ref": "civicworks_status"},
			"version": map[string]interface{}{"type": "integer", "minimum": 1}, "occurred_at": timestampProperty(),
		}),
		"event_acknowledgement": object([]string{}, map[string]interface{}{}),
		"timestamp":             map[string]interface{}{"type": "string", "format": "date-time", "explicit_utc_offset": true},
		"error": object([]string{"error", "message", "retryable"}, map[string]interface{}{
			"error": map[string]interface{}{"enum_ref": "error_code"}, "message": map[string]interface{}{"type": "string"}, "retryable": map[string]interface{}{"type": "boolean"},
			"errors": map[string]interface{}{"type": "array", "items": object([]string{"field", "code"}, map[string]interface{}{
				"field": map[string]interface{}{"type": "string", "format": "json-pointer", "examples": []string{"/requester/email", "/attachments/2/media_type"}},
				"code":  map[string]interface{}{"enum_ref": "validation_code"},
			})},
			"current_version":    map[string]interface{}{"type": "integer", "minimum": 1},
			"failing_request_id": map[string]interface{}{"type": "string"},
			"operation_id":       map[string]interface{}{"type": "string"},
		}),
	}
}

func portalRequestWriteSchema(draft bool) map[string]interface{} {
	required := []string{"summary", "description", "service_type", "requester"}
	if draft {
		required = []string{}
	}
	return map[string]interface{}{
		"type":                       "object",
		"required":                   required,
		"conditional_validation_ref": "service_type_rules",
		"properties": map[string]interface{}{
			"summary": stringProperty(5, 160), "description": map[string]interface{}{"type": "string", "min_length": 10, "max_length": 5000, "content": "plain_text"},
			"service_type": map[string]interface{}{"enum_ref": "service_type"}, "requester": map[string]interface{}{"schema_ref": "requester"},
			"location": map[string]interface{}{"schema_ref": "location_input"}, "attachment_tokens": map[string]interface{}{"type": "array", "max_items": 5, "items": map[string]interface{}{"type": "string"}},
			"custom_fields": map[string]interface{}{"type": "object", "additional_properties": true},
		},
	}
}

func versionedObject(required []string, properties map[string]interface{}) map[string]interface{} {
	properties["version"] = map[string]interface{}{"type": "integer", "minimum": 1}
	properties["updated_at"] = timestampProperty()
	required = appendUnique(required, "version", "updated_at")
	return object(required, properties)
}

func arrayEnum(enum string) map[string]interface{} {
	return map[string]interface{}{"type": "array", "unique_items": true, "items": map[string]interface{}{"enum_ref": enum}}
}

func stringArray() map[string]interface{} {
	return map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}}
}

func passwordProperty() map[string]interface{} {
	return map[string]interface{}{"type": "string", "min_length": 12, "max_length": 128, "minimum_character_classes": 3, "character_classes": []string{"uppercase", "lowercase", "digit", "symbol"}}
}

func attachmentMediaTypeProperty() map[string]interface{} {
	return map[string]interface{}{"enum": []string{"image/jpeg", "image/png", "application/pdf", "text/plain", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"}}
}

func localizedStringsProperty() map[string]interface{} {
	return map[string]interface{}{"type": "object", "required": []string{"EN"}, "properties": map[string]interface{}{"EN": map[string]interface{}{"type": "string"}, "ES": map[string]interface{}{"type": "string"}, "VI": map[string]interface{}{"type": "string"}}}
}
