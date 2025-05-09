package db

import (
	"context"
	_ "embed"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
	"h-project/internal/entity"
	"log"
	"log/slog"
)

//go:embed sql/selectCompanies.sql
var selectCompanies string

//go:embed sql/insertCompany.sql
var insertCompany string

type DB struct {
	pg  *sqlx.DB
	ctx context.Context
}

func NewDB(ctx context.Context, dsn string) (*DB, error) {
	log.Println(dsn, "DSN 2")
	conn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db := DB{
		pg:  conn,
		ctx: ctx,
	}
	return &db, nil
}

func (db *DB) Close(logger *slog.Logger) error {
	logger.Info("Closing DB")
	return db.pg.Close()

}
func (db *DB) GetCompanies() (*[]entity.Company, error) {
	var companies []entity.Company
	err := db.pg.Get(&companies, selectCompanies)
	return &companies, err
}

func (db *DB) AddCompany(company *entity.Company) error {
	var id int

	res, err := db.pg.NamedQuery(insertCompany, company)
	if err != nil {
		return err
	}

	for res.Rows.Next() {
		if err = res.Rows.Scan(&id); err != nil {
			return err
		}
	}

	return nil
}
