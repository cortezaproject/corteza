package types

type (
	Permission struct {
		Resource  string `json:"resource"`
		Operation string `json:"operation"`
	}
)
