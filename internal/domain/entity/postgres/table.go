package postgres

type Table struct {
	Name string `db:"table_name"`
}

type Column struct {
	Name string `db:"column_name"`
	Type string `db:"data_type"`
}
