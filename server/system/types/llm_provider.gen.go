package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/pkg/sql"
	"reflect"
	"time"
)

type (
	LlmProvider struct {
		ID           uint64            `json:"llmProviderID,string"`
		TenantID     uint64            `json:"tenantID,string,omitempty"`
		ProjectID    uint64            `json:"projectID,string,omitempty"`
		Handle       string            `json:"handle"`
		Status       string            `json:"status"`
		Provider     string            `json:"provider"`
		CredentialID uint64            `json:"credentialID,string"`
		Meta         LLMProviderMeta   `json:"meta"`
		Config       LLMProviderConfig `json:"config"`
		CreatedAt    time.Time         `json:"createdAt,omitempty"`
		UpdatedAt    *time.Time        `json:"updatedAt,omitempty"`
		DeletedAt    *time.Time        `json:"deletedAt,omitempty"`
		CreatedBy    uint64            `json:"createdBy,string"`
		UpdatedBy    uint64            `json:"updatedBy,string,omitempty"`
		DeletedBy    uint64            `json:"deletedBy,string,omitempty"`
	}

	LLMProviderMeta struct {
		Short       string `json:"short"`
		Description string `json:"description"`
	}

	LLMProviderConfig struct {
		PromptURL   string                  `json:"promptURL"`
		Model       string                  `json:"model"`
		Temperature *float64                `json:"temperature,omitempty"`
		Timeout     string                  `json:"timeout"`
		Guard       *LLMProviderGuardConfig `json:"guard,omitempty"`
	}

	LLMProviderGuardConfig struct {
		Enabled    bool               `json:"enabled"`
		Provider   string             `json:"provider"`
		Model      string             `json:"model"`
		Thresholds map[string]float64 `json:"thresholds,omitempty"`
	}
)

func (r LlmProvider) Clone() *LlmProvider {
	dup := r
	dup.Meta = *r.Meta.Clone()

	dup.Config = *r.Config.Clone()

	if r.UpdatedAt != nil {
		v := *r.UpdatedAt
		dup.UpdatedAt = &v
	}

	if r.DeletedAt != nil {
		v := *r.DeletedAt
		dup.DeletedAt = &v
	}

	return &dup
}

func (r LlmProvider) Diff(cmp *LlmProvider) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &LlmProvider{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "llmProviderID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.Handle != cmp.Handle {
		out = append(out, &revisions.Change{Key: "handle", Old: []any{cmp.Handle}, New: []any{r.Handle}})
	}

	if r.Status != cmp.Status {
		out = append(out, &revisions.Change{Key: "status", Old: []any{cmp.Status}, New: []any{r.Status}})
	}

	if r.Provider != cmp.Provider {
		out = append(out, &revisions.Change{Key: "provider", Old: []any{cmp.Provider}, New: []any{r.Provider}})
	}

	if r.CredentialID != cmp.CredentialID {
		out = append(out, &revisions.Change{Key: "credentialID", Old: []any{cmp.CredentialID}, New: []any{r.CredentialID}})
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
	}

	for _, c := range r.Config.Diff(&cmp.Config) {
		c.Key = "config." + c.Key
		out = append(out, c)
	}

	if !reflect.DeepEqual(r.CreatedAt, cmp.CreatedAt) {
		out = append(out, &revisions.Change{Key: "createdAt", Old: []any{cmp.CreatedAt}, New: []any{r.CreatedAt}})
	}

	if !reflect.DeepEqual(r.UpdatedAt, cmp.UpdatedAt) {
		out = append(out, &revisions.Change{Key: "updatedAt", Old: []any{cmp.UpdatedAt}, New: []any{r.UpdatedAt}})
	}

	if !reflect.DeepEqual(r.DeletedAt, cmp.DeletedAt) {
		out = append(out, &revisions.Change{Key: "deletedAt", Old: []any{cmp.DeletedAt}, New: []any{r.DeletedAt}})
	}

	if r.CreatedBy != cmp.CreatedBy {
		out = append(out, &revisions.Change{Key: "createdBy", Old: []any{cmp.CreatedBy}, New: []any{r.CreatedBy}})
	}

	if r.UpdatedBy != cmp.UpdatedBy {
		out = append(out, &revisions.Change{Key: "updatedBy", Old: []any{cmp.UpdatedBy}, New: []any{r.UpdatedBy}})
	}

	if r.DeletedBy != cmp.DeletedBy {
		out = append(out, &revisions.Change{Key: "deletedBy", Old: []any{cmp.DeletedBy}, New: []any{r.DeletedBy}})
	}

	return out
}

func (r LLMProviderMeta) Clone() *LLMProviderMeta {
	dup := r
	return &dup
}

func (r LLMProviderMeta) Diff(cmp *LLMProviderMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &LLMProviderMeta{}
	}
	if r.Short != cmp.Short {
		out = append(out, &revisions.Change{Key: "short", Old: []any{cmp.Short}, New: []any{r.Short}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	return out
}

func (r *LLMProviderMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r LLMProviderMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r LLMProviderConfig) Clone() *LLMProviderConfig {
	dup := r
	if r.Temperature != nil {
		v := *r.Temperature
		dup.Temperature = &v
	}

	if r.Guard != nil {
		dup.Guard = r.Guard.Clone()
	}

	return &dup
}

func (r LLMProviderConfig) Diff(cmp *LLMProviderConfig) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &LLMProviderConfig{}
	}
	if r.PromptURL != cmp.PromptURL {
		out = append(out, &revisions.Change{Key: "promptURL", Old: []any{cmp.PromptURL}, New: []any{r.PromptURL}})
	}

	if r.Model != cmp.Model {
		out = append(out, &revisions.Change{Key: "model", Old: []any{cmp.Model}, New: []any{r.Model}})
	}

	if !reflect.DeepEqual(r.Temperature, cmp.Temperature) {
		out = append(out, &revisions.Change{Key: "temperature", Old: []any{cmp.Temperature}, New: []any{r.Temperature}})
	}

	if r.Timeout != cmp.Timeout {
		out = append(out, &revisions.Change{Key: "timeout", Old: []any{cmp.Timeout}, New: []any{r.Timeout}})
	}

	if (r.Guard == nil) != (cmp.Guard == nil) {
		out = append(out, &revisions.Change{Key: "guard", Old: []any{cmp.Guard}, New: []any{r.Guard}})
	} else if r.Guard != nil {
		for _, c := range r.Guard.Diff(cmp.Guard) {
			c.Key = "guard." + c.Key
			out = append(out, c)
		}
	}

	return out
}

func (r *LLMProviderConfig) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r LLMProviderConfig) Value() (driver.Value, error) { return json.Marshal(r) }

func (r LLMProviderGuardConfig) Clone() *LLMProviderGuardConfig {
	dup := r
	if r.Thresholds != nil {
		dup.Thresholds = make(map[string]float64, len(r.Thresholds))
		for k, v := range r.Thresholds {
			dup.Thresholds[k] = v
		}
	}

	return &dup
}

func (r LLMProviderGuardConfig) Diff(cmp *LLMProviderGuardConfig) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &LLMProviderGuardConfig{}
	}
	if r.Enabled != cmp.Enabled {
		out = append(out, &revisions.Change{Key: "enabled", Old: []any{cmp.Enabled}, New: []any{r.Enabled}})
	}

	if r.Provider != cmp.Provider {
		out = append(out, &revisions.Change{Key: "provider", Old: []any{cmp.Provider}, New: []any{r.Provider}})
	}

	if r.Model != cmp.Model {
		out = append(out, &revisions.Change{Key: "model", Old: []any{cmp.Model}, New: []any{r.Model}})
	}

	if !reflect.DeepEqual(r.Thresholds, cmp.Thresholds) {
		out = append(out, &revisions.Change{Key: "thresholds", Old: []any{cmp.Thresholds}, New: []any{r.Thresholds}})
	}

	return out
}

func (r *LLMProviderGuardConfig) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r LLMProviderGuardConfig) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseLLMProviderMeta(ss []string) (p LLMProviderMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParseLLMProviderConfig(ss []string) (p LLMProviderConfig, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
