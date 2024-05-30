package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
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
	dsn := os.Getenv("DSN")
	log.Println(dsn, "DSN")
	db, err := dbSqlx.NewDB(ctx, dsn)

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

	app.AddCloser(func() error {
		return db.Close(logger)
	})

	fmt.Println("App Starting")

	app.Run()
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
