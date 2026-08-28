package city311

import (
	"reflect"
	"testing"
)

func TestEveryEndpointHasDirectionAuthenticationAndResolvedSchemas(t *testing.T) {
	contract := NewContractDocument()
	for name, endpoint := range contract.Endpoints {
		if endpoint.Direction != EndpointProvidedByCRM && endpoint.Direction != EndpointConsumedByCRM {
			t.Errorf("endpoint %s has invalid direction %q", name, endpoint.Direction)
		}
		if endpoint.Authentication.Mode == "" {
			t.Errorf("endpoint %s has no authentication declaration", name)
		}
		assertEnumValue(t, contract, name+" authentication mode", "authentication_mode", endpoint.Authentication.Mode)
		if endpoint.Authentication.ActorClass != "" {
			assertEnumValue(t, contract, name+" actor class", "authenticated_actor_class", endpoint.Authentication.ActorClass)
		}
		for _, alternative := range endpoint.Authentication.Alternatives {
			for _, role := range alternative.ApplicationRoles {
				assertEnumValue(t, contract, name+" application role", "application_role", role)
			}
			if alternative.ResourceRelationship != "" {
				assertEnumValue(t, contract, name+" resource relationship", "resource_relationship", alternative.ResourceRelationship)
			}
		}
		if len(endpoint.SuccessStatuses) == 0 {
			t.Errorf("endpoint %s has no success status", name)
		}
		assertSchemaReference(t, contract, name+" request", endpoint.RequestSchema)
		assertSchemaReference(t, contract, name+" response", endpoint.ResponseSchema)
		for entity, schema := range endpoint.EntityResponseSchemas {
			assertSchemaReference(t, contract, name+" entity "+entity, schema)
		}
		for parameter, property := range endpoint.PathParameters {
			walkContractValue(t, contract, name+" path parameter "+parameter, property)
		}
	}
}

func TestEverySchemaAndEnumReferenceResolves(t *testing.T) {
	contract := NewContractDocument()
	for name, schema := range contract.Schemas {
		walkContractValue(t, contract, "schema "+name, schema)
	}
}

func TestCrossCuttingBrowserProtocolIsFrozen(t *testing.T) {
	contract := NewContractDocument()
	want := []string{
		"session_and_authorization", "optimistic_concurrency", "validation_errors", "lists",
		"idempotency", "asynchronous_operations", "atomic_bulk_operations", "attachment_transport",
	}
	for _, key := range want {
		if _, present := contract.Protocol[key]; !present {
			t.Errorf("missing protocol contract %s", key)
		}
	}

	concurrency := contract.Protocol["optimistic_concurrency"].(map[string]interface{})
	if concurrency["request_header"] != IfMatchHeader || concurrency["missing_status"] != 428 || concurrency["stale_status"] != 409 {
		t.Fatal("optimistic concurrency must freeze If-Match, missing-version, and stale-version behavior")
	}
	validation := contract.Protocol["validation_errors"].(map[string]interface{})
	if validation["code_enum_ref"] != "validation_code" {
		t.Fatal("field errors must use the controlled validation-code vocabulary")
	}

	errorProperties := contract.Schemas["error"]["properties"].(map[string]interface{})
	fieldError := errorProperties["errors"].(map[string]interface{})["items"].(map[string]interface{})
	fieldProperties := fieldError["properties"].(map[string]interface{})
	if fieldProperties["field"].(map[string]interface{})["format"] != "json-pointer" {
		t.Fatal("field error paths must use JSON Pointer")
	}
	if _, present := errorProperties["failing_request_id"]; !present {
		t.Fatal("bulk failures must identify the failing request")
	}
	serverErrors := contract.Protocol["server_error_policy"].(map[string]interface{})
	if !reflect.DeepEqual(serverErrors["declared_statuses"], []int{503}) || !reflect.DeepEqual(serverErrors["undeclared_statuses"], []int{500, 502, 504}) {
		t.Fatalf("server error policy has unexpected status boundary: %#v", serverErrors)
	}
}

func TestCivicWorksCallbackAndLocationRules(t *testing.T) {
	contract := NewContractDocument()

	for _, current := range []ServiceRequestStatus{
		ServiceRequestStatusResolved, ServiceRequestStatusClosed, ServiceRequestStatusReopened,
	} {
		plan, ok := PlanCivicWorksTransition(current, CivicWorksStatusCompleted)
		if !ok || len(plan) != 0 {
			t.Errorf("redelivered COMPLETED must be acknowledged as a no-op from %s", current)
		}
		key := string(current) + "+" + string(CivicWorksStatusCompleted)
		if published, present := contract.CivicWorksTransitionPlans[key]; !present || len(published) != 0 {
			t.Errorf("published transition plans must expose %s as a no-op", key)
		}
	}

	callback := contract.Endpoints["civicworks_event_callback"]
	if callback.Direction != EndpointProvidedByCRM || callback.Path != "/integrations/civicworks/events" {
		t.Fatal("CivicWorks callback must be a CRM-provided endpoint")
	}
	for _, header := range []string{"X-CivicWorks-Event-Id", "X-CivicWorks-Signature"} {
		if !contains(callback.RequiredHeaders, header) {
			t.Errorf("CivicWorks callback is missing %s", header)
		}
	}
	if contract.Endpoints["workflow_action_execute"].Direction != EndpointConsumedByCRM || contract.Endpoints["mapping_geocode"].Direction != EndpointConsumedByCRM {
		t.Fatal("fixture endpoints called by CRM must be marked consumed_by_crm")
	}
	if contract.Endpoints["mapping_geocode"].Path != "/internal/integrations/mapping/geocode" {
		t.Fatal("mapping geocode must use a server-only path separate from the browser proxy")
	}
	if callback.ResponseSchema != "empty_response" {
		t.Fatal("CivicWorks callback must publish an empty response schema for HTTP 204")
	}

	for serviceType, required := range map[string]bool{
		"TREE_MAINTENANCE": true, "POTHOLE": true, "MISSED_TRASH": true, "GENERAL_INQUIRY": false,
	} {
		rule, present := contract.ServiceTypeRules[serviceType]
		if !present || rule.LocationRequired != required || rule.ConfirmedCoordinatesRequired != required {
			t.Errorf("incorrect location rule for %s: %#v", serviceType, rule)
		}
	}
}

func TestProvidedAndConsumedEndpointPathsDoNotCollide(t *testing.T) {
	contract := NewContractDocument()
	seen := map[string]string{}
	for name, endpoint := range contract.Endpoints {
		key := endpoint.Method + " " + endpoint.Path
		if previous, present := seen[key]; present && previous != endpoint.Direction {
			t.Fatalf("endpoint %s collides with %s at %s", name, previous, key)
		}
		seen[key] = endpoint.Direction
	}
}

func TestApplicationAndIdentityRoleVocabulariesStayDistinct(t *testing.T) {
	contract := NewContractDocument()
	wantApplicationRoles := []string{
		"public_visitor", "constituent", "service_agent", "supervisor", "department_manager",
		"platform_administrator", "workflow_designer", "integration_client",
	}
	if !reflect.DeepEqual(contract.Enums["application_role"], wantApplicationRoles) {
		t.Fatalf("application roles do not match Table 4.1.1-A: %#v", contract.Enums["application_role"])
	}
	if !reflect.DeepEqual(contract.Enums["oidc_actor_type"], []string{"constituent"}) {
		t.Fatal("OIDC actor type must remain the narrow identity claim vocabulary")
	}
	if _, ambiguous := contract.Enums["actor_type"]; ambiguous {
		t.Fatal("ambiguous actor_type vocabulary must not be published")
	}
}

func TestRequiredClientSurfaceInventoryIsFrozen(t *testing.T) {
	contract := NewContractDocument()
	required := []string{
		"session_current", "session_sign_in", "account_register", "federated_sign_in_start", "portal_service_request_submit",
		"portal_draft_create", "portal_my_requests", "portal_link_anonymous_request", "profile_get", "password_change",
		"geocode_proxy", "portal_attachment_upload", "public_branding_get", "public_content_get", "public_help_get",
		"staff_request_queue", "staff_request_detail", "staff_request_transition", "staff_request_bulk",
		"staff_constituent_search", "admin_branding_update", "admin_content_update", "admin_custom_fields_create",
		"workflow_create", "workflow_execution_get", "identity_configuration_update", "integration_update",
		"report_run", "audit_list", "mail_send", "calendar_import", "health",
	}
	for _, name := range required {
		if _, present := contract.Endpoints[name]; !present {
			t.Errorf("required client surface %s is not frozen", name)
		}
	}
}

func TestMutatingVersionedEndpointsRequireIfMatch(t *testing.T) {
	contract := NewContractDocument()
	for _, name := range []string{
		"portal_draft_update", "profile_update", "staff_request_transition", "admin_branding_update", "admin_help_update",
		"admin_content_update", "admin_custom_fields_update", "workflow_update", "integration_update",
	} {
		endpoint := contract.Endpoints[name]
		if !contains(endpoint.RequiredHeaders, IfMatchHeader) {
			t.Errorf("versioned update %s does not require If-Match", name)
		}
		if endpoint.ErrorStatuses[string(ErrorExpectedVersionRequired)] != 428 || endpoint.ErrorStatuses[string(ErrorVersionConflict)] != 409 {
			t.Errorf("versioned update %s has incomplete concurrency errors", name)
		}
	}
}

func TestSubmissionAndRetryIdempotencyIsConsistent(t *testing.T) {
	contract := NewContractDocument()
	for _, name := range []string{"service_request_create", "portal_service_request_submit", "staff_request_bulk", "mail_send", "workflow_action_execute", "civicworks_work_order_create"} {
		endpoint := contract.Endpoints[name]
		if !contains(endpoint.RequiredHeaders, IdempotencyHeader) {
			t.Errorf("retryable operation %s does not require Idempotency-Key", name)
		}
	}
}

func TestBulkConcurrencyAndReminderActionsAreReachable(t *testing.T) {
	contract := NewContractDocument()

	bulkProperties := contract.Schemas["bulk_request"]["properties"].(map[string]interface{})
	items := bulkProperties["request_items"].(map[string]interface{})
	if items["items_ref"] != "bulk_request_item" || items["unique_by"] != "request_id" {
		t.Fatal("bulk requests must carry unique versioned request items")
	}
	itemProperties := contract.Schemas["bulk_request_item"]["properties"].(map[string]interface{})
	if _, present := itemProperties["expected_version"]; !present {
		t.Fatal("each bulk request item must carry its expected version")
	}
	bulkEndpoint := contract.Endpoints["staff_request_bulk"]
	if bulkEndpoint.ErrorStatuses[string(ErrorVersionConflict)] != 409 || !contains(bulkEndpoint.RequiredHeaders, IdempotencyHeader) {
		t.Fatal("bulk request must expose concurrency failure and retry safety")
	}
	bulkConflict := contract.Mocks["bulk_concurrency_failure"]
	if bulkConflict.HTTPStatus != 409 {
		t.Fatal("bulk concurrency failure must have a deterministic 409 mock")
	}

	reminder := contract.Endpoints["staff_reminder_action"]
	action := reminder.PathParameters["action"]
	if action["enum_ref"] != "reminder_action" || !reflect.DeepEqual(contract.Enums["reminder_action"], []string{"SNOOZE", "COMPLETE", "CANCEL"}) {
		t.Fatal("reminder action path must bind to the controlled action vocabulary")
	}
}

func TestLocalSessionCreationIsExplicit(t *testing.T) {
	contract := NewContractDocument()
	endpoint := contract.Endpoints["session_sign_in"]
	if endpoint.Method != "POST" || endpoint.Path != "/api/v1/session" || endpoint.RequestSchema != "local_sign_in" || endpoint.ResponseSchema != "session" {
		t.Fatal("local sign-in must explicitly create the shared session")
	}
	if endpoint.Authentication.Mode != "none" || endpoint.ErrorStatuses[string(ErrorUnauthenticated)] != 401 || len(endpoint.ErrorStatuses) != 2 {
		t.Fatal("local sign-in must be public and expose invalid credentials as UNAUTHENTICATED")
	}
	if endpoint.PrivacyRule == "" || endpoint.ErrorStatuses[string(ErrorNotFound)] != 0 {
		t.Fatal("local sign-in must not reveal whether the identifier exists")
	}
	identifier := contract.Schemas["local_sign_in"]["properties"].(map[string]interface{})["login_identifier"].(map[string]interface{})
	if !reflect.DeepEqual(identifier["accepted_forms"], []string{"verified_login_email", "local_login_identifier"}) {
		t.Fatal("local sign-in must declare whether it accepts email, identifier, or both")
	}
}

func TestPublicEndpointErrorsAndAccountPrivacyAreExact(t *testing.T) {
	contract := NewContractDocument()
	health := contract.Endpoints["health"]
	if !reflect.DeepEqual(health.ErrorStatuses, map[string]int{string(ErrorTemporarilyUnavailable): 503}) {
		t.Fatalf("health must publish only its required unavailable response: %#v", health.ErrorStatuses)
	}
	for _, name := range []string{"public_branding_get", "public_content_get", "public_help_get", "federated_sign_in_start"} {
		endpoint := contract.Endpoints[name]
		if _, present := endpoint.ErrorStatuses[string(ErrorUnauthenticated)]; present {
			t.Errorf("public endpoint %s declares an unreachable UNAUTHENTICATED response", name)
		}
		if _, present := endpoint.ErrorStatuses[string(ErrorForbidden)]; present {
			t.Errorf("public endpoint %s declares an unreachable FORBIDDEN response", name)
		}
	}

	registration := contract.Endpoints["account_register"]
	if registration.SuccessStatuses["accepted"] != 202 || registration.SuccessStatuses["existing_account_indistinguishable"] != 202 {
		t.Fatal("new and existing registration identifiers must have the same accepted status")
	}
	if registration.ResponseSchema != "account_registration_acknowledgement" || registration.PrivacyRule == "" {
		t.Fatal("registration must use the privacy-safe acknowledgement contract")
	}
}

func TestOptionalSessionErrorsAndCookieFallbackAreExact(t *testing.T) {
	contract := NewContractDocument()
	sessionProtocol := contract.Protocol["session_and_authorization"].(map[string]interface{})
	if sessionProtocol["optional_cookie_rule"] == "" {
		t.Fatal("optional-session endpoints must define invalid-cookie behavior")
	}
	for _, name := range []string{"session_current", "language_update", "portal_attachment_upload", "portal_service_request_submit", "geocode_proxy"} {
		endpoint := contract.Endpoints[name]
		if endpoint.Authentication.Mode != "session_cookie_optional" {
			t.Errorf("endpoint %s must remain anonymous-capable", name)
		}
		if _, present := endpoint.ErrorStatuses[string(ErrorUnauthenticated)]; present {
			t.Errorf("optional-session endpoint %s declares unreachable UNAUTHENTICATED", name)
		}
		if _, present := endpoint.ErrorStatuses[string(ErrorForbidden)]; present {
			t.Errorf("optional-session endpoint %s declares unreachable FORBIDDEN", name)
		}
	}
	geocode := contract.Endpoints["geocode_proxy"]
	want := map[string]int{
		string(ErrorAddressNotFound):           404,
		string(ErrorMapTemporarilyUnavailable): 503,
		string(ErrorValidation):                422,
	}
	if !reflect.DeepEqual(geocode.ErrorStatuses, want) {
		t.Fatalf("geocode proxy errors do not match the browser failure contract: %#v", geocode.ErrorStatuses)
	}
}

func TestSemanticVersioningStartsAtFirstPublication(t *testing.T) {
	versioning := NewContractDocument().Versioning
	if versioning.EffectiveAt != "merge_to_2024.9.x" {
		t.Fatal("the initial semantic-version scheme must take effect at merge")
	}
	if versioning.FirstPublishedVersion != ContractVersion || versioning.MajorVersionRule == "" {
		t.Fatal("the first publication and post-publication major-version rule must be explicit")
	}
}

func TestIdentityAdministrationPublishesEffectiveRuntimeConfiguration(t *testing.T) {
	contract := NewContractDocument()
	readEndpoint := contract.Endpoints["identity_configuration_get"]
	writeEndpoint := contract.Endpoints["identity_configuration_update"]
	if readEndpoint.ResponseSchema != "identity_configuration" || writeEndpoint.RequestSchema != "identity_configuration_write" || writeEndpoint.ResponseSchema != "identity_configuration" {
		t.Fatal("identity administration endpoints must separate effective configuration from permitted writes")
	}

	properties := contract.Schemas["identity_configuration"]["properties"].(map[string]interface{})
	for _, name := range []string{
		"oidc_enabled", "saml_enabled", "oidc_issuer_url", "oidc_staff_client_id", "oidc_public_client_id",
		"oidc_client_secret_configured", "saml_metadata_url", "saml_sp_entity_id", "actor_role_mappings",
	} {
		if _, present := properties[name]; !present {
			t.Errorf("identity configuration omits %s", name)
		}
	}
	if _, present := properties["oidc_client_secret"]; present {
		t.Fatal("identity configuration must never expose the OIDC client secret")
	}
	secretStatus := properties["oidc_client_secret_configured"].(map[string]interface{})
	if secretStatus["secret_value_returned"] != false || secretStatus["runtime_key"] != "OIDC_CLIENT_SECRET" {
		t.Fatal("identity configuration must expose only OIDC secret status")
	}

	writeProperties := contract.Schemas["identity_configuration_write"]["properties"].(map[string]interface{})
	if len(writeProperties) != 2 {
		t.Fatalf("identity administration may update enablement only: %#v", writeProperties)
	}
	for _, name := range []string{"oidc_enabled", "saml_enabled"} {
		if _, present := writeProperties[name]; !present {
			t.Errorf("identity write schema omits %s", name)
		}
	}

	mapping := contract.Schemas["actor_role_mapping"]["properties"].(map[string]interface{})
	if mapping["asserted_role"].(map[string]interface{})["enum_ref"] != "actor_role" || mapping["application_role"].(map[string]interface{})["enum_ref"] != "application_role" {
		t.Fatal("identity role mapping must bind inbound actor roles to application roles")
	}

	decisionPresent := false
	for _, decision := range contract.Decisions {
		if decision.ID == "IDENTITY-CONFIGURATION-BOUNDARY" {
			decisionPresent = true
		}
	}
	if !decisionPresent {
		t.Fatal("runtime-supplied identity configuration boundary must be recorded")
	}
}

func TestEveryTemplatedPathBindsAllParameters(t *testing.T) {
	contract := NewContractDocument()
	for name, endpoint := range contract.Endpoints {
		parameters := pathParameterNames(endpoint.Path)
		if len(parameters) != len(endpoint.PathParameters) {
			t.Errorf("endpoint %s binds %d of %d path parameters", name, len(endpoint.PathParameters), len(parameters))
			continue
		}
		for _, parameter := range parameters {
			if _, present := endpoint.PathParameters[parameter]; !present {
				t.Errorf("endpoint %s does not bind path parameter %s", name, parameter)
			}
		}
	}
}

func assertSchemaReference(t *testing.T, contract ContractDocument, context, schema string) {
	t.Helper()
	if schema == "" {
		return
	}
	if _, present := contract.Schemas[schema]; !present {
		t.Errorf("%s references missing schema %s", context, schema)
	}
}

func assertEnumValue(t *testing.T, contract ContractDocument, context, enum, value string) {
	t.Helper()
	if !contains(contract.Enums[enum], value) {
		t.Errorf("%s uses %q outside enum %s", context, value, enum)
	}
}

func walkContractValue(t *testing.T, contract ContractDocument, context string, value interface{}) {
	t.Helper()
	switch typed := value.(type) {
	case map[string]interface{}:
		for key, child := range typed {
			switch key {
			case "schema_ref", "items_ref":
				if name, ok := child.(string); ok {
					assertSchemaReference(t, contract, context, name)
				}
			case "enum_ref":
				if name, ok := child.(string); ok {
					if _, present := contract.Enums[name]; !present {
						t.Errorf("%s references missing enum %s", context, name)
					}
				}
			}
			walkContractValue(t, contract, context, child)
		}
	case []interface{}:
		for _, child := range typed {
			walkContractValue(t, contract, context, child)
		}
	case []map[string]interface{}:
		for _, child := range typed {
			walkContractValue(t, contract, context, child)
		}
	}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
