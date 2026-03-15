package automation

import (
	"context"
	"encoding/json"
	"strings"

	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
)

type (
	jsonParserHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}

	jsonParserHandler struct {
		reg jsonParserHandlerRegistry
	}

	jsonParserParseArgs struct {
		hasJsonText bool
		JsonText    string

		hasFields bool
		Fields    []string
	}

	jsonParserParseResults struct {
		Data  map[string]interface{}
		Error string
	}
)

func JsonParserHandler(reg jsonParserHandlerRegistry) *jsonParserHandler {
	h := &jsonParserHandler{
		reg: reg,
	}

	h.register()
	return h
}

func (h jsonParserHandler) register() {
	h.reg.AddFunctions(
		h.Parse(),
	)
}

func (h jsonParserHandler) parse(ctx context.Context, args *jsonParserParseArgs) (*jsonParserParseResults, error) {
	r := &jsonParserParseResults{
		Data: make(map[string]interface{}),
	}

	if args.JsonText == "" {
		r.Error = "jsonText argument is required and cannot be empty"
		return r, nil
	}

	// Trim and parse JSON text
	jsonText := strings.TrimSpace(args.JsonText)
	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(jsonText), &parsed)
	if err != nil {
		// Try to complete partial JSON (using a simplified version of bifrost's parser)
		completedJson := completePartialJSON(jsonText)
		err = json.Unmarshal([]byte(completedJson), &parsed)
		if err != nil {
			r.Error = err.Error()
			return r, nil
		}
	}

	if args.hasFields && len(args.Fields) > 0 {
		// Extract requested fields
		for _, field := range args.Fields {
			field = strings.TrimSpace(field)
			if field == "" {
				continue
			}

			// Handle nested fields using dot notation
			value, exists := getNestedField(parsed, strings.Split(field, "."))
			if exists {
				r.Data[field] = value
			}
		}
	} else {
		// Return entire parsed JSON
		r.Data = parsed
	}

	return r, nil
}

// completePartialJSON completes partial JSON strings (simplified implementation)
func completePartialJSON(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "{}"
	}

	// Quick check: if it starts with { or [, it might be JSON
	if s[0] != '{' && s[0] != '[' {
		return s
	}

	var stack []byte
	inString := false
	escaped := false

	for i := 0; i < len(s); i++ {
		char := s[i]

		if escaped {
			escaped = false
			continue
		}

		if char == '\\' {
			escaped = true
			continue
		}

		if char == '"' {
			inString = !inString
			continue
		}

		if inString {
			continue
		}

		switch char {
		case '{', '[':
			if char == '{' {
				stack = append(stack, '}')
			} else {
				stack = append(stack, ']')
			}
		case '}', ']':
			if len(stack) > 0 && stack[len(stack)-1] == char {
				stack = stack[:len(stack)-1]
			}
		}
	}

	// Close any unclosed strings
	if inString {
		s += "\""
	}

	// Add closing characters in reverse order
	for i := len(stack) - 1; i >= 0; i-- {
		s += string(stack[i])
	}

	return s
}

// getNestedField retrieves a field from nested map using path
func getNestedField(data map[string]interface{}, path []string) (interface{}, bool) {
	current := interface{}(data)

	for _, key := range path {
		switch v := current.(type) {
		case map[string]interface{}:
			val, exists := v[key]
			if !exists {
				return nil, false
			}
			current = val
		default:
			return nil, false
		}
	}

	return current, true
}
