package appstore

// ConnectionSummary is returned by GET /v1/connections — lightweight, no config blobs.
type ConnectionSummary struct {
	ID          string             `json:"id"`
	Handle      string             `json:"handle"`
	Status      string             `json:"status"`
	Short       string             `json:"short"`
	Description string             `json:"description,omitempty"`
	Icon        string             `json:"icon,omitempty"`
	Tags        []string           `json:"tags,omitempty"`
	Operations  []OperationSummary `json:"operations"`
	Resources   []ResourceSummary  `json:"resources,omitempty"`
	Triggers    []TriggerSummary   `json:"triggers,omitempty"`
}

// Connection is returned by GET /v1/connections/{id}/config — full, with config blobs.
type Connection struct {
	ID         string      `json:"id"`
	Handle     string      `json:"handle"`
	Status     string      `json:"status"`
	Meta       Meta        `json:"meta"`
	Service    interface{} `json:"service"`
	Operations interface{} `json:"operations"`
	Resources  interface{} `json:"resources"`
}

type Meta struct {
	Short       string   `json:"short"`
	Description string   `json:"description,omitempty"`
	Icon        string   `json:"icon,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// OAuthApp is a provider blueprint from GET /v1/oauth-apps/{handle} — the public
// OAuth endpoints keyed by the connector's oauthApp name. Secrets (client id/
// secret) are never here; they stay in the instance settings.
type OAuthApp struct {
	Handle             string            `json:"handle"`
	AuthURL            string            `json:"authURL"`
	TokenURL           string            `json:"tokenURL"`
	IdentityURL        string            `json:"identityURL"`
	IdentityEmailField string            `json:"identityEmailField"`
	AuthParams         map[string]string `json:"authParams"`
	PKCE               bool              `json:"pkce"`

	// Setup guidance for the admin form — public, provider-specific.
	ConsoleURL string   `json:"consoleURL,omitempty"`
	SetupSteps []string `json:"setupSteps,omitempty"`
	DocsURL    string   `json:"docsURL,omitempty"`
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
