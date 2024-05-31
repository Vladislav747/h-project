package consumer

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
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
	logger.Debug("Starting kafka consumer")
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"group.id":          "myGroup",
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
	//c.readMessageLoop()
}

func (c *kafkaConsumer) Stop() {
	c.logger.Info("kafka transport stopped")
	c.isRunning = false
}

//func (c *kafkaConsumer) readMessageLoop() {
//	for c.isRunning {
//		msg, err := c.consumer.ReadMessage(-1)
//		if err != nil {
//			c.logger.Error("kafka consume error " + err.Error())
//			continue
//		}
//		fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
//		var data entity.Company
//		if err := json.Unmarshal(msg.Value, &data); err != nil {
//			logrus.Errorf("JSON serialization error %s", err)
//			//Что отправить в прометеус кол-во неудачных сообщений
//			//разграничивать ошибки
//			// circuit breaker -
//			continue
//		}
//		//try := 0
//		var (
//			distance float64
//		)
//		//for {
//		//	distance, err = c.calcService.CalculateDistance(data)
//		//	if err == nil {
//		//
//		//	}
//		//	try = try + 1
//		//	if try > 10 {
//		//		logrus.Warnf("Error calculating distance, try %d, error %s", try, err)
//		//	}
//		//
//		//}
//
//		if err != nil {
//			logrus.Errorf("Calculation error %s", err)
//			continue
//		}
//		fmt.Println(distance, "distance")
//	}
//}
