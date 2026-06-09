package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"time"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/sql"
)

type (
	KnowledgeBase struct {
		ID        uint64 `json:"knowledgeBaseID,string"`
		TenantID  uint64 `json:"tenantID,string,omitempty"`
		ProjectID uint64 `json:"projectID,string,omitempty"`
		Handle    string `json:"handle"`

		Title       string `json:"title"`
		Description string `json:"description,omitempty"`

		Context *KnowledgeBaseContext `json:"context,omitempty"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
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

	KnowledgeBaseFilter struct {
		KnowledgeBaseID []uint64     `json:"knowledgeBaseID"`
		Handle          string       `json:"handle"`
		Query           string       `json:"query"`
		Deleted         filter.State `json:"deleted"`

		Check func(*KnowledgeBase) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)

func (m *KnowledgeBaseContext) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m KnowledgeBaseContext) Value() (driver.Value, error) { return json.Marshal(m) }

func (ll KnowledgeBaseIDList) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, id := range ll {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte('"')
		buf.WriteString(strconv.FormatUint(id, 10))
		buf.WriteByte('"')
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
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
