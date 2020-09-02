package types

type (
	Node struct {
		ID   uint64 `json:"recordID,string"`
		Name string `json:"name,string"`
	}
)
