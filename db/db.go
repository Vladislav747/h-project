package db

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
	"h-project/config"
	"strconv"
)

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
