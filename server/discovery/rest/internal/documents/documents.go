package documents

import (
	"github.com/cortezaproject/corteza-server/pkg/filter"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"time"
)

type (
	Response struct {
		Filter    Filter     `json:"filter"`
		Documents []Document `json:"documents"`
	}

	Document struct {
		ID     uint64      `json:"id,string"`
		URL    string      `json:"url,omitempty"`
		Source interface{} `json:"source"`
	}

	Filter struct {
		Limit    uint                 `json:"limit"`
		NextPage *filter.PagingCursor `json:"nextPage"`
	}

	docPartialUser struct {
		UserID   uint64 `json:"userID,string"`
		Email    string `json:"email,omitempty"`
		Name     string `json:"name,omitempty"`
		Username string `json:"username,omitempty"`
		Handle   string `json:"handle,omitempty"`
	}

	docPartialChange struct {
		At *time.Time      `json:"at,omitempty"`
		By *docPartialUser `json:"by,omitempty"`
	}

	docUser struct {
		ResourceType string            `json:"resourceType"`
		UserID       uint64            `json:"userID,string"`
		Email        string            `json:"email"`
		Name         string            `json:"name,omitempty"`
		Handle       string            `json:"handle,omitempty"`
		Suspended    *time.Time        `json:"suspendedAt,omitempty"`
		Url          string            `json:"url,omitempty"`
		Updated      *docPartialChange `json:"updated,omitempty"`
		Created      *docPartialChange `json:"created,omitempty"`
		Deleted      *docPartialChange `json:"deleted,omitempty"`
		Security     docSecurity       `json:"security"`
	}

	docComposeNamespace struct {
		ResourceType string                         `json:"resourceType"`
		NamespaceID  uint64                         `json:"namespaceID,string"`
		Name         string                         `json:"name,omitempty"`
		Handle       string                         `json:"handle,omitempty"`
		Url          string                         `json:"url,omitempty"`
		Enabled      bool                           `json:"enabled"`
		Meta         docPartialComposeNamespaceMeta `json:"meta"`

		Updated  *docPartialChange `json:"updated,omitempty"`
		Created  *docPartialChange `json:"created,omitempty"`
		Deleted  *docPartialChange `json:"deleted,omitempty"`
		Security docSecurity       `json:"security"`

		// Aggregation update
		Namespace docPartialComposeNamespace `json:"namespace"`
	}

	docPartialComposeNamespaceMeta struct {
		Subtitle    string `json:"subtitle,omitempty"`
		Description string `json:"description,omitempty"`
	}

	docPartialComposeNamespace struct {
		NamespaceID uint64 `json:"namespaceID,string"`
		Name        string `json:"name,omitempty"`
		Handle      string `json:"handle,omitempty"`
	}

	docComposeModule struct {
		ResourceType string                          `json:"resourceType"`
		ModuleID     uint64                          `json:"moduleID,string"`
		Name         string                          `json:"name,omitempty"`
		Handle       string                          `json:"handle,omitempty"`
		Url          string                          `json:"url,omitempty"`
		Labels       map[string]string               `json:"labels,omitempty"`
		Fields       []*docPartialComposeModuleField `json:"fields"`
		Updated      *docPartialChange               `json:"updated,omitempty"`
		Created      *docPartialChange               `json:"created,omitempty"`
		Deleted      *docPartialChange               `json:"deleted,omitempty"`
		Security     docSecurity                     `json:"security"`

		// Aggregation update
		Namespace docPartialComposeNamespace `json:"namespace"`
		Module    docPartialComposeModule    `json:"module"`
	}

	docPartialComposeModuleField struct {
		Name  string `json:"name,omitempty"`
		Label string `json:"label,omitempty"`
	}

	docPartialComposeModule struct {
		ModuleID uint64 `json:"moduleID,string"`
		Name     string `json:"name,omitempty"`
		Handle   string `json:"handle,omitempty"`
	}

	docComposeRecord struct {
		ResourceType string                   `json:"resourceType"`
		RecordID     uint64                   `json:"recordID,string"`
		Url          string                   `json:"url,omitempty"`
		Labels       map[string]string        `json:"labels,omitempty"`
		ValueLabels  map[string]string        `json:"valueLabels,omitempty"`
		Values       map[string][]interface{} `json:"values"`
		Updated      *docPartialChange        `json:"updated,omitempty"`
		Created      *docPartialChange        `json:"created,omitempty"`
		Deleted      *docPartialChange        `json:"deleted,omitempty"`
		Security     docSecurity              `json:"security"`

		// Aggregation update
		Namespace docPartialComposeNamespace `json:"namespace"`
		Module    docPartialComposeModule    `json:"module"`
	}

	docSecurity struct {
		// list of roles that are allowed to read the resource
		AllowedRoles []uint64 `json:"allowedRoles"`

		// list of roles that are disallowed to read the resource
		DeniedRoles []uint64 `json:"deniedRoles"`
	}
)

func makePartialChange(at *time.Time) *docPartialChange {
	if at == nil {
		return nil
	}

	// @todo handle unreadable (access-control) modules
	// @todo attach security info (allow/deny roles)
	return &docPartialChange{At: at}
}

func makeChange(at *time.Time, user *sysTypes.User) (out *docPartialChange) {
	out = &docPartialChange{}
	if at != nil {
		out.At = at
	}
	if user != nil {
		out.By = &docPartialUser{
			UserID:   user.ID,
			Email:    user.Email,
			Name:     user.Name,
			Username: user.Username,
			Handle:   user.Handle,
		}
	}

	return
}
