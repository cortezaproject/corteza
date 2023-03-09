package types

type (
	Route struct {
		ID       uint64
		Endpoint string
		Method   string
		Meta     RouteMeta
	}

	RouteMeta struct {
		Debug bool
		Async bool
	}

	RouteFilter struct {
		ID      uint64            `json:"filterID,string"`
		Route   uint64            `json:"routeID,string"`
		Weight  uint64            `json:"weight,string"`
		Ref     string            `json:"ref,omitempty"`
		Kind    string            `json:"kind,omitempty"`
		Enabled bool              `json:"enabled,omitempty"`
		Params  RouteFilterParams `json:"params"`
	}

	RouteFilterParams map[string]interface{}
)
