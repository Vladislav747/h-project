package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	dbSqlx "h-project/db"
	"h-project/internal/application"
	"h-project/internal/file"
	"h-project/internal/file/storage/minio"
	consumerLib "h-project/internal/kafka/consumer"
	"h-project/internal/kafka/producer"
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

	port := os.Getenv("APPLICATION_PORT")

	app := application.NewApplication(logger)

	app.Name = name
	app.Version = versionName
	app.Port = port

	fileStorage, err := minio.NewStorage(os.Getenv("MINIO_ENDPOINT"), os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"), logger)
	if err != nil {
		logger.Error(err.Error())
	}
	fileService, err := file.NewService(fileStorage, logger)
	if err != nil {
		logger.Error(err.Error())
	}

	producerClient, err := producer.NewKafkaProducer(os.Getenv("KAFKA_TOPIC"), logger)
	if err != nil {
		logger.Error("Error launching Kafka Producer" + err.Error())
	}
	fmt.Print(producerClient)

	consumerClient, err := consumerLib.NewKafkaConsumer(os.Getenv("KAFKA_TOPIC"), logger)
	if err != nil {
		logger.Error("Error launching Kafka Consumer" + err.Error())
	}
	fmt.Print(consumerClient)

	handler := application.NewHTTPHandler(name, versionName, db, fileService, logger)
	app.RegisterHTTPHandler(handler)

	app.AddCloser(func() error {
		return db.Close(logger)
	})

	logger.Info("App Starting")

	app.Run()
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
