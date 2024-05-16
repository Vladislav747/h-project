package db

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
	"h-project/config"
	"h-project/internal/entity"
	"strconv"
)

//go:embed sql/selectCompanies.sql
var selectCompanies string

//go:embed sql/insertCompany.sql
var insertCompany string

type DB struct {
	pg  *sqlx.DB
	ctx context.Context
}

func NewDB(ctx context.Context, conf config.Config) (*DB, error) {
	port, err := strconv.ParseInt(conf.Port, 10, 0)
	if err != nil {
		fmt.Println("Error converting string to integer:", err)
		return nil, err
	}
	conn, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%d password=%s user=%s dbname=%s sslmode=disable timezone='%s'", conf.Host, port, conf.Password, conf.User, conf.DBName, conf.TimeZone))
	if err != nil {
		return nil, err
	}

	db := DB{
		pg:  conn,
		ctx: ctx,
	}
	return &db, nil
}

func (db *DB) Exec(query string, args ...interface{}) error {
	_, err := db.pg.ExecContext(db.ctx, query, args...)
	return err
}
func (db *DB) NamedQuery(query string, args ...interface{}) (*sqlx.Rows, error) {
	return db.pg.NamedQueryContext(db.ctx, query, args)
}
func (db *DB) Select(dest interface{}, query string, args ...interface{}) error {
	return db.pg.SelectContext(db.ctx, dest, query, args...)
}
func (db *DB) Get(dest interface{}, query string, args ...interface{}) error {
	return db.pg.GetContext(db.ctx, dest, query, args...)
}
func (db *DB) Close() {
	_ = db.pg.Close()
}
func (db *DB) GetCompanies() (*[]entity.Company, error) {
	var companies []entity.Company
	err := db.Select(&companies, selectCompanies)
	return &companies, err
}

func (db *DB) AddCompany(company *entity.Company) error {
	var id int

	res, err := db.NamedQuery(insertCompany, company)
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
