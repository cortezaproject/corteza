package appstore

// ConnectionSummary is returned by GET /v1/connections — lightweight, no config blobs.
type ConnectionSummary struct {
	ID          string             `json:"id"`
	Handle      string             `json:"handle"`
	Status      string             `json:"status"`
	Short       string             `json:"short"`
	Description string             `json:"description,omitempty"`
	Operations  []OperationSummary `json:"operations"`
	Resources   []ResourceSummary  `json:"resources,omitempty"`
	Triggers    []TriggerSummary   `json:"triggers,omitempty"`
}

// Connection is returned by GET /v1/connections/{id}/config — full, with config blobs.
type Connection struct {
	ID         string          `json:"id"`
	Handle     string          `json:"handle"`
	Status     string          `json:"status"`
	Meta       Meta            `json:"meta"`
	Service    interface{}     `json:"service"`
	Operations interface{}     `json:"operations"`
	Resources  interface{}     `json:"resources"`
}

type Meta struct {
	Short       string   `json:"short"`
	Description string   `json:"description,omitempty"`
	Icon        string   `json:"icon,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type OperationSummary struct {
	Handle      string `json:"handle"`
	Short       string `json:"short"`
	Description string `json:"description,omitempty"`
}

type ResourceSummary struct {
	Handle string `json:"handle"`
	Short  string `json:"short"`
}

type TriggerSummary struct {
	Event       string `json:"event"`
	Resource    string `json:"resource"`
	Description string `json:"description,omitempty"`
}
