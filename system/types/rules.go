package types

type Rules struct {
	TeamID    uint64 `db:"rel_team"`
	Resource  string `db:"resource"`
	Operation string `db:"operation"`
	Value     string `db:"value"`
}
