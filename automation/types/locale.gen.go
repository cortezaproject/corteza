package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// Definitions file that controls how this file is generated:
// - automation.workflow.yaml

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"strconv"
)

type (
	LocaleKey struct {
		Name          string
		Resource      string
		Path          string
		CustomHandler string
	}
)

// Types and stuff
const (
	WorkflowResourceTranslationType = "automation:workflow"
)

var (
	LocaleKeyWorkflowName = LocaleKey{
		Name:     "name",
		Resource: WorkflowResourceTranslationType,
		Path:     "name",
	}
	LocaleKeyWorkflowDescription = LocaleKey{
		Name:     "description",
		Resource: WorkflowResourceTranslationType,
		Path:     "description",
	}
)

// ResourceTranslation returns string representation of Locale resource for Workflow by calling WorkflowResourceTranslation fn
//
// Locale resource is in the automation:workflow/... format
//
// This function is auto-generated
func (r Workflow) ResourceTranslation() string {
	return WorkflowResourceTranslation(r.ID)
}

// WorkflowResourceTranslation returns string representation of Locale resource for Workflow
//
// Locale resource is in the automation:workflow/... format
//
// This function is auto-generated
func WorkflowResourceTranslation(id uint64) string {
	cpts := []interface{}{WorkflowResourceTranslationType}
	cpts = append(cpts, strconv.FormatUint(id, 10))

	return fmt.Sprintf(WorkflowResourceTranslationTpl(), cpts...)
}

// @todo template
func WorkflowResourceTranslationTpl() string {
	return "%s/%s"
}

func (r *Workflow) DecodeTranslations(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation
	if aux = tt.FindByKey(LocaleKeyWorkflowName.Path); aux != nil {
		r.Meta.Name = aux.Msg
	} else {
		r.Meta.Name = LocaleKeyWorkflowName.Path
	}
	if aux = tt.FindByKey(LocaleKeyWorkflowDescription.Path); aux != nil {
		r.Meta.Description = aux.Msg
	} else {
		r.Meta.Description = LocaleKeyWorkflowDescription.Path
	}
}

func (r *Workflow) EncodeTranslations() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{
		{
			Resource: r.ResourceTranslation(),
			Key:      LocaleKeyWorkflowName.Path,
			Msg:      firstOkString(r.Meta.Name, LocaleKeyWorkflowName.Path),
		},
		{
			Resource: r.ResourceTranslation(),
			Key:      LocaleKeyWorkflowDescription.Path,
			Msg:      firstOkString(r.Meta.Description, LocaleKeyWorkflowDescription.Path),
		},
	}

	return out
}

func firstOkString(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}
