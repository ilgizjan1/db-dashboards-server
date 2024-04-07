package mysqlRepo

import (
	"context"
	"database/sql"
	"db-dashboards/internal/domain/entity/mysqlEntity"
	"fmt"
	"strings"
)

type Repo struct {
	DB *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{
		DB: db,
	}
}

func (r *Repo) GetAllTables(ctx context.Context) ([]*mysqlEntity.Table, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT * FROM information_schema.tables WHERE table_schema = 'koa'") // todo: return only public tables?
	if err != nil {
		return nil, err
	}

	columnNames, _ := rows.Columns()

	var tables []*mysqlEntity.Table

	for rows.Next() {
		columns := make([]string, len(columnNames))
		columnPointers := make([]any, len(columnNames))

		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		var table mysqlEntity.Table

		for i, columnName := range columnNames {
			if columnName == "TABLE_NAME" {
				table.Name = columns[i]
				break
			}
		}

		tables = append(tables, &table)
	}

	return tables, nil
}

func (r *Repo) GetColumnsFromTable(ctx context.Context, tableName string) ([]*mysqlEntity.Column, error) {

	rows, err := r.DB.QueryContext(ctx,
		fmt.Sprintf("SELECT * FROM information_schema.columns WHERE table_name   = '%v' order by ordinal_position", tableName))

	if err != nil {
		return nil, err
	}

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	var columns []*mysqlEntity.Column

	for rows.Next() {
		cols := make([]string, len(columnNames))
		columnPointers := make([]any, len(columnNames))

		for i, _ := range cols {
			columnPointers[i] = &cols[i]
		}

		rows.Scan(columnPointers...)

		var column mysqlEntity.Column

		for i, columnName := range columnNames {
			if columnName == "COLUMN_NAME" {
				column.Name = cols[i]
			}

			if columnName == "DATA_TYPE" {
				column.Type = strings.ToLower(columnTypes[i].DatabaseTypeName())
			}

			if column.Name != "" && column.Type != "" {
				break
			}
		}

		columns = append(columns, &column)
	}

	return columns, nil
}
