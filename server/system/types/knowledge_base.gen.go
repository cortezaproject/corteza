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
	"strconv"
	"time"
)

type (
	KnowledgeBase struct {
		ID          uint64                `json:"knowledgeBaseID,string"`
		TenantID    uint64                `json:"tenantID,string,omitempty"`
		ProjectID   uint64                `json:"projectID,string,omitempty"`
		Handle      string                `json:"handle"`
		Title       string                `json:"title"`
		Description string                `json:"description,omitempty"`
		Context     *KnowledgeBaseContext `json:"context,omitempty"`
		CreatedAt   time.Time             `json:"createdAt,omitempty"`
		UpdatedAt   *time.Time            `json:"updatedAt,omitempty"`
		DeletedAt   *time.Time            `json:"deletedAt,omitempty"`
		CreatedBy   uint64                `json:"createdBy,string"`
		UpdatedBy   uint64                `json:"updatedBy,string,omitempty"`
		DeletedBy   uint64                `json:"deletedBy,string,omitempty"`
	}

	KnowledgeBaseContext struct {
		Namespaces []KnowledgeBaseNamespaceContext `json:"namespaces"`
	}

	KnowledgeBaseNamespaceContext struct {
		NamespaceID uint64              `json:"namespaceID,string"`
		ModuleIDs   KnowledgeBaseIDList `json:"moduleIDs"`
	}

	// KnowledgeBaseIDList is a []uint64 that serializes each element as a JSON string.
	KnowledgeBaseIDList []uint64
)

func (r KnowledgeBase) Clone() *KnowledgeBase {
	dup := r
	if r.Context != nil {
		dup.Context = r.Context.Clone()
	}

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

func (r KnowledgeBase) Diff(cmp *KnowledgeBase) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &KnowledgeBase{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "knowledgeBaseID", Old: []any{cmp.ID}, New: []any{r.ID}})
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

	if r.Title != cmp.Title {
		out = append(out, &revisions.Change{Key: "title", Old: []any{cmp.Title}, New: []any{r.Title}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if (r.Context == nil) != (cmp.Context == nil) {
		out = append(out, &revisions.Change{Key: "context", Old: []any{cmp.Context}, New: []any{r.Context}})
	} else if r.Context != nil {
		for _, c := range r.Context.Diff(cmp.Context) {
			c.Key = "context." + c.Key
			out = append(out, c)
		}
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

func (r KnowledgeBaseContext) Clone() *KnowledgeBaseContext {
	dup := r
	if r.Namespaces != nil {
		dup.Namespaces = make([]KnowledgeBaseNamespaceContext, len(r.Namespaces))
		for i := range r.Namespaces {
			dup.Namespaces[i] = *r.Namespaces[i].Clone()
		}
	}

	return &dup
}

func (r KnowledgeBaseContext) Diff(cmp *KnowledgeBaseContext) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &KnowledgeBaseContext{}
	}
	if !reflect.DeepEqual(r.Namespaces, cmp.Namespaces) {
		out = append(out, &revisions.Change{Key: "namespaces", Old: []any{cmp.Namespaces}, New: []any{r.Namespaces}})
	}

	return out
}

func (r *KnowledgeBaseContext) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r KnowledgeBaseContext) Value() (driver.Value, error) { return json.Marshal(r) }

func (r KnowledgeBaseNamespaceContext) Clone() *KnowledgeBaseNamespaceContext {
	dup := r
	return &dup
}

func (r KnowledgeBaseNamespaceContext) Diff(cmp *KnowledgeBaseNamespaceContext) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &KnowledgeBaseNamespaceContext{}
	}
	if r.NamespaceID != cmp.NamespaceID {
		out = append(out, &revisions.Change{Key: "namespaceID", Old: []any{cmp.NamespaceID}, New: []any{r.NamespaceID}})
	}

	if !reflect.DeepEqual(r.ModuleIDs, cmp.ModuleIDs) {
		out = append(out, &revisions.Change{Key: "moduleIDs", Old: []any{cmp.ModuleIDs}, New: []any{r.ModuleIDs}})
	}

	return out
}

func (r *KnowledgeBaseNamespaceContext) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r KnowledgeBaseNamespaceContext) Value() (driver.Value, error) { return json.Marshal(r) }

func (ll KnowledgeBaseIDList) MarshalJSON() ([]byte, error) {
	ss := make([]string, len(ll))
	for i, id := range ll {
		ss[i] = strconv.FormatUint(id, 10)
	}
	return json.Marshal(ss)
}

func (ll *KnowledgeBaseIDList) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*ll = make(KnowledgeBaseIDList, 0, len(raw))
	for _, r := range raw {
		s := string(r)
		if len(s) >= 2 && s[0] == '"' {
			s = s[1 : len(s)-1]
		}
		id, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		*ll = append(*ll, id)
	}
	return nil
}

func ParseKnowledgeBaseContext(ss []string) (p KnowledgeBaseContext, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
