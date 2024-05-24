package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"h-project/internal/entity"
	"log/slog"
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
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	return &kafkaConsumer{
		consumer: c,
		logger:   logger,
	}, nil
}

func (c *kafkaConsumer) Start() {
	c.logger.Info("kafka transport started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *kafkaConsumer) Stop() {
	c.logger.Info("kafka transport stopped")
	c.isRunning = false
}

func (c *kafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
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
		//try := 0
		var (
			distance float64
		)
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

		if err != nil {
			logrus.Errorf("Calculation error %s", err)
			continue
		}
		fmt.Println(distance, "distance")
	}
}
