package main

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"
)

func main() {

	app := fiber.New()

	listenAddr := flag.String("listenAddr", ":3001", "The listen address of the API server")
	flag.Parse()

	//Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	fmt.Println("App Starting")

	app.Listen(*listenAddr)
}
