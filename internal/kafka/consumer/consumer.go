package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"h-project/internal/entity"
	"h-project/internal/kafka/constants"
	"log/slog"
	"os"
)

type DataConsumer interface {
	ConsumeData()
}

type kafkaConsumer struct {
	consumer  *kafka.Consumer
	isRunning bool
	logger    *slog.Logger
}

func NewKafkaConsumer(topic string, logger *slog.Logger) (*kafkaConsumer, error) {
	logger.Info("Starting kafka consumer")
	logger.Info("Topic: ", topic)
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv(""),
		"group.id":          constants.KafkaGroupId,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		logger.Error("Kafka consumer error when starting ", "port", os.Getenv("KAFKA_BOOTSTRAP_SERVERS"))
		return nil, err
	}

	logger.Info("Kafka consumer started on server ", "port", os.Getenv("KAFKA_BOOTSTRAP_SERVERS"))

	c.SubscribeTopics([]string{topic}, nil)

	return &kafkaConsumer{
		consumer: c,
		logger:   logger,
	}, nil
}

func (c *kafkaConsumer) Start() {
	c.logger.Info("kafka transport started")
	c.isRunning = true
	fmt.Println(c.consumer.IsClosed(), "IsClosed")

	go c.readMessageLoop()
}

func (c *kafkaConsumer) Stop() error {
	c.logger.Info("kafka transport stopped")
	c.isRunning = false
	return nil
}

func (c *kafkaConsumer) Commit() error {
	offsets, err := c.consumer.Commit()
	if err != nil {
		c.logger.Error("Error committing offsets", "error", err)
		return err
	}

	c.logger.Info("Offsets committed", "offsets", offsets)
	return nil
}

func (c *kafkaConsumer) readMessageLoop() {
	c.logger.Info("readMessageLoop launched")
	for c.isRunning {
		fmt.Println(c.isRunning, "here")
		msg, err := c.consumer.ReadMessage(-1)
		fmt.Println(msg, "msg")
		if err != nil {
			c.logger.Error("kafka consume error " + err.Error())
			continue
		}
		fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		var data entity.Company
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serialization error %s", err)
			//Что отправить в прометеус кол-во неудачных сообщений
			//разграничивать ошибки
			// circuit breaker -
			continue
		}
		fmt.Println(msg.Value, "value")

		fmt.Println(data, "data")

		// Commit the offsets
		if err := c.Commit(); err != nil {
			c.logger.Error("Error committing offsets", "error", err)
			// Handle the error or continue without committing
		}

		//try := 0
		//var (
		//	distance float64
		//)
		//for {
		//	distance, err = c.calcService.CalculateDistance(data)
		//	if err == nil {
		//
		//	}
		//	try = try + 1
		//	if try > 10 {
		//		logrus.Warnf("Error calculating distance, try %d, error %s", try, err)
		//	}
		//
		//}
	}
}
