package compose

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

component: schema.#component & {
	handle: "compose"

	resources: {
		"attachment":                  attachment
		"chart":                       chart
		"city311-actor-profile":       actorProfile
		"city311-audit-event":         auditEvent
		"city311-idempotency-record":  idempotencyRecord
		"city311-public-history-item": publicHistoryItem
		"city311-request-attachment":  requestAttachment
		"city311-request-sequence":    requestSequence
		"city311-service-request":     serviceRequest
		"module":                      module
		"module-field":                moduleField
		"namespace":                   namespace
		"page":                        page
		"page-layout":                 pageLayout
		"record":                      record
		"record-revision":             record_revision
	}

	rbac: operations: {
		"settings.read": description:                "Read settings"
		"settings.manage": description:              "Manage settings"
		"namespace.create": description:             "Create namespace"
		"namespaces.search": description:            "List, search or filter namespaces"
		"resource-translations.manage": description: "List, search, create, or update resource translations"
	}
}
