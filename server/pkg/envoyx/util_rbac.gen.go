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

	gRef := func(pp []string, i int) string {
		if pp[i] == "*" {
			return ""
		}
		return pp[i]
	}

	switch rt {

	case "corteza::system:apigwFilter":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:apigw-filter",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:apigwRoute":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:apigw-route",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:application":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:application",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:attachment":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:attachment",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:authClient":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:auth-client",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:authConfirmedClient":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:auth-confirmed-client",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:authOa2token":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:auth-oa2token",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:authSession":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:auth-session",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:credential":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:credential",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:dalConnection":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:dal-connection",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:dalSensitivityLevel":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:dal-sensitivity-level",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:dataPrivacyRequest":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:data-privacy-request",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:dataPrivacyRequestComment":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:data-privacy-request-comment",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:queue":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:queue",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:queueMessage":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:queue-message",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:reminder":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:reminder",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:report":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:report",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:resourceTranslation":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:resource-translation",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:role":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:role",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:roleMember":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:role-member",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:settingValue":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:settings",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:template":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:template",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::system:user":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::system:user",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::compose:attachment":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::compose:attachment",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::compose:chart":

		if gRef(pp, 0) == "" {
			return
		}

		if gRef(pp, 1) == "" {
			return
		}
		out["1"] = Ref{
			ResourceType: "corteza::compose:chart",
			Identifiers:  MakeIdentifiers(gRef(pp, 1)),
		}

	case "corteza::compose:module":

		if gRef(pp, 0) == "" {
			return
		}

		if gRef(pp, 1) == "" {
			return
		}
		out["1"] = Ref{
			ResourceType: "corteza::compose:module",
			Identifiers:  MakeIdentifiers(gRef(pp, 1)),
		}

	case "corteza::compose:moduleField":

		if gRef(pp, 0) == "" {
			return
		}

		if gRef(pp, 1) == "" {
			return
		}

		if gRef(pp, 2) == "" {
			return
		}
		out["2"] = Ref{
			ResourceType: "corteza::compose:module-field",
			Identifiers:  MakeIdentifiers(gRef(pp, 2)),
		}

	case "corteza::compose:namespace":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::compose:namespace",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::compose:page":

		if gRef(pp, 0) == "" {
			return
		}

		if gRef(pp, 1) == "" {
			return
		}
		out["1"] = Ref{
			ResourceType: "corteza::compose:page",
			Identifiers:  MakeIdentifiers(gRef(pp, 1)),
		}

	case "corteza::compose:record":

		if gRef(pp, 0) == "" {
			return
		}

		if gRef(pp, 1) == "" {
			return
		}

		if gRef(pp, 2) == "" {
			return
		}
		out["2"] = Ref{
			ResourceType: "corteza::compose:record",
			Identifiers:  MakeIdentifiers(gRef(pp, 2)),
		}

	case "corteza::compose:recordRevision":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::compose:record-revision",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::automation:session":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::automation:session",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::automation:trigger":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::automation:trigger",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::automation:workflow":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::automation:workflow",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::federation:exposedModule":

		if gRef(pp, 0) == "" {
			return
		}

		if gRef(pp, 1) == "" {
			return
		}
		out["1"] = Ref{
			ResourceType: "corteza::federation:exposed-module",
			Identifiers:  MakeIdentifiers(gRef(pp, 1)),
		}

	case "corteza::federation:moduleMapping":

		if gRef(pp, 0) == "" {
			return
		}

		if gRef(pp, 1) == "" {
			return
		}
		out["1"] = Ref{
			ResourceType: "corteza::federation:module-mapping",
			Identifiers:  MakeIdentifiers(gRef(pp, 1)),
		}

	case "corteza::federation:node":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::federation:node",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::federation:nodeSync":

		if gRef(pp, 0) == "" {
			return
		}
		out["0"] = Ref{
			ResourceType: "corteza::federation:node-sync",
			Identifiers:  MakeIdentifiers(gRef(pp, 0)),
		}

	case "corteza::federation:sharedModule":

		if gRef(pp, 0) == "" {
			return
		}

		if gRef(pp, 1) == "" {
			return
		}
		out["1"] = Ref{
			ResourceType: "corteza::federation:shared-module",
			Identifiers:  MakeIdentifiers(gRef(pp, 1)),
		}

	}

	return
}
