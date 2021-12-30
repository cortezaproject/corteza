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

	PreprocessorHandleResourceRemove preprocessor = "resourceRemove"
	PreprocessorHandleResourceLoad   preprocessor = "resourceLoad"
	PreprocessorHandleNamespaceLoad  preprocessor = "namespaceLoad"
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

func PreprocessorResourceRemove(resourceType, identifier string) (out preprocessorResourceRemove) {
	if identifier == "" {
		identifier = "*"
	}

	out = preprocessorResourceRemove{
		Resource:   resourceType,
		Identifier: identifier,
	}

	return
}

func (t preprocessorResourceRemove) Ref() preprocessor {
	return PreprocessorHandleResourceRemove
}

func (t preprocessorResourceRemove) Worker() []string {
	return []string{EnvoyWorkerName}
}

func (t preprocessorResourceRemove) Params() interface{} {
	return t
}

func PreprocessorResourceLoadByID(resourceType string, id uint64) (out preprocessorResourceLoad) {
	out = preprocessorResourceLoad{
		Resource: resourceType,
		ID:       id,
	}
	return
}

func PreprocessorResourceLoadByHandle(resourceType, handle string) (out preprocessorResourceLoad) {
	out = preprocessorResourceLoad{
		Resource: resourceType,
		Handle:   handle,
	}
	return
}

func PreprocessorResourceLoadByQuery(resourceType, query string) (out preprocessorResourceLoad) {
	out = preprocessorResourceLoad{
		Resource: resourceType,
		Query:    query,
	}
	return
}

func (t preprocessorResourceLoad) Ref() preprocessor {
	return PreprocessorHandleResourceLoad
}

func (t preprocessorResourceLoad) Worker() []string {
	return []string{EnvoyWorkerName}
}

func (t preprocessorResourceLoad) Params() interface{} {
	return t
}

func PreprocessorNamespaceLoadByID(id uint64) (out preprocessorNamespaceLoad) {
	out = preprocessorNamespaceLoad{
		ID: id,
	}

	return
}

func PreprocessorNamespaceLoadByHandle(handle string) (out preprocessorNamespaceLoad) {
	out = preprocessorNamespaceLoad{
		Handle: handle,
	}

	return
}

func (t preprocessorNamespaceLoad) Ref() preprocessor {
	return PreprocessorHandleNamespaceLoad
}

func (t preprocessorNamespaceLoad) Worker() []string {
	return []string{EnvoyWorkerName}
}

func (t preprocessorNamespaceLoad) Params() interface{} {
	return t
}
