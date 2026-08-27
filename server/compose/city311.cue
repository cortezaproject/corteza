package compose

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

serviceRequest: {
	model: {
		ident:            "compose_city311_service_request"
		omitGetterSetter: true
		attributes: {
			id: schema.IdField
			request_number: {goType: "string", sortable: true, dal: {type: "Text", length: 16}}
			summary: {goType: "string", sortable: true, dal: {type: "Text", length: 160}}
			description: {goType: "string", dal: {type: "Text"}}
			service_type: {goType: "types.ServiceType", sortable: true, dal: {type: "Text", length: 64}}
			owning_department: {goType: "types.DepartmentCode", sortable: true, dal: {type: "Text", length: 64}}
			council_district: {goType: "types.DistrictCode", sortable: true, dal: {type: "Text", length: 32}}
			source_channel: {goType: "types.SourceChannel", sortable: true, dal: {type: "Text", length: 32}}
			origin_class: {goType: "types.OriginClass", sortable: true, dal: {type: "Text", length: 16}}
			status: {goType: "types.ServiceRequestStatus", sortable: true, dal: {type: "Text", length: 32}}
			primary_requester: {goType: "types.City311JSON", dal: {type: "JSON", defaultEmptyObject: true}}
			location: {goType: "types.City311JSON", dal: {type: "JSON", defaultEmptyObject: true}}
			custom_fields: {goType: "types.City311JSON", dal: {type: "JSON", defaultEmptyObject: true}}
			primary_assignee_id: {ident: "primaryAssigneeID", goType: "uint64", sortable: true, dal: {type: "ID", default: 0}}
			collaborator_ids: {ident: "collaboratorIDs", goType: "types.City311Uint64Set", dal: {type: "JSON"}}
			duplicate_group_id: {ident: "duplicateGroupID", goType: "string", sortable: true, dal: {type: "Text", length: 64}}
			version: {goType: "int", sortable: true, dal: {type: "Number", meta: {"rdbms:type": "integer"}, default: 1}}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNowField
		}
		indexes: {
			primary: {attribute: "id"}
			unique_request_number: {attribute: "request_number"}
			request_scope: {attributes: ["owning_department", "council_district", "status"]}
		}
	}
	filter: {
		struct: {
			request_number: {goType: "string"}
			service_type: {goType: "string"}
			owning_department: {goType: "string"}
			council_district: {goType: "string"}
			source_channel: {goType: "string"}
			origin_class: {goType: "string"}
			status: {goType: "string"}
			primary_assignee_id: {ident: "primaryAssigneeID", goType: "uint64"}
		}
		byValue: ["request_number", "service_type", "owning_department", "council_district", "source_channel", "origin_class", "status", "primary_assignee_id"]
	}
	features: {labels: false, flags: false}
	envoy: {omit: true}
	store: {
		ident: "city311ServiceRequest"
		api: lookups: [
			{fields: ["id"]},
			{fields: ["request_number"], constraintCheck: true},
		]
	}
}

requestSequence: {
	model: {
		ident:            "compose_city311_request_sequence"
		omitGetterSetter: true
		attributes: {
			id: schema.IdField
			next_number: {goType: "uint64", dal: {type: "Number", meta: {"rdbms:type": "bigint"}}}
		}
		indexes: {primary: {attribute: "id"}}
	}
	filter: {struct: {}}
	features: {labels: false, flags: false}
	envoy: {omit: true}
	store: {ident: "city311RequestSequence", api: lookups: [{fields: ["id"]}]}
}

idempotencyRecord: {
	model: {
		ident:            "compose_city311_idempotency"
		omitGetterSetter: true
		attributes: {
			id: schema.IdField
			operation: {goType: "string", dal: {type: "Text", length: 96}}
			key_hash: {goType: "string", dal: {type: "Text", length: 64}}
			request_hash: {goType: "string", dal: {type: "Text", length: 64}}
			response_status: {goType: "int", dal: {type: "Number", meta: {"rdbms:type": "integer"}}}
			response_body: {goType: "types.City311JSON", dal: {type: "JSON", defaultEmptyObject: true}}
			request_id: {ident: "requestID", goType: "uint64", dal: {type: "ID"}}
			created_at: schema.SortableTimestampNowField
			expires_at: schema.SortableTimestampField
		}
		indexes: {
			primary: {attribute: "id"}
			unique_operation_key: {attributes: ["operation", "key_hash"]}
			expires: {attribute: "expires_at"}
		}
	}
	filter: {
		struct: {operation: {goType: "string"}, key_hash: {goType: "string"}}
		byValue: ["operation", "key_hash"]
	}
	features: {labels: false, flags: false}
	envoy: {omit: true}
	store: {
		ident: "city311IdempotencyRecord"
		api: lookups: [
			{fields: ["id"]},
			{fields: ["operation", "key_hash"], constraintCheck: true},
		]
	}
}

auditEvent: {
	model: {
		ident:            "compose_city311_audit_event"
		omitGetterSetter: true
		attributes: {
			id: schema.IdField
			request_id: {ident: "requestID", goType: "uint64", sortable: true, dal: {type: "ID"}}
			event_type: {goType: "string", sortable: true, dal: {type: "Text", length: 96}}
			actor_type: {goType: "types.AuditActorType", dal: {type: "Text", length: 32}}
			actor_id: {ident: "actorID", goType: "uint64", dal: {type: "ID"}}
			source_channel: {goType: "types.SourceChannel", dal: {type: "Text", length: 32}}
			before: {goType: "types.City311JSON", dal: {type: "JSON", defaultEmptyObject: true}}
			after: {goType: "types.City311JSON", dal: {type: "JSON", defaultEmptyObject: true}}
			created_at: schema.SortableTimestampNowField
		}
		indexes: {primary: {attribute: "id"}, request: {attributes: ["request_id", "created_at"]}}
	}
	filter: {
		struct: {request_id: {ident: "requestID", goType: "uint64"}, event_type: {goType: "string"}}
		byValue: ["request_id", "event_type"]
	}
	features: {labels: false, flags: false}
	envoy: {omit: true}
	store: {ident: "city311AuditEvent", api: lookups: [{fields: ["id"]}]}
}

requestAttachment: {
	model: {
		ident:            "compose_city311_request_attachment"
		omitGetterSetter: true
		attributes: {
			id: schema.IdField
			request_id: {ident: "requestID", goType: "uint64", sortable: true, dal: {type: "ID"}}
			filename: {goType: "string", dal: {type: "Text", length: 120}}
			media_type: {goType: "string", dal: {type: "Text", length: 128}}
			size: {goType: "uint64", dal: {type: "Number", meta: {"rdbms:type": "bigint"}}}
			content: {goType: "[]byte", dal: {type: "Blob"}}
			created_at: schema.SortableTimestampNowField
		}
		indexes: {primary: {attribute: "id"}, request: {attributes: ["request_id", "created_at"]}}
	}
	filter: {
		struct: {request_id: {ident: "requestID", goType: "uint64"}}
		byValue: ["request_id"]
	}
	features: {labels: false, flags: false}
	envoy: {omit: true}
	store: {ident: "city311RequestAttachment", api: lookups: [{fields: ["id"]}]}
}

publicHistoryItem: {
	model: {
		ident:            "compose_city311_public_history"
		omitGetterSetter: true
		attributes: {
			id: schema.IdField
			request_id: {ident: "requestID", goType: "uint64", sortable: true, dal: {type: "ID"}}
			action: {goType: "string", dal: {type: "Text", length: 64}}
			responsible_department: {goType: "types.DepartmentCode", dal: {type: "Text", length: 64}}
			occurred_at: schema.SortableTimestampNowField
		}
		indexes: {primary: {attribute: "id"}, request: {attributes: ["request_id", "occurred_at"]}}
	}
	filter: {
		struct: {request_id: {ident: "requestID", goType: "uint64"}}
		byValue: ["request_id"]
	}
	features: {labels: false, flags: false}
	envoy: {omit: true}
	store: {ident: "city311PublicHistoryItem", api: lookups: [{fields: ["id"]}]}
}

actorProfile: {
	model: {
		ident:            "compose_city311_actor_profile"
		omitGetterSetter: true
		attributes: {
			id: schema.IdField
			application_roles: {goType: "types.City311ApplicationRoleSet", dal: {type: "JSON"}}
			department: {goType: "types.DepartmentCode", sortable: true, dal: {type: "Text", length: 64}}
			districts: {goType: "types.City311DistrictCodeSet", dal: {type: "JSON"}}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNowField
		}
		indexes: {primary: {attribute: "id"}, department: {attribute: "department"}}
	}
	filter: {struct: {department: {goType: "string"}}, byValue: ["department"]}
	features: {labels: false, flags: false}
	envoy: {omit: true}
	store: {ident: "city311ActorProfile", api: lookups: [{fields: ["id"]}]}
}
