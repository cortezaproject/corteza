package types

var (
	NodeStatusPending      = "pending"
	NodeStatusPairRequest  = "pair_requested"
	NodeStatusPairComplete = "paired"
)

type (
	Node struct {
		ID       uint64 `json:"nodeID,string"`
		SharedID uint64 `json:"sharedNodeID,string"`
		Name     string `json:"name,string"`
		Domain   string `json:"domain,string"`
		Status   string `json:"status,string"`
		Token    string `json:"-"`
		NodeURI  string `json:"nodeURI,string"`
	}
)
