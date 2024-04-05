package main

import (
	"context"
	"database/sql"
	"db-dashboards/internal/repository/postgres"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	logger := logrus.New()
	ctx := context.Background()

	connStr := "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		logger.Fatalf("cannot open database connection with connection string: %v, err: %v", connStr, err)
	}

	db := sqlx.NewDb(conn, "postgres")

	repo := postgres.New(db)

	tables, err := repo.GetAllTables(ctx)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	for _, table := range tables {
		logger.Infof("table: %v", table.Name)

		columns, err := repo.GetColumnsFromTable(ctx, table.Name)
		if err != nil {
			logger.Fatalf(err.Error())
		}

		for _, col := range columns {
			logger.Infof("column %v of type %v", col.Name, col.Type)
		}

	}

}
