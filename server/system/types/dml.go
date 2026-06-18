package types

type (
	DmlConnection struct {
		ID              uint64   `json:"connectionID,string"`
		DalConnectionID uint64   `json:"dalConnectionID,string"`
		Handle          string   `json:"handle"`
		Type            string   `json:"type"`
		Label           string   `json:"label"`
		Driver          string   `json:"driver"`
		Capabilities    []string `json:"capabilities"`
		ModelIdent      string   `json:"modelIdent,omitempty"`
	}

	DmlConnectionFilter struct {
		ConnectionID []string `json:"connectionID"`
		Handle       string   `json:"handle"`
		Type         string   `json:"type"`
	}

	DmlModel struct {
		ConnectionID uint64          `json:"connectionID,string"`
		Ident        string          `json:"ident"`
		Label        string          `json:"label"`
		ResourceType string          `json:"resourceType"`
		Attributes   []*DmlAttribute `json:"attributes"`
	}

	DmlAttribute struct {
		Ident      string `json:"ident"`
		Label      string `json:"label"`
		PrimaryKey bool   `json:"primaryKey"`
		Sortable   bool   `json:"sortable"`
		Filterable bool   `json:"filterable"`
		Type       string `json:"type"`
		Store      string `json:"store"`
	}
)
