package city311

import (
	"encoding/json"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestOpenAPISnapshotMatchesAuthoritativeContract(t *testing.T) {
	actualRaw, err := os.ReadFile("openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	var actual interface{}
	if err = json.Unmarshal(actualRaw, &actual); err != nil {
		t.Fatal(err)
	}
	expectedRaw, err := json.Marshal(NewOpenAPIDocument())
	if err != nil {
		t.Fatal(err)
	}
	var expected interface{}
	if err = json.Unmarshal(expectedRaw, &expected); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatal("openapi.json drifted from NewOpenAPIDocument; regenerate both BE-contract-v1 artifacts")
	}
}

func TestOpenAPICoversEveryProvidedOperationAndDeclaredExample(t *testing.T) {
	contract := NewContractDocument()
	document := NewOpenAPIDocument()
	paths := document["paths"].(map[string]interface{})
	covered := 0
	for name, endpoint := range contract.Endpoints {
		if endpoint.Direction != EndpointProvidedByCRM {
			continue
		}
		pathItem, present := paths[endpoint.Path].(map[string]interface{})
		if !present {
			t.Errorf("OpenAPI omits path for %s", name)
			continue
		}
		operation, present := pathItem[strings.ToLower(endpoint.Method)].(map[string]interface{})
		if !present || operation["operationId"] != name {
			t.Errorf("OpenAPI omits operation %s", name)
			continue
		}
		covered++
		responses := operation["responses"].(map[string]interface{})
		for _, status := range endpoint.SuccessStatuses {
			response := responses[strconv.Itoa(status)].(map[string]interface{})
			if _, present := response["x-city311-example"]; !present {
				t.Errorf("OpenAPI operation %s has no success example for %d", name, status)
			}
		}
		for _, status := range endpoint.ErrorStatuses {
			response := responses[strconv.Itoa(status)].(map[string]interface{})
			if _, present := response["x-city311-example"]; !present {
				t.Errorf("OpenAPI operation %s has no failure example for %d", name, status)
			}
		}
	}
	if covered == 0 {
		t.Fatal("OpenAPI contains no CRM-provided operations")
	}
}

func TestOpenAPIUsesStandardSchemaKeywordsAndResolvedReferences(t *testing.T) {
	document := NewOpenAPIDocument()
	raw, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	text := string(raw)
	for _, bespoke := range []string{"schema_ref", "items_ref", "enum_ref", "one_of", "min_length", "max_length", "additional_properties"} {
		if strings.Contains(text, `"`+bespoke+`"`) {
			t.Errorf("OpenAPI leaks bespoke schema keyword %s", bespoke)
		}
	}
	components := document["components"].(map[string]interface{})["schemas"].(map[string]interface{})
	walkOpenAPIReferences(t, document, components)
}

func TestCapabilitiesRoutesAndScopesAreBound(t *testing.T) {
	contract := NewContractDocument()
	for _, name := range []string{"capability", "route", "oauth_scope"} {
		if len(contract.Enums[name]) == 0 {
			t.Errorf("missing client vocabulary %s", name)
		}
	}
	actor := contract.Schemas["current_actor"]["properties"].(map[string]interface{})
	for property, enum := range map[string]string{"capabilities": "capability", "available_routes": "route", "scopes": "oauth_scope"} {
		items := actor[property].(map[string]interface{})["items"].(map[string]interface{})
		if items["enum_ref"] != enum {
			t.Errorf("current_actor.%s is not bound to %s", property, enum)
		}
	}
	for name, endpoint := range contract.Endpoints {
		if endpoint.Direction == EndpointProvidedByCRM && endpoint.Authentication.Mode == "session_cookie" {
			if endpoint.RequiredCapability == "" || !contains(contract.Enums["capability"], endpoint.RequiredCapability) {
				t.Errorf("protected endpoint %s has no published capability", name)
			}
		}
	}
}

func walkOpenAPIReferences(t *testing.T, value interface{}, components map[string]interface{}) {
	t.Helper()
	switch typed := value.(type) {
	case map[string]interface{}:
		for key, child := range typed {
			if key == "$ref" {
				ref, _ := child.(string)
				const prefix = "#/components/schemas/"
				if strings.HasPrefix(ref, prefix) {
					if _, present := components[strings.TrimPrefix(ref, prefix)]; !present {
						t.Errorf("OpenAPI contains dangling schema reference %s", ref)
					}
				}
			}
			walkOpenAPIReferences(t, child, components)
		}
	case []interface{}:
		for _, child := range typed {
			walkOpenAPIReferences(t, child, components)
		}
	}
}
