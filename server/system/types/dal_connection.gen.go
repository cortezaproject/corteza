package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/geolocation"
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/pkg/sql"
	"reflect"
	"time"
)

type (
	DalConnection struct {
		ID        uint64              `json:"connectionID,string"`
		Handle    string              `json:"handle"`
		Type      string              `json:"type"`
		Config    DalConnectionConfig `json:"config"`
		Meta      DalConnectionMeta   `json:"meta"`
		Issues    []dal.Issue         `json:"issues,omitempty"`
		Labels    map[string]string   `json:"labels,omitempty"`
		CreatedAt time.Time           `json:"createdAt,omitempty"`
		UpdatedAt *time.Time          `json:"updatedAt,omitempty"`
		DeletedAt *time.Time          `json:"deletedAt,omitempty"`
		CreatedBy uint64              `json:"createdBy,string"`
		UpdatedBy uint64              `json:"updatedBy,string,omitempty"`
		DeletedBy uint64              `json:"deletedBy,string,omitempty"`
	}

	DalConnectionConfig struct {
		DAL     *DalConnectionConfigDAL    `json:"dal,omitempty"`
		Privacy DalConnectionConfigPrivacy `json:"privacy"`
	}

	DalConnectionConfigPrivacy struct {
		SensitivityLevelID uint64 `json:"sensitivityLevelID,string,omitempty"`
	}

	DalConnectionConfigDAL struct {
		Type            string         `json:"type"`
		Params          map[string]any `json:"params"`
		ModelIdent      string         `json:"modelIdent"`
		ModelIdentCheck []string       `json:"modelIdentCheck"`
	}

	DalConnectionMeta struct {
		Name       string                      `json:"name"`
		Ownership  string                      `json:"ownership"`
		Location   geolocation.Full            `json:"location"`
		Properties DalConnectionMetaProperties `json:"properties"`
	}

	DalConnectionMetaProperties struct {
		DataAtRestEncryption    DalConnectionMetaProperty `json:"dataAtRestEncryption"`
		DataAtRestProtection    DalConnectionMetaProperty `json:"dataAtRestProtection"`
		DataAtTransitEncryption DalConnectionMetaProperty `json:"dataAtTransitEncryption"`
		DataRestoration         DalConnectionMetaProperty `json:"dataRestoration"`
	}

	DalConnectionMetaProperty struct {
		Enabled bool   `json:"enabled"`
		Notes   string `json:"notes"`
	}
)

func (r DalConnection) Clone() *DalConnection {
	dup := r
	dup.Config = *r.Config.Clone()

	dup.Meta = *r.Meta.Clone()

	if r.Issues != nil {
		dup.Issues = make([]dal.Issue, len(r.Issues))
		copy(dup.Issues, r.Issues)
	}

	if r.Labels != nil {
		dup.Labels = make(map[string]string, len(r.Labels))
		for k, v := range r.Labels {
			dup.Labels[k] = v
		}
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

func (r DalConnection) Diff(cmp *DalConnection) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DalConnection{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "connectionID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.Handle != cmp.Handle {
		out = append(out, &revisions.Change{Key: "handle", Old: []any{cmp.Handle}, New: []any{r.Handle}})
	}

	if r.Type != cmp.Type {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	for _, c := range r.Config.Diff(&cmp.Config) {
		c.Key = "config." + c.Key
		out = append(out, c)
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
	}

	if !reflect.DeepEqual(r.Issues, cmp.Issues) {
		out = append(out, &revisions.Change{Key: "issues", Old: []any{cmp.Issues}, New: []any{r.Issues}})
	}

	if !reflect.DeepEqual(r.Labels, cmp.Labels) {
		out = append(out, &revisions.Change{Key: "labels", Old: []any{cmp.Labels}, New: []any{r.Labels}})
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

func (r DalConnectionConfig) Clone() *DalConnectionConfig {
	dup := r
	if r.DAL != nil {
		dup.DAL = r.DAL.Clone()
	}

	dup.Privacy = *r.Privacy.Clone()

	return &dup
}

func (r DalConnectionConfig) Diff(cmp *DalConnectionConfig) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DalConnectionConfig{}
	}
	if (r.DAL == nil) != (cmp.DAL == nil) {
		out = append(out, &revisions.Change{Key: "dal", Old: []any{cmp.DAL}, New: []any{r.DAL}})
	} else if r.DAL != nil {
		for _, c := range r.DAL.Diff(cmp.DAL) {
			c.Key = "dal." + c.Key
			out = append(out, c)
		}
	}

	for _, c := range r.Privacy.Diff(&cmp.Privacy) {
		c.Key = "privacy." + c.Key
		out = append(out, c)
	}

	return out
}

func (r *DalConnectionConfig) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r DalConnectionConfig) Value() (driver.Value, error) { return json.Marshal(r) }

func (r DalConnectionConfigPrivacy) Clone() *DalConnectionConfigPrivacy {
	dup := r
	return &dup
}

func (r DalConnectionConfigPrivacy) Diff(cmp *DalConnectionConfigPrivacy) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DalConnectionConfigPrivacy{}
	}
	if r.SensitivityLevelID != cmp.SensitivityLevelID {
		out = append(out, &revisions.Change{Key: "sensitivityLevelID", Old: []any{cmp.SensitivityLevelID}, New: []any{r.SensitivityLevelID}})
	}

	return out
}

func (r *DalConnectionConfigPrivacy) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r DalConnectionConfigPrivacy) Value() (driver.Value, error) { return json.Marshal(r) }

func (r DalConnectionConfigDAL) Clone() *DalConnectionConfigDAL {
	dup := r
	if r.Params != nil {
		dup.Params = make(map[string]any, len(r.Params))
		for k, v := range r.Params {
			dup.Params[k] = v
		}
	}

	if r.ModelIdentCheck != nil {
		dup.ModelIdentCheck = make([]string, len(r.ModelIdentCheck))
		copy(dup.ModelIdentCheck, r.ModelIdentCheck)
	}

	return &dup
}

func (r DalConnectionConfigDAL) Diff(cmp *DalConnectionConfigDAL) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DalConnectionConfigDAL{}
	}
	if r.Type != cmp.Type {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	if !reflect.DeepEqual(r.Params, cmp.Params) {
		out = append(out, &revisions.Change{Key: "params", Old: []any{cmp.Params}, New: []any{r.Params}})
	}

	if r.ModelIdent != cmp.ModelIdent {
		out = append(out, &revisions.Change{Key: "modelIdent", Old: []any{cmp.ModelIdent}, New: []any{r.ModelIdent}})
	}

	if !reflect.DeepEqual(r.ModelIdentCheck, cmp.ModelIdentCheck) {
		out = append(out, &revisions.Change{Key: "modelIdentCheck", Old: []any{cmp.ModelIdentCheck}, New: []any{r.ModelIdentCheck}})
	}

	return out
}

func (r *DalConnectionConfigDAL) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r DalConnectionConfigDAL) Value() (driver.Value, error) { return json.Marshal(r) }

func (r DalConnectionMeta) Clone() *DalConnectionMeta {
	dup := r
	dup.Properties = *r.Properties.Clone()

	return &dup
}

func (r DalConnectionMeta) Diff(cmp *DalConnectionMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DalConnectionMeta{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Ownership != cmp.Ownership {
		out = append(out, &revisions.Change{Key: "ownership", Old: []any{cmp.Ownership}, New: []any{r.Ownership}})
	}

	if !reflect.DeepEqual(r.Location, cmp.Location) {
		out = append(out, &revisions.Change{Key: "location", Old: []any{cmp.Location}, New: []any{r.Location}})
	}

	for _, c := range r.Properties.Diff(&cmp.Properties) {
		c.Key = "properties." + c.Key
		out = append(out, c)
	}

	return out
}

func (r *DalConnectionMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r DalConnectionMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r DalConnectionMetaProperties) Clone() *DalConnectionMetaProperties {
	dup := r
	dup.DataAtRestEncryption = *r.DataAtRestEncryption.Clone()

	dup.DataAtRestProtection = *r.DataAtRestProtection.Clone()

	dup.DataAtTransitEncryption = *r.DataAtTransitEncryption.Clone()

	dup.DataRestoration = *r.DataRestoration.Clone()

	return &dup
}

func (r DalConnectionMetaProperties) Diff(cmp *DalConnectionMetaProperties) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DalConnectionMetaProperties{}
	}
	for _, c := range r.DataAtRestEncryption.Diff(&cmp.DataAtRestEncryption) {
		c.Key = "dataAtRestEncryption." + c.Key
		out = append(out, c)
	}

	for _, c := range r.DataAtRestProtection.Diff(&cmp.DataAtRestProtection) {
		c.Key = "dataAtRestProtection." + c.Key
		out = append(out, c)
	}

	for _, c := range r.DataAtTransitEncryption.Diff(&cmp.DataAtTransitEncryption) {
		c.Key = "dataAtTransitEncryption." + c.Key
		out = append(out, c)
	}

	for _, c := range r.DataRestoration.Diff(&cmp.DataRestoration) {
		c.Key = "dataRestoration." + c.Key
		out = append(out, c)
	}

	return out
}

func (r *DalConnectionMetaProperties) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r DalConnectionMetaProperties) Value() (driver.Value, error) { return json.Marshal(r) }

func (r DalConnectionMetaProperty) Clone() *DalConnectionMetaProperty {
	dup := r
	return &dup
}

func (r DalConnectionMetaProperty) Diff(cmp *DalConnectionMetaProperty) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DalConnectionMetaProperty{}
	}
	if r.Enabled != cmp.Enabled {
		out = append(out, &revisions.Change{Key: "enabled", Old: []any{cmp.Enabled}, New: []any{r.Enabled}})
	}

	if r.Notes != cmp.Notes {
		out = append(out, &revisions.Change{Key: "notes", Old: []any{cmp.Notes}, New: []any{r.Notes}})
	}

	return out
}

func (r *DalConnectionMetaProperty) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r DalConnectionMetaProperty) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseDalConnectionConfig(ss []string) (p DalConnectionConfig, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParseDalConnectionMeta(ss []string) (p DalConnectionMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
