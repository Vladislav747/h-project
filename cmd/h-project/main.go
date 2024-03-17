package main

import (
	"fmt"
	"h-project/internal/application"
	"h-project/version"
	"os"
)

func main() {

	name := version.APIName
	version := version.APIVersion

	port := os.Getenv("APPLICATION_PORT")

	app := application.NewApplication()

	app.Name = name
	app.Version = version
	app.Port = port

	app.Run()

	fmt.Println("App Starting")
}
