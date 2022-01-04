package gig

// @todo generate the file; this is a placeholder until qlang support is added

import (
	"context"
	"fmt"

	automationTypes "github.com/cortezaproject/corteza-server/automation/types"
	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/envoy/store"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/spf13/cast"
)

type (
	preprocessorResourceRemove struct {
		Resource   string `json:"resource"`
		Identifier string `json:"identifier"`
	}

	preprocessorResourceLoad struct {
		Resource string `json:"resource"`
		ID       uint64 `json:"id,omitempty"`
		Handle   string `json:"handle,omitempty"`
		Query    string `json:"query,omitempty"`
	}

	preprocessorNamespaceLoad struct {
		ID     uint64 `json:"id,omitempty"`
		Handle string `json:"handle,omitempty"`
	}
)

var (
	EnvoyWorkerName = "envoy"

	PreprocessorHandleResourceRemove = "resourceRemove"
	PreprocessorHandleResourceLoad   = "resourceLoad"
	PreprocessorHandleNamespaceLoad  = "namespaceLoad"
)

var (
	ApigwRouteResourceType  = systemTypes.ApigwRouteResourceType
	ApplicationResourceType = systemTypes.ApplicationResourceType
	ReportResourceType      = systemTypes.ReportResourceType
	RoleResourceType        = systemTypes.RoleResourceType
	TemplateResourceType    = systemTypes.TemplateResourceType
	UserResourceType        = systemTypes.UserResourceType

	ComposeChartResourceType       = composeTypes.ChartResourceType
	ComposeModuleResourceType      = composeTypes.ModuleResourceType
	ComposeModuleFieldResourceType = composeTypes.ModuleFieldResourceType
	ComposeNamespaceResourceType   = composeTypes.NamespaceResourceType
	ComposePageResourceType        = composeTypes.PageResourceType

	AutomationWorkflowResourceType = automationTypes.WorkflowResourceType

	SettingsResourceType            = resource.SettingsResourceType
	RbacResourceType                = resource.RbacResourceType
	ResourceTranslationResourceType = resource.ResourceTranslationType
)

// Utilities

func (w *workerEnvoy) preprocess(ctx context.Context, tasks ...Preprocessor) (err error) {
	for _, t := range tasks {
		switch tc := t.(type) {
		case preprocessorResourceRemove:
			err = w.resourceRemove(ctx, tc)
		case preprocessorResourceLoad:
			err = w.resourceLoad(ctx, tc)
		case preprocessorNamespaceLoad:
			err = w.namespaceLoad(ctx, tc)
		default:
			err = fmt.Errorf("unknown preprocessor: %s", w.Ref())
		}

		if err != nil {
			return
		}
	}

	return nil
}

func (w *workerEnvoy) filterComposeNamespace(base *store.DecodeFilter, defs preprocessorResourceLoad) (out *store.DecodeFilter) {
	out = base

	if defs.ID != 0 {
		out = out.ComposeNamespace(&composeTypes.NamespaceFilter{
			NamespaceID: []uint64{defs.ID},
		})
	}

	if defs.Handle != "" {
		out = out.ComposeNamespace(&composeTypes.NamespaceFilter{
			Slug: defs.Handle,
		})
	}

	return
}

func (w *workerEnvoy) filterComposeModule(base *store.DecodeFilter, defs preprocessorResourceLoad) (out *store.DecodeFilter) {
	out = base

	if defs.ID != 0 {
		out = out.ComposeModule(&composeTypes.ModuleFilter{
			ModuleID: []uint64{defs.ID},
		})
	}

	if defs.Handle != "" {
		out = out.ComposeModule(&composeTypes.ModuleFilter{
			Handle: defs.Handle,
		})
	}

	return
}

func (w *workerEnvoy) filterComposeChart(base *store.DecodeFilter, defs preprocessorResourceLoad) (out *store.DecodeFilter) {
	out = base

	if defs.ID != 0 {
		out = out.ComposeChart(&composeTypes.ChartFilter{
			ChartID: []uint64{defs.ID},
		})
	}

	if defs.Handle != "" {
		out = out.ComposeChart(&composeTypes.ChartFilter{
			Handle: defs.Handle,
		})
	}

	return
}

func (w *workerEnvoy) filterComposePage(base *store.DecodeFilter, defs preprocessorResourceLoad) (out *store.DecodeFilter) {
	out = base

	// if defs.id != 0 {
	// 	out = out.ComposePage(&composeTypes.PageFilter{
	// 		PageID: defs.id,
	// 	})
	// }

	if defs.Handle != "" {
		out = out.ComposePage(&composeTypes.PageFilter{
			Handle: defs.Handle,
		})
	}

	return
}

// Preprocessors

func PreprocessorResourceRemoveParams(params map[string]interface{}) (out preprocessorResourceRemove) {
	rt := cast.ToString(params["resourceType"])
	idf := cast.ToString(params["identifier"])

	return PreprocessorResourceRemove(rt, idf)
}

func PreprocessorResourceRemove(resource, identifier string) (out preprocessorResourceRemove) {
	out.Resource = resource
	out.Identifier = identifier

	if out.Identifier == "" {
		out.Identifier = "*"
	}

	return
}

func (t preprocessorResourceRemove) Ref() string {
	return PreprocessorHandleResourceRemove
}

func (t preprocessorResourceRemove) Worker() []string {
	return []string{EnvoyWorkerName}
}

func (t preprocessorResourceRemove) Params() map[string]interface{} {
	return map[string]interface{}{
		"resource":   t.Resource,
		"identifier": t.Identifier,
	}
}

func PreprocessorResourceLoadParams(params map[string]interface{}) (out preprocessorResourceLoad) {
	out.Resource = cast.ToString(params["resource"])

	if id, ok := params["id"]; ok {
		out.ID = cast.ToUint64(id)
	} else if handle, ok := params["handle"]; ok {
		out.Handle = cast.ToString(handle)
	} else if query, ok := params["query"]; ok {
		out.Query = cast.ToString(query)
	}

	return
}

func PreprocessorResourceLoadID(resource string, id uint64) (out preprocessorResourceLoad) {
	out.Resource = resource
	out.ID = id

	return
}

func PreprocessorResourceLoadHandle(resource, handle string) (out preprocessorResourceLoad) {
	out.Resource = resource
	out.Handle = handle

	return
}

func PreprocessorResourceLoadQuery(resource, query string) (out preprocessorResourceLoad) {
	out.Resource = resource
	out.Query = query

	return
}

func (t preprocessorResourceLoad) Ref() string {
	return PreprocessorHandleResourceLoad
}

func (t preprocessorResourceLoad) Worker() []string {
	return []string{EnvoyWorkerName}
}

func (t preprocessorResourceLoad) Params() map[string]interface{} {
	return map[string]interface{}{
		"resource": t.Resource,
		"id":       cast.ToString(t.ID),
		"handle":   t.Handle,
		"query":    t.Query,
	}
}

func PreprocessorNamespaceLoadParams(params map[string]interface{}) (out preprocessorNamespaceLoad) {
	if id, ok := params["id"]; ok {
		return PreprocessorNamespaceLoadID(cast.ToUint64(id))
	} else if handle, ok := params["handle"]; ok {
		return PreprocessorNamespaceLoadHandle(cast.ToString(handle))
	}
	return
}

func PreprocessorNamespaceLoadID(id uint64) (out preprocessorNamespaceLoad) {
	out.ID = id
	return
}

func PreprocessorNamespaceLoadHandle(handle string) (out preprocessorNamespaceLoad) {
	out.Handle = handle
	return
}

func (t preprocessorNamespaceLoad) Ref() string {
	return PreprocessorHandleNamespaceLoad
}

func (t preprocessorNamespaceLoad) Worker() []string {
	return []string{EnvoyWorkerName}
}

func (t preprocessorNamespaceLoad) Params() map[string]interface{} {
	return map[string]interface{}{
		"id":     cast.ToString(t.ID),
		"handle": t.Handle,
	}
}
