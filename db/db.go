package db

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
	"h-project/config"
)

type DB struct {
	pg  *sqlx.DB
	ctx context.Context
}

func NewDB(ctx context.Context, conf config.Config) (*DB, error) {

	conn, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%d password=%s user=%s dbname=%s sslmode=disable timezone='%s'", conf.Host, conf.Port, conf.Password, conf.User, conf.DBName, conf.TimeZone))
	if err != nil {
		return nil, err
	}

	db := DB{
		pg:  conn,
		ctx: ctx,
	}
	return &db, nil
}
