package envoyx

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"strings"
)

// SplitResourceIdentifier takes an identifier string and splices it into path
// identifiers as defined by the resource
func SplitResourceIdentifier(ref string) (out map[string]Ref) {
	out = make(map[string]Ref, 3)

	ref = strings.TrimRight(ref, "/")
	pp := strings.Split(ref, "/")
	rt := pp[0]
	pp = pp[1:]

	// @todo this is to support res. tr. resources also.
	//       Split it into a separate function and remove this.
	if !strings.HasPrefix(rt, "corteza::") {
		rt = "corteza::" + rt
	}

	gRef := func(pp []string, i int) string {
		if pp[i] == "*" {
			return ""
		}
		return pp[i]
	}

	switch rt {

	case "corteza::system:apigw-filter":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:apigw-filter",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:apigw-route":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:apigw-route",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:application":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:application",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:attachment":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:attachment",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:auth-client":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:auth-client",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:auth-confirmed-client":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:auth-confirmed-client",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:auth-oa2token":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:auth-oa2token",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:auth-session":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:auth-session",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:credential":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:credential",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:dal-connection":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:dal-connection",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:dal-schema-alteration":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:dal-schema-alteration",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:dal-sensitivity-level":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:dal-sensitivity-level",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:data-privacy-request":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:data-privacy-request",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:data-privacy-request-comment":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:data-privacy-request-comment",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:queue":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:queue",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:queue-message":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:queue-message",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:reminder":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:reminder",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:report":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:report",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:resource-translation":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:resource-translation",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:role":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:role",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:role-member":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:role-member",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:settings":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:settings",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:template":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:template",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::system:user":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::system:user",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::compose:attachment":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::compose:attachment",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::compose:chart":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		aux := gRef(pp, 0)
		if aux != "" {
			scope.ResourceType = "corteza::compose:namespace"
			scope.Identifiers = MakeIdentifiers(aux)
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::compose:namespace",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

		if gRef(pp, 1) == "" {
			return
		}

		out["Path.1"] = Ref{
			ResourceType: "corteza::compose:chart",
			Identifiers:  MakeIdentifiers(gRef(pp, 1)),
			Scope:        scope,
		}

	case "corteza::compose:module":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		aux := gRef(pp, 0)
		if aux != "" {
			scope.ResourceType = "corteza::compose:namespace"
			scope.Identifiers = MakeIdentifiers(aux)
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::compose:namespace",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

		if gRef(pp, 1) == "" {
			return
		}

		out["Path.1"] = Ref{
			ResourceType: "corteza::compose:module",
			Identifiers:  MakeIdentifiers(gRef(pp, 1)),
			Scope:        scope,
		}

	case "corteza::compose:module-field":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		aux := gRef(pp, 0)
		if aux != "" {
			scope.ResourceType = "corteza::compose:namespace"
			scope.Identifiers = MakeIdentifiers(aux)
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::compose:namespace",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

		if gRef(pp, 1) == "" {
			return
		}

		out["Path.1"] = Ref{
			ResourceType: "corteza::compose:module",
			Identifiers:  MakeIdentifiers(gRef(pp, 1)),
			Scope:        scope,
		}

		if gRef(pp, 2) == "" {
			return
		}

		out["Path.2"] = Ref{
			ResourceType: "corteza::compose:module-field",
			Identifiers:  MakeIdentifiers(gRef(pp, 2)),
			Scope:        scope,
		}

	case "corteza::compose:namespace":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		scope.ResourceType = "corteza::compose:namespace"
		scope.Identifiers = MakeIdentifiers(gRef(pp, 0))

		out["Path.0"] = Ref{
			ResourceType: "corteza::compose:namespace",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::compose:page":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		aux := gRef(pp, 0)
		if aux != "" {
			scope.ResourceType = "corteza::compose:namespace"
			scope.Identifiers = MakeIdentifiers(aux)
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::compose:namespace",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

		if gRef(pp, 1) == "" {
			return
		}

		out["Path.1"] = Ref{
			ResourceType: "corteza::compose:page",
			Identifiers:  MakeIdentifiers(gRef(pp, 1)),
			Scope:        scope,
		}

	case "corteza::compose:page-layout":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		aux := gRef(pp, 0)
		if aux != "" {
			scope.ResourceType = "corteza::compose:namespace"
			scope.Identifiers = MakeIdentifiers(aux)
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::compose:namespace",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

		if gRef(pp, 1) == "" {
			return
		}

		out["Path.1"] = Ref{
			ResourceType: "corteza::compose:page",
			Identifiers:  MakeIdentifiers(gRef(pp, 1)),
			Scope:        scope,
		}

		if gRef(pp, 2) == "" {
			return
		}

		out["Path.2"] = Ref{
			ResourceType: "corteza::compose:page-layout",
			Identifiers:  MakeIdentifiers(gRef(pp, 2)),
			Scope:        scope,
		}

	case "corteza::compose:record":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		aux := gRef(pp, 0)
		if aux != "" {
			scope.ResourceType = "corteza::compose:namespace"
			scope.Identifiers = MakeIdentifiers(aux)
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::compose:namespace",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

		if gRef(pp, 1) == "" {
			return
		}

		out["Path.1"] = Ref{
			ResourceType: "corteza::compose:module",
			Identifiers:  MakeIdentifiers(gRef(pp, 1)),
			Scope:        scope,
		}

		if gRef(pp, 2) == "" {
			return
		}

		out["Path.2"] = Ref{
			ResourceType: "corteza::compose:record",
			Identifiers:  MakeIdentifiers(gRef(pp, 2)),
			Scope:        scope,
		}

	case "corteza::compose:record-revision":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::compose:record-revision",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::automation:session":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::automation:session",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::automation:trigger":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::automation:trigger",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	case "corteza::automation:workflow":
		scope := Scope{}

		if gRef(pp, 0) == "" {
			return
		}

		out["Path.0"] = Ref{
			ResourceType: "corteza::automation:workflow",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
			Scope:        scope,
		}

	}

	return
}
