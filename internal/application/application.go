package application

import (
	"context"
	"h-project/api"
	"h-project/version"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	Name    string
	Version string
	Port    string

	ctx    context.Context
	cancel context.CancelFunc
	exit   chan bool

	httpServer *http.Server
	logger     *slog.Logger
}

func NewApplication() *Application {
	ctx, cancel := context.WithCancel(context.Background())

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := &Application{
		logger: logger,
		cancel: cancel,
		exit:   make(chan bool),
		ctx:    ctx,
	}

	// Create a new instance of http.ServeMux
	mux := http.NewServeMux()
	// Register your handlers
	mux.HandleFunc("/", api.HomeHandler)
	mux.HandleFunc("/status", api.StatusHandler)

	app.httpServer = &http.Server{
		Handler:      mux, // Use the newly created mux with your routes
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return app
}

func (app *Application) Run() {
	if app.Port == "" {
		app.Port = "80"
	}
	app.httpServer.Addr = ":" + app.Port

	app.logger.Info(
		"application started",
		"name", app.Name,
		"version", app.Version,
		"port", app.Port,
		"commit", version.Commit,
		"buildTime", version.BuildTime,
	)

	err := app.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		app.logger.Error(
			"error listening and serving",
			"error", err,
		)
	}

	app.cancel()
	<-app.exit
}

func (app *Application) gracefulShutdown() {
	// SIGTERM for docker container default signal
	signalCtx, cancel := signal.NotifyContext(app.ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// wait for parent or signal context to cancel
	<-signalCtx.Done()
	app.logger.Info("shutting down http server...")

	// make a new context for the shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.httpServer.Shutdown(shutdownCtx); err != nil {
		app.logger.Error(
			"error shutting down http server",
			"error", err,
		)
	}
	// Запись в канал который блокирует метод Run
	app.exit <- true

}
