package federation

import (
	"fmt"
	"io"

	"github.com/cortezaproject/corteza/server/pkg/options"
)

type (
	EncoderAdapterActivityStreams struct{}
)

// Build an activity streams format from default internal Corteza
// payload, including the author, activitystreams metadata and paging
// custom metadata
func (a EncoderAdapterActivityStreams) BuildStructure(w io.Writer, o options.FederationOpt, p interface{}) (interface{}, error) {
	var (
		next, prev *listResponsePagingActivityStreams
		items      []listResponseItemActivityStreams
	)

	payload := p.(ListStructurePayload)

	if payload.Filter.Paging.NextPage != nil {
		next = &listResponsePagingActivityStreams{
			Type: "Link",
			Name: "Next page",
			Href: fmt.Sprintf("https://%s/nodes/%d/modules/exposed/?pageCursor=%s", o.Host, payload.NodeID, payload.Filter.Paging.NextPage.Encode()),
		}
	}

	if payload.Filter.Paging.PrevPage != nil {
		prev = &listResponsePagingActivityStreams{
			Type: "Link",
			Name: "Previous page",
			Href: fmt.Sprintf("https://%s/nodes/%d/modules/exposed/?pageCursor=%s", o.Host, payload.NodeID, payload.Filter.Paging.PrevPage.Encode()),
		}
	}

	// loop through items and format them
	for _, v := range *payload.Set {
		item := listResponseItemActivityStreams{
			Context:          "https://www.w3.org/ns/activitystreams",
			Type:             "Module",
			Summary:          fmt.Sprintf("Structure for module %s on node %d", v.Name, v.NodeID),
			Url:              fmt.Sprintf("https://%s/nodes/%d/modules/%d", o.Host, v.NodeID, v.ID),
			Name:             v.Name,
			Handle:           v.Handle,
			Node:             v.NodeID,
			FederationModule: v.ID,
			ComposeModule:    v.ComposeModuleID,
			ComposeNamespace: v.ComposeNamespaceID,

			Attribution: []listResponseItemAttribution{
				{
					Context: "https://www.w3.org/ns/activitystreams",
					Type:    "User",
					Id:      fmt.Sprintf("https://%s/system/users/%d", o.Host, v.CreatedBy),
				},
			},

			CreatedAt: v.CreatedAt,
			CreatedBy: v.CreatedBy,

			UpdatedAt: v.UpdatedAt,
			UpdatedBy: v.UpdatedBy,

			DeletedAt: v.DeletedAt,
			DeletedBy: v.DeletedBy,

			Fields: v.Fields,
		}

		items = append(items, item)
	}

	return listModuleResponseActivityStreams{
		Context:      "https://www.w3.org/ns/activitystreams",
		ItemsPerPage: payload.Filter.Limit,
		Items:        items,
		Next:         next,
		Prev:         prev,
	}, nil
}

// Build an activity streams format from default internal Corteza
// payload, including the author, activitystreams metadata and paging
// custom metadata
func (a EncoderAdapterActivityStreams) BuildData(w io.Writer, o options.FederationOpt, p interface{}) (interface{}, error) {
	var (
		next, prev *listResponsePagingActivityStreams
		items      []listResponseItemActivityStreams
	)

	payload := p.(ListDataPayload)

	if payload.Filter.Paging.NextPage != nil {
		next = &listResponsePagingActivityStreams{
			Type: "Link",
			Name: "Next page",
			Href: fmt.Sprintf("https://%s/nodes/%d/modules/%d/records/social/?pageCursor=%s", o.Host, payload.NodeID, payload.ModuleID, payload.Filter.Paging.NextPage.Encode()),
		}
	}

	if payload.Filter.Paging.PrevPage != nil {
		prev = &listResponsePagingActivityStreams{
			Type: "Link",
			Name: "Previous page",
			Href: fmt.Sprintf("https://%s/nodes/%d/modules/%d/records/social/?pageCursor=%s", o.Host, payload.NodeID, payload.ModuleID, payload.Filter.Paging.PrevPage.Encode()),
		}
	}

	// loop through items and format them
	for _, v := range *payload.Set {
		item := listResponseItemActivityStreams{

			Context:          "https://www.w3.org/ns/activitystreams",
			Type:             "Record",
			Summary:          fmt.Sprintf("Data for module %d on node %d", v.ModuleID, payload.NodeID),
			Url:              fmt.Sprintf("https://%s/nodes/%d/modules/%d/records/social/", o.Host, payload.NodeID, payload.ModuleID),
			Node:             payload.NodeID,
			FederationModule: v.ModuleID,
			ComposeModule:    v.ModuleID,
			ComposeNamespace: v.NamespaceID,
			Attribution: []listResponseItemAttribution{
				{
					Context: "https://www.w3.org/ns/activitystreams",
					Type:    "User",
					Id:      fmt.Sprintf("https://%s/system/users/%d", o.Host, v.CreatedBy),
				},
			},

			CreatedAt: v.CreatedAt,
			CreatedBy: v.CreatedBy,

			UpdatedAt: v.UpdatedAt,
			UpdatedBy: v.UpdatedBy,

			DeletedAt: v.DeletedAt,
			DeletedBy: v.DeletedBy,
			Values:    v.Values,
		}

		items = append(items, item)
	}

	return listModuleResponseActivityStreams{
		Context:      "https://www.w3.org/ns/activitystreams",
		ItemsPerPage: payload.Filter.Limit,
		Items:        items,
		Next:         next,
		Prev:         prev,
	}, nil
}
