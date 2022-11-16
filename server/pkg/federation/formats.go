package federation

import (
	"time"

	ct "github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/federation/types"
)

type (
	listResponsePagingActivityStreams struct {
		Type string `json:"type"`
		Name string `json:"name"`
		Href string `json:"href"`
	}

	listResponseItemActivityStreams struct {
		Context string `json:"@context"`
		Type    string `json:"type"`
		Summary string `json:"summary"`
		Name    string `json:"name,omitempty"`
		Handle  string `json:"handle,omitempty"`
		Url     string `json:"url"`

		Node             uint64 `json:"node,string"`
		FederationModule uint64 `json:"federationModule,string"`
		ComposeModule    uint64 `json:"composeModule,string"`
		ComposeNamespace uint64 `json:"composeNamespace,string"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty" `

		Attribution []listResponseItemAttribution `json:"attributedTo"`

		Fields types.ModuleFieldSet `json:"fields,omitempty"`
		Values ct.RecordValueSet    `json:"values,omitempty"`
	}

	listResponseItemAttribution struct {
		Context string `json:"@context"`
		Id      string `json:"id"`
		Type    string `json:"type"`
	}

	listModuleResponseActivityStreams struct {
		Context      string                             `json:"@context"`
		ItemsPerPage uint                               `json:"itemsPerPage"`
		Next         *listResponsePagingActivityStreams `json:"next,omitempty"`
		Prev         *listResponsePagingActivityStreams `json:"prev,omitempty"`
		Items        interface{}                        `json:"items"`
	}

	listModuleResponseCortezaInternal struct {
		Filter *types.ExposedModuleFilter `json:"filter"`
		Set    *types.ExposedModuleSet    `json:"set"`
	}

	listRecordResponseCortezaInternal struct {
		Filter *ct.RecordFilter `json:"filter"`
		Set    *ct.RecordSet    `json:"set"`
	}

	ListStructurePayload struct {
		NodeID uint64
		Filter *types.ExposedModuleFilter `json:"filter"`
		Set    *types.ExposedModuleSet    `json:"set"`
	}

	ListDataPayload struct {
		NodeID   uint64
		ModuleID uint64
		Filter   *ct.RecordFilter `json:"filter"`
		Set      *ct.RecordSet    `json:"set"`
	}
)
