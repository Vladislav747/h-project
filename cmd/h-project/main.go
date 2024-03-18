package main

import (
	"fmt"
	"h-project/internal/application"
	"h-project/version"
	"os"
)

func main() {

	name := version.APIName
	versionName := version.APIVersion

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
