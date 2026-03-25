package automation

import (
	"context"
	"encoding/json"
	"fmt"
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
	}

	jsonParserParseResults struct {
		Result interface{}
		Error  string
	}

	jsonParserStringifyArgs struct {
		hasJsonObject bool
		JsonObject    interface{}
	}

	jsonParserStringifyResults struct {
		Result string
		Error  string
	}

	jsonParserTemplateArgs struct {
		hasVars  bool
		Template string
		Vars     map[string]interface{}
	}

	jsonParserTemplateResults struct {
		Result string
		Error  string
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
		h.Stringify(),
		h.Template(),
	)
}

// parse implements the jsonParse function
func (h jsonParserHandler) parse(ctx context.Context, args *jsonParserParseArgs) (*jsonParserParseResults, error) {
	r := &jsonParserParseResults{}

	if args.JsonText == "" {
		r.Error = "jsonText is empty"
		return r, nil
	}

	var result interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(args.JsonText)), &result); err != nil {
		r.Error = err.Error()
		return r, nil
	}

	r.Result = result
	return r, nil
}

// Smart stringify: if input is already valid JSON string, return as-is
func (h jsonParserHandler) stringify(ctx context.Context, args *jsonParserStringifyArgs) (*jsonParserStringifyResults, error) {

	r := &jsonParserStringifyResults{}

	if args.JsonObject == nil {
		r.Error = "jsonObject is empty"
		return r, nil
	}

	// If the value is already a string, check if it's valid JSON
	if s, ok := args.JsonObject.(string); ok {
		// Trim whitespace and try to parse
		trimmed := strings.TrimSpace(s)
		if len(trimmed) == 0 {
			r.Error = "jsonObject is empty"
			return r, nil
		}

		// Preprocess JS-like object literals into valid JSON
		if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
			trimmed = preprocessJSObjectToJSON(trimmed)
		}

		var probe interface{}
		if err := json.Unmarshal([]byte(trimmed), &probe); err == nil {
			// Already valid JSON string - return as-is (no extra encoding)
			r.Result = trimmed
			return r, nil
		}
		// Plain text string - convert to JSON string (wrap in quotes and escape)
		b, err := json.Marshal(s)
		if err != nil {
			r.Error = err.Error()
			return r, nil
		}
		r.Result = string(b)
		return r, nil
	}

	// Not a string (object, array, number, bool, null) - marshal normally
	b, err := json.Marshal(args.JsonObject)
	if err != nil {
		r.Error = err.Error()
		return r, nil
	}

	r.Result = string(b)
	return r, nil
}

// substituteVariables replaces ${varName} with values from the vars map
// Supports dot-notation for nested values (e.g., ${var.filename}, ${var.data.content})
// Also supports any prefix like ${jsonVars.filename}, ${data.content} etc.
func substituteVariables(s string, vars map[string]interface{}) string {
	var result strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '$' && i+1 < len(s) && s[i+1] == '{' {
			// Find closing brace
			end := strings.Index(s[i+2:], "}")
			if end != -1 {
				placeholder := s[i+2 : i+2+end]

				// Strip any leading "word." prefix (e.g. var., jsonVars., data.)
				// Use the part after the FIRST dot as the lookup key,
				// but preserve dot-notation for nested maps
				key := placeholder
				if dotIdx := strings.Index(placeholder, "."); dotIdx != -1 {
					key = placeholder[dotIdx+1:]
				}

				if val := getNestedValue(vars, key); val != nil {
					// Convert value to string and sanitize for JSON
					strVal := fmt.Sprintf("%v", val)
					result.WriteString(sanitizeForJSON(strVal))
				} else {
					// Variable not found - leave as-is
					result.WriteString("${" + placeholder + "}")
				}
				i += end + 3 // Skip ${...}
				continue
			}
		}
		result.WriteByte(s[i])
		i++
	}
	return result.String()
}

// getNestedValue retrieves a value from a map using dot-notation (e.g., "filename" or "data.content")
func getNestedValue(vars map[string]interface{}, key string) interface{} {
	parts := strings.Split(key, ".")
	var current interface{} = vars

	for _, part := range parts {
		if currentMap, ok := current.(map[string]interface{}); ok {
			current = currentMap[part]
			if current == nil {
				return nil
			}
		} else {
			return nil
		}
	}
	return current
}

// sanitizeForJSON escapes special characters for use in JSON strings
func sanitizeForJSON(s string) string {
	var result strings.Builder
	for _, c := range s {
		switch c {
		case '"':
			result.WriteString(`\"`)
		case '\\':
			result.WriteString(`\\`)
		case '\n':
			result.WriteString(`\n`)
		case '\r':
			// skip
		case '\t':
			result.WriteString(`\t`)
		default:
			result.WriteRune(c)
		}
	}
	return result.String()
}

// preprocessJSObjectToJSON converts JavaScript-like object literals to valid JSON
func preprocessJSObjectToJSON(s string) string {
	// Step 1: Replace backtick strings with proper JSON double-quoted strings
	s = replaceBacktickStrings(s)

	// Step 2: Fix bare backslashes (e.g. openai\gpt -> openai\\gpt)
	// Only fix backslashes that are NOT valid JSON escape sequences
	s = fixBareBackslashes(s)

	return s
}

// replaceBacktickStrings converts backtick strings to JSON double-quoted strings
func replaceBacktickStrings(s string) string {
	var result strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '`' {
			// Find closing backtick
			result.WriteByte('"')
			i++
			for i < len(s) && s[i] != '`' {
				switch s[i] {
				case '\n':
					result.WriteString(`\n`)
				case '\r':
					// skip \r, handle \r\n as just \n
				case '"':
					result.WriteString(`\"`)
				case '\\':
					result.WriteString(`\\`)
				default:
					result.WriteByte(s[i])
				}
				i++
			}
			result.WriteByte('"')
			i++ // skip closing backtick
		} else {
			result.WriteByte(s[i])
			i++
		}
	}
	return result.String()
}

// fixBareBackslashes escapes bare backslashes that aren't valid JSON escapes
func fixBareBackslashes(s string) string {
	// Valid JSON escape chars after backslash
	validEscapes := map[byte]bool{
		'"': true, '\\': true, '/': true, 'b': true,
		'f': true, 'n': true, 'r': true, 't': true, 'u': true,
	}

	var result strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\\' && i+1 < len(s) {
			if validEscapes[s[i+1]] {
				// Valid escape sequence - keep as-is
				result.WriteByte(s[i])
			} else {
				// Bare backslash - escape it
				result.WriteString(`\\`)
			}
		} else {
			result.WriteByte(s[i])
		}
		i++
	}
	return result.String()
}

// template implements the jsonTemplate function with variable substitution
func (h jsonParserHandler) template(ctx context.Context, args *jsonParserTemplateArgs) (*jsonParserTemplateResults, error) {
	r := &jsonParserTemplateResults{}

	if args.Template == "" {
		r.Error = "template is empty"
		return r, nil
	}

	// If no vars provided, just return the template as-is
	if args.Vars == nil {
		r.Result = args.Template
		return r, nil
	}

	// Perform variable substitution
	result := substituteVariables(args.Template, args.Vars)
	r.Result = result

	return r, nil
}
