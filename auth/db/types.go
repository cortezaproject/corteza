package db

type migration struct {
	Project        string `db:"project"`
	Filename       string `db:"filename"`
	StatementIndex int    `db:"statement_index"`
	Status         string `db:"status"`
}
