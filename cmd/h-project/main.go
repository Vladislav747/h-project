package main

import (
	"context"
	"fmt"
	"h-project/config"
	dbSqlx "h-project/db"
	"h-project/internal/application"
	"h-project/version"
	"log"
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
	_, err := dbSqlx.NewDB(ctx, conf)

	if err != nil {
		log.Fatalln(err)
	}

	handler := application.NewHTTPHandler(name, versionName)

	port := os.Getenv("APPLICATION_PORT")

	app := application.NewApplication()

	app.Name = name
	app.Version = versionName
	app.Port = port

	app.RegisterHTTPHandler(handler)

	app.Run()

	fmt.Println("App Starting")
}
