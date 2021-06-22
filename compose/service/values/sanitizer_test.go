package values

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/cortezaproject/corteza-server/compose/types"
)

func Test_sanitizer_Run(t *testing.T) {
	tests := []struct {
		name    string
		kind    string
		options map[string]interface{}
		input   string
		output  string
		outref  uint64
	}{
		{
			name:   "numbers should be trimmed",
			kind:   "Number",
			input:  " 42 ",
			output: "42",
		},
		{
			name:   "object reference should be processed",
			kind:   "Record",
			input:  " 133569629112020995 ",
			output: "133569629112020995",
			outref: 133569629112020995,
		},
		{
			name:   "object reference should be numeric",
			kind:   "Record",
			input:  " foo ",
			output: "",
		},
		{
			name:   "user reference should be processed",
			kind:   "User",
			input:  " 133569629112020995 ",
			output: "133569629112020995",
			outref: 133569629112020995,
		},
		{
			name:   "user reference should be numeric",
			kind:   "User",
			input:  " foo ",
			output: "",
		},
		{
			name:   "strings should be kept intact",
			kind:   "String",
			input:  " The answer ",
			output: " The answer ",
		},
		{
			name:   "booleans should be converted (t)",
			kind:   "Bool",
			input:  "t",
			output: "1",
		},
		{
			name:   "booleans should be converted (false)",
			kind:   "Bool",
			input:  "false",
			output: "0",
		},
		{
			name:   "booleans should be converted (garbage)",
			kind:   "Bool",
			input:  "%%#)%)')$)'",
			output: "0",
		},
		{
			name:   "dates should be converted to ISO",
			kind:   "DateTime",
			input:  "Mon Jan 2 15:04:05 2006",
			output: "2006-01-02T15:04:05Z",
		},
		{
			name:   "dates should be converted to UTC",
			kind:   "DateTime",
			input:  "2020-03-02T20:20:20+05:00",
			output: "2020-03-02T15:20:20Z",
		},
		{
			name:   "micro/mili seconds should be cut off",
			kind:   "DateTime",
			input:  "2020-03-11T11:20:08.471Z",
			output: "2020-03-11T11:20:08Z",
		},
		{
			name:   "number space trim",
			kind:   "Number",
			input:  "  42 ",
			output: "42",
		},
		{
			name:   "number negative",
			kind:   "Number",
			input:  "-42",
			output: "-42",
		},
		{
			name:   "number positive",
			kind:   "Number",
			input:  "+42",
			output: "42",
		},
		{
			name:    "number precision",
			kind:    "Number",
			options: map[string]interface{}{"precision": 3},
			input:   "42.44455",
			output:  "42.445",
		},
		{
			name:    "number precision; bigger precision then provided fracture",
			kind:    "Number",
			options: map[string]interface{}{"precision": 3},
			input:   "42.4",
			output:  "42.4",
		},
		{
			name:    "number precision; default to 0",
			kind:    "Number",
			options: map[string]interface{}{},
			input:   "42.4",
			output:  "42",
		},
		{
			name:    "number precision; clamped between [0, 6]",
			kind:    "Number",
			options: map[string]interface{}{"precision": 10},
			input:   "42.5555555555",
			output:  "42.555556",
		},
		{
			name:    "number precision; round",
			kind:    "Number",
			options: map[string]interface{}{"precision": 0},
			input:   "41.6",
			output:  "42",
		},
		{
			name:    "number precision; trailing .00",
			kind:    "Number",
			options: map[string]interface{}{"precision": 2},
			input:   "42.00",
			output:  "42",
		},
		{
			name:    "number precision; trailing .040",
			kind:    "Number",
			options: map[string]interface{}{"precision": 2},
			input:   "42.040",
			output:  "42.04",
		},
		{
			name:    "string escaping; html",
			kind:    "String",
			options: map[string]interface{}{},
			input:   "<span onerror=alert()>Title here</span>",
			output:  "<span>Title here</span>",
		},
		{
			name:    "string escaping; html a.href with javascript alert",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a>`,
			output:  "XSS",
		},
		{
			name:    "string escaping; a.href with javascript",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<a href="javascript:document.location='https://cortezaproject.org/'">XSS</A>`,
			output:  "XSS",
		},
		{
			name:    "string escaping; script with script",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<script>document.write("<scri");</script>pt src="https://cortezaproject.org/script.js"></script>`,
			output:  "pt src=&#34;https://cortezaproject.org/script.js&#34;&gt;",
		},
		{
			name:    "string escaping; script with a",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<script a=">'>" src="https://cortezaproject.org/xss.js"></script>`,
			output:  "",
		},
		{
			name:    "string escaping; meta with script",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<meta http-equiv="set-cookie" content="<script>alert('xss')</script>">`,
			output:  "",
		},
		{
			name:    "string escaping; object",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<object type="text/x-scriptlet" data="https://cortezaproject.org/xss.html"></object>`,
			output:  "",
		},
		{
			name:    "string escaping; base href",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<base href="javascript:alert('xss');//">`,
			output:  "",
		},
		{
			name:    "string escaping; script",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<!--[if gte ie 4]><script>alert('xss');</script><![endif]-->`,
			output:  "",
		},
		{
			name:    "string escaping; div",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<div style="background-image: url(javascript:alert('XSS'))"></div>`,
			output:  "<div></div>",
		},
		{
			name:    "string escaping; frameset",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<frameset><frame src="javascript:alert('XSS');"></frameset>`,
			output:  "",
		},
		{
			name:    "string escaping; iframe",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<iframe src=# onmouseover="alert(document.cookie)"></iframe>`,
			output:  "",
		},
		{
			name:    "string escaping; meta",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<meta http-equiv="refresh" content="0; url=https://;url=javascript:alert('xss');">`,
			output:  "",
		},
		{
			name:    "string escaping; br",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<br size="&{alert('XSS')}">`,
			output:  "<br>",
		},
		{
			name:    "string escaping; bgsound",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<bgsound src="javascript:alert('XSS');">`,
			output:  "",
		},
		{
			name:    "string escaping; input type image",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<input type="image" src="javascript:alert('XSS');">`,
			output:  "",
		},
		{
			name:    "string escaping; style",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<style>@import 'https://cortezaproject.org/xss.css';</style>`,
			output:  "",
		},
		{
			name:    "string escaping; link",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<link rel="stylesheet" href="javascript:alert('xss');">`,
			output:  "",
		},
		{
			name:    "string escaping; html body onload event",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<body onload=alert('XSS')>`,
			output:  "",
		},
		{
			name:    "string escaping; xss element",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `'';!--"<XSS>=&{()}`,
			output:  "&#39;&#39;;!--&#34;=&amp;{()}",
		},
		{
			name:    "string escaping; xss element",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `Hello <span class="><script src='https://cortezaproject.org/XSS.js'></script>">there</span> world.`,
			output:  "Hello <span>there</span> world.",
		},
		{
			name:    "string escaping; xss element",
			kind:    "String",
			options: map[string]interface{}{},
			input:   `<tag1>cor<tag2></tag2>teza</tag1><tag1>server</tag1><tag2>123</tag2>`,
			output:  "cortezaserver123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := sanitizer{}
			m := &types.Module{Fields: types.ModuleFieldSet{&types.ModuleField{Name: "testField", Kind: tt.kind, Options: tt.options}}}
			v := types.RecordValueSet{&types.RecordValue{Name: "testField", Value: tt.input}}
			o := types.RecordValueSet{&types.RecordValue{Name: "testField", Value: tt.output, Ref: tt.outref}}

			// Need to mark values as updated to trigger sanitization.
			v.SetUpdatedFlag(true)
			o.SetUpdatedFlag(true)
			if sanitized := s.Run(m, v); !reflect.DeepEqual(sanitized, o) {
				t.Errorf("\ninput value:\n%v\n\nresult of sanitization:\n%v\n\nexpected:\n%v\n", tt.input, sanitized, o)
			}
		})
	}
}

func TestSanitizerExpr(t *testing.T) {
	tests := []struct {
		name   string
		kind   string
		expr   []string
		input  string
		output string
		outref uint64
	}{
		{
			name:   "cap numbers",
			kind:   "Number",
			expr:   []string{`value > 42 ? 42 : value`},
			input:  "43",
			output: "42",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := sanitizer{}
			m := &types.Module{Fields: types.ModuleFieldSet{&types.ModuleField{Name: "testField", Kind: tt.kind}}}
			m.Fields.FindByName("testField").Expressions.Sanitizers = tt.expr
			v := types.RecordValueSet{&types.RecordValue{Name: "testField", Value: tt.input}}
			o := types.RecordValueSet{&types.RecordValue{Name: "testField", Value: tt.output, Ref: tt.outref}}

			// Need to mark values as updated to trigger sanitization.
			v.SetUpdatedFlag(true)
			o.SetUpdatedFlag(true)
			if sanitized := s.Run(m, v); !reflect.DeepEqual(sanitized, o) {
				t.Errorf("\ninput value:\n%v\n\nresult of sanitization:\n%v\n\nexpected:\n%v\n", tt.input, sanitized, o)
			}
		})
	}
}

func TestDatetimeSanitizer(t *testing.T) {
	tests := []struct {
		input              interface{}
		onlyDate, onlyTime bool
		rval               string
	}{
		{time.Date(1999, 9, 9, 9, 9, 9, 9, time.UTC), false, false, "1999-09-09T09:09:09Z"},
		{"2021-03-23T20:21:15Z", false, false, "2021-03-23T20:21:15Z"},
		{"2021-03-23T20:21:15+01:00", false, false, "2021-03-23T19:21:15Z"},
		{"2021-03-23", true, false, "2021-03-23"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.input), func(t *testing.T) {
			assert.New(t).Equal(tt.rval, sDatetime(tt.input, tt.onlyDate, tt.onlyTime))
		})
	}
}
