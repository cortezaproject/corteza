package city311

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// NewOpenAPIDocument translates the authoritative City 311 contract into a
// standard OpenAPI 3.1 document consumable by client generators, validators,
// mock servers, Swagger UI, and Redoc.
func NewOpenAPIDocument() map[string]interface{} {
	contract := NewContractDocument()
	paths := map[string]interface{}{}
	consumed := map[string]interface{}{}
	for _, name := range sortedEndpointNames(contract.Endpoints) {
		endpoint := contract.Endpoints[name]
		operation := openAPIOperation(contract, name, endpoint)
		if endpoint.Direction == EndpointConsumedByCRM {
			consumed[name] = map[string]interface{}{
				"method": endpoint.Method, "path": endpoint.Path, "operation": operation,
			}
			continue
		}
		pathItem, _ := paths[endpoint.Path].(map[string]interface{})
		if pathItem == nil {
			pathItem = map[string]interface{}{}
			paths[endpoint.Path] = pathItem
		}
		method := strings.ToLower(endpoint.Method)
		if _, duplicate := pathItem[method]; duplicate {
			panic("duplicate OpenAPI operation: " + endpoint.Method + " " + endpoint.Path)
		}
		pathItem[method] = operation
	}

	components := map[string]interface{}{}
	for name, schema := range contract.Schemas {
		components[name] = openAPISchema(contract, schema)
	}
	for name, values := range contract.Enums {
		components["enum_"+name] = map[string]interface{}{
			"type": "string", "enum": values, "x-city311-enum-name": name,
		}
	}

	return map[string]interface{}{
		"openapi":           "3.1.0",
		"jsonSchemaDialect": "https://json-schema.org/draft/2020-12/schema",
		"info": map[string]interface{}{
			"title": "City 311 CRM client/server API", "version": contract.ContractVersion,
			"description": "Developer 1-owned BE-contract-v1 generated from NewContractDocument().",
		},
		"servers": []map[string]interface{}{{"url": "/", "description": "City 311 CRM origin"}},
		"paths":   paths,
		"components": map[string]interface{}{
			"schemas": components,
			"securitySchemes": map[string]interface{}{
				"sessionCookie": map[string]interface{}{"type": "apiKey", "in": "cookie", "name": "corteza_session"},
				"oauth2": map[string]interface{}{
					"type": "oauth2", "flows": map[string]interface{}{
						"clientCredentials": map[string]interface{}{
							"tokenUrl": "/oauth2/token", "scopes": map[string]interface{}{
								ScopeRequestWrite: "Create service requests", ScopeCRMExport: "Export CRM data", ScopeWorkflowExecute: "Execute workflow actions",
							},
						},
					},
				},
				"serverApiToken":      map[string]interface{}{"type": "http", "scheme": "bearer"},
				"civicworksSignature": map[string]interface{}{"type": "apiKey", "in": "header", "name": "X-CivicWorks-Signature"},
			},
		},
		"x-city311-contract-versioning": contract.Versioning,
		"x-city311-consumed-operations": consumed,
	}
}

func openAPIOperation(contract ContractDocument, name string, endpoint EndpointContract) map[string]interface{} {
	operation := map[string]interface{}{
		"operationId":              name,
		"summary":                  strings.ReplaceAll(name, "_", " "),
		"tags":                     []string{operationTag(endpoint.Path)},
		"security":                 openAPISecurity(endpoint),
		"responses":                openAPIResponses(contract, name, endpoint),
		"x-city311-direction":      endpoint.Direction,
		"x-city311-authentication": endpoint.Authentication,
	}
	if endpoint.RequiredCapability != "" {
		operation["x-city311-required-capability"] = endpoint.RequiredCapability
	}
	if endpoint.Scope != "" {
		operation["x-city311-oauth-scope"] = endpoint.Scope
	}
	if endpoint.PrivacyRule != "" {
		operation["x-city311-privacy-rule"] = endpoint.PrivacyRule
	}

	parameters := openAPIParameters(contract, endpoint)
	if len(parameters) > 0 {
		operation["parameters"] = parameters
	}
	if endpoint.RequestSchema != "" {
		contentType := "application/json"
		if name == "portal_attachment_upload" {
			contentType = "multipart/form-data"
		}
		examples := map[string]interface{}{
			"request": map[string]interface{}{"value": exampleForSchema(contract, endpoint.RequestSchema, map[string]bool{})},
		}
		for mockName, mock := range contract.Mocks {
			if mock.Endpoint == name && mock.Role == "request" {
				examples[mockName] = map[string]interface{}{"value": normalizedJSON(mock.Body)}
			}
		}
		operation["requestBody"] = map[string]interface{}{
			"required": endpoint.RequestSchema != "empty_request",
			"content": map[string]interface{}{
				contentType: map[string]interface{}{
					"schema": map[string]interface{}{"$ref": schemaRef(endpoint.RequestSchema)}, "examples": examples,
				},
			},
		}
	}
	return operation
}

func openAPIParameters(contract ContractDocument, endpoint EndpointContract) []interface{} {
	parameters := make([]interface{}, 0, len(endpoint.PathParameters)+len(endpoint.QueryParameters)+len(endpoint.RequiredHeaders))
	for _, name := range sortedSchemaPropertyNames(endpoint.PathParameters) {
		parameters = append(parameters, map[string]interface{}{
			"name": name, "in": "path", "required": true, "schema": openAPISchema(contract, endpoint.PathParameters[name]),
		})
	}
	for _, name := range sortedSchemaPropertyNames(endpoint.QueryParameters) {
		parameters = append(parameters, map[string]interface{}{
			"name": name, "in": "query", "required": false, "schema": openAPISchema(contract, endpoint.QueryParameters[name]),
		})
	}
	headers := append([]string(nil), endpoint.RequiredHeaders...)
	sort.Strings(headers)
	for _, name := range headers {
		if name == "Authorization" || name == "Content-Type" {
			continue
		}
		schema := map[string]interface{}{"type": "string"}
		if name == IfMatchHeader {
			schema["pattern"] = "^\\\"[1-9][0-9]*\\\"$"
			schema["example"] = "\"7\""
		}
		parameters = append(parameters, map[string]interface{}{
			"name": name, "in": "header", "required": true, "schema": schema,
		})
	}
	return parameters
}

func openAPIResponses(contract ContractDocument, endpointName string, endpoint EndpointContract) map[string]interface{} {
	responses := map[string]interface{}{}
	for _, label := range sortedStatusLabels(endpoint.SuccessStatuses) {
		status := endpoint.SuccessStatuses[label]
		key := strconv.Itoa(status)
		response := responseAt(responses, key, label)
		response["x-city311-example"] = map[string]interface{}{"status": status, "body": nil}
		if status != 204 && endpoint.ResponseSchema != "" && endpoint.ResponseSchema != "empty_response" {
			body := responseExampleBody(contract, endpointName, status, exampleForSchema(contract, endpoint.ResponseSchema, map[string]bool{}))
			response["x-city311-example"] = map[string]interface{}{"status": status, "body": body}
			content := responseContent(response, map[string]interface{}{"$ref": schemaRef(endpoint.ResponseSchema)})
			examples := content["examples"].(map[string]interface{})
			examples[label] = map[string]interface{}{"value": body}
			attachResponseMocks(contract, endpointName, status, false, examples)
		}
	}
	for _, code := range sortedStatusLabels(endpoint.ErrorStatuses) {
		status := endpoint.ErrorStatuses[code]
		key := strconv.Itoa(status)
		response := responseAt(responses, key, code)
		body := responseExampleBody(contract, endpointName, status, map[string]interface{}{
			"error": code, "message": humanizeError(code), "retryable": status >= 500,
		})
		response["x-city311-example"] = map[string]interface{}{"status": status, "body": body}
		content := responseContent(response, map[string]interface{}{"$ref": schemaRef("error")})
		examples := content["examples"].(map[string]interface{})
		examples[strings.ToLower(code)] = map[string]interface{}{"value": body}
		attachResponseMocks(contract, endpointName, status, true, examples)
	}
	return responses
}

func responseExampleBody(contract ContractDocument, endpointName string, status int, fallback interface{}) interface{} {
	for _, name := range sortedMockNames(contract.Mocks) {
		mock := contract.Mocks[name]
		if mock.Endpoint == endpointName && mock.Role == "response" && mock.HTTPStatus == status {
			return normalizedJSON(mock.Body)
		}
	}
	return fallback
}

func sortedMockNames(mocks map[string]MockContract) []string {
	names := make([]string, 0, len(mocks))
	for name := range mocks {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func responseAt(responses map[string]interface{}, status, description string) map[string]interface{} {
	if existing, present := responses[status].(map[string]interface{}); present {
		return existing
	}
	response := map[string]interface{}{"description": strings.ReplaceAll(description, "_", " ")}
	responses[status] = response
	return response
}

func responseContent(response map[string]interface{}, schema map[string]interface{}) map[string]interface{} {
	contentMap, _ := response["content"].(map[string]interface{})
	if contentMap == nil {
		contentMap = map[string]interface{}{}
		response["content"] = contentMap
	}
	media, _ := contentMap["application/json"].(map[string]interface{})
	if media == nil {
		media = map[string]interface{}{"schema": schema, "examples": map[string]interface{}{}}
		contentMap["application/json"] = media
	}
	return media
}

func attachResponseMocks(contract ContractDocument, endpointName string, status int, wantError bool, examples map[string]interface{}) {
	for name, mock := range contract.Mocks {
		if mock.Endpoint != endpointName || mock.Role != "response" || mock.HTTPStatus != status {
			continue
		}
		body := normalizedJSON(mock.Body)
		bodyMap, _ := body.(map[string]interface{})
		_, hasError := bodyMap["error"]
		if hasError != wantError {
			continue
		}
		examples[name] = map[string]interface{}{"value": body}
	}
}

func openAPISecurity(endpoint EndpointContract) []map[string][]string {
	scope := []string{}
	if endpoint.Scope != "" {
		scope = []string{endpoint.Scope}
	}
	switch endpoint.Authentication.Mode {
	case "none":
		return []map[string][]string{}
	case "session_cookie_optional":
		return []map[string][]string{{"sessionCookie": {}}, {}}
	case "session_cookie":
		return []map[string][]string{{"sessionCookie": {}}}
	case "oauth2_bearer", "oauth2_client_credentials":
		return []map[string][]string{{"oauth2": scope}}
	case "server_api_token":
		return []map[string][]string{{"serverApiToken": {}}}
	case "request_signature":
		return []map[string][]string{{"civicworksSignature": {}}}
	default:
		panic("unsupported OpenAPI authentication mode: " + endpoint.Authentication.Mode)
	}
}

func openAPISchema(contract ContractDocument, value interface{}) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		out := map[string]interface{}{}
		nullable, _ := typed["nullable"].(bool)
		for key, child := range typed {
			if key == "nullable" {
				continue
			}
			switch key {
			case "properties":
				properties := map[string]interface{}{}
				for name, property := range child.(map[string]interface{}) {
					properties[name] = openAPISchema(contract, property)
				}
				out["properties"] = properties
			case "schema_ref":
				out["$ref"] = schemaRef(child.(string))
			case "items_ref":
				out["items"] = map[string]interface{}{"$ref": schemaRef(child.(string))}
			case "enum_ref":
				name := child.(string)
				out["$ref"] = schemaRef("enum_" + name)
				out["x-city311-enum-ref"] = name
			case "enum_source":
				if child == "help_keys" {
					keys := make([]string, 0, len(contract.HelpKeys))
					for name := range contract.HelpKeys {
						keys = append(keys, name)
					}
					sort.Strings(keys)
					out["enum"] = keys
					out["x-city311-enum-source"] = child
				}
			default:
				out[openAPIKeyword(key)] = openAPISchema(contract, child)
			}
		}
		if nullable {
			if typeName, ok := out["type"].(string); ok {
				out["type"] = []string{typeName, "null"}
			} else {
				return map[string]interface{}{"anyOf": []interface{}{out, map[string]interface{}{"type": "null"}}}
			}
		}
		return out
	case []map[string]interface{}:
		out := make([]interface{}, 0, len(typed))
		for _, child := range typed {
			out = append(out, openAPISchema(contract, child))
		}
		return out
	case []interface{}:
		out := make([]interface{}, 0, len(typed))
		for _, child := range typed {
			out = append(out, openAPISchema(contract, child))
		}
		return out
	default:
		return value
	}
}

func openAPIKeyword(key string) string {
	standard := map[string]string{
		"min_length": "minLength", "max_length": "maxLength", "min_items": "minItems", "max_items": "maxItems",
		"unique_items": "uniqueItems", "additional_properties": "additionalProperties", "one_of": "oneOf",
		"read_only": "readOnly", "write_only": "writeOnly",
	}
	if mapped := standard[key]; mapped != "" {
		return mapped
	}
	for _, name := range []string{"type", "required", "properties", "items", "enum", "const", "format", "pattern", "minimum", "maximum", "default", "description", "examples", "example", "allOf", "anyOf", "oneOf", "additionalProperties"} {
		if key == name {
			return key
		}
	}
	return "x-city311-" + strings.ReplaceAll(key, "_", "-")
}

func exampleForSchema(contract ContractDocument, schemaName string, visiting map[string]bool) interface{} {
	if visiting[schemaName] {
		return map[string]interface{}{}
	}
	visiting[schemaName] = true
	defer delete(visiting, schemaName)
	return exampleForValue(contract, contract.Schemas[schemaName], visiting)
}

func exampleForValue(contract ContractDocument, value interface{}, visiting map[string]bool) interface{} {
	schema, ok := value.(map[string]interface{})
	if !ok {
		return nil
	}
	if ref, ok := schema["schema_ref"].(string); ok {
		return exampleForSchema(contract, ref, visiting)
	}
	if constant, present := schema["const"]; present {
		return constant
	}
	if enumRef, ok := schema["enum_ref"].(string); ok && len(contract.Enums[enumRef]) > 0 {
		return contract.Enums[enumRef][0]
	}
	if source, ok := schema["enum_source"].(string); ok && source == "help_keys" {
		keys := make([]string, 0, len(contract.HelpKeys))
		for key := range contract.HelpKeys {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		if len(keys) > 0 {
			return keys[0]
		}
	}
	if values, ok := schema["enum"].([]string); ok && len(values) > 0 {
		return values[0]
	}
	if values, ok := schema["enum"].([]interface{}); ok && len(values) > 0 {
		return values[0]
	}
	switch schema["type"] {
	case "object":
		result := map[string]interface{}{}
		properties, _ := schema["properties"].(map[string]interface{})
		required, _ := schema["required"].([]string)
		for _, name := range required {
			if property, present := properties[name]; present {
				result[name] = exampleForValue(contract, property, visiting)
			}
		}
		return result
	case "array":
		if ref, ok := schema["items_ref"].(string); ok {
			return []interface{}{exampleForSchema(contract, ref, visiting)}
		}
		if item, present := schema["items"]; present {
			return []interface{}{exampleForValue(contract, item, visiting)}
		}
		return []interface{}{}
	case "integer":
		if minimum, present := schema["minimum"]; present {
			return minimum
		}
		return 1
	case "number":
		return 1.0
	case "boolean":
		return false
	case "string":
		switch schema["pattern"] {
		case "^SR-[0-9]{4}-[0-9]{5}$":
			return "SR-2026-00041"
		case "^[a-z0-9._-]+$":
			return "resident.user"
		}
		switch schema["format"] {
		case "email":
			return "resident@example.test"
		case "date-time":
			return "2026-08-25T12:00:00Z"
		case "uri":
			return "https://city311.example.test/resource"
		case "binary":
			return "binary-content"
		}
		if _, password := schema["minimum_character_classes"]; password {
			return "ExamplePass1!"
		}
		minimum, _ := schema["min_length"].(int)
		maximum, _ := schema["max_length"].(int)
		length := len("string")
		if minimum > length {
			length = minimum
		}
		if maximum > 0 && length > maximum {
			length = maximum
		}
		if length < 1 {
			length = 1
		}
		return strings.Repeat("x", length)
	default:
		if alternatives, ok := schema["one_of"].([]map[string]interface{}); ok && len(alternatives) > 0 {
			return exampleForValue(contract, alternatives[0], visiting)
		}
		return map[string]interface{}{}
	}
}

func normalizedJSON(value interface{}) interface{} {
	raw, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	var normalized interface{}
	if err = json.Unmarshal(raw, &normalized); err != nil {
		panic(err)
	}
	return normalized
}

func schemaRef(name string) string { return "#/components/schemas/" + name }

func humanizeError(code string) string {
	return fmt.Sprintf("%s.", strings.ToLower(strings.ReplaceAll(code, "_", " ")))
}

func operationTag(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 && parts[0] == "api" && parts[1] == "v1" {
		return parts[2]
	}
	if len(parts) > 0 {
		return parts[0]
	}
	return "root"
}

func sortedEndpointNames(endpoints map[string]EndpointContract) []string {
	names := make([]string, 0, len(endpoints))
	for name := range endpoints {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func sortedSchemaPropertyNames(properties map[string]map[string]interface{}) []string {
	names := make([]string, 0, len(properties))
	for name := range properties {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func sortedStatusLabels(statuses map[string]int) []string {
	labels := make([]string, 0, len(statuses))
	for label := range statuses {
		labels = append(labels, label)
	}
	sort.Strings(labels)
	return labels
}
