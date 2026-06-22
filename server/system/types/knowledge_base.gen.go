package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/crusttech/human/server/pkg/sql"
	"strconv"
	"time"
)

type KnowledgeBase struct {
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

func (r KnowledgeBase) Clone() *KnowledgeBase {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}

// KnowledgeBaseIDList is a []uint64 that serializes each element as a JSON string.
type KnowledgeBaseIDList []uint64

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

func (m *KnowledgeBaseContext) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m KnowledgeBaseContext) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseKnowledgeBaseContext(ss []string) (p KnowledgeBaseContext, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
