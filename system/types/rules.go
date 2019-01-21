package types

type Access int

const (
	Allow   Access = 1
	Deny    Access = 0
	Inherit Access = -1
)

type Rules struct {
	TeamID    uint64 `db:"rel_team"`
	Resource  string `db:"resource"`
	Operation string `db:"operation"`
	Value     Access `db:"value"`
}
