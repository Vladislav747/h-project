package main

import (
	"context"
	"fmt"
	"h-project/config"
	dbSqlx "h-project/db"
	"h-project/internal/application"
	"h-project/version"
	"log"
	"log/slog"
	"os"
)

func main() {

	name := version.APIName
	versionName := version.APIVersion
	ctx := context.Background()
	conf := config.Config{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		TimeZone: os.Getenv("TIMEZONE"),
	}
	db, err := dbSqlx.NewDB(ctx, conf)

	if err != nil {
		log.Fatalln(err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	handler := application.NewHTTPHandler(name, versionName, db, logger)

	port := os.Getenv("APPLICATION_PORT")

	app := application.NewApplication(logger)

	app.Name = name
	app.Version = versionName
	app.Port = port

	app.RegisterHTTPHandler(handler)

	app.Run()

	fmt.Println("App Starting")
}
