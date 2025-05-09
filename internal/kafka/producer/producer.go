package producer

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"h-project/internal/entity"
	"log/slog"
	"os"
)

type DataProducer interface {
	ProduceData(entity.Company) error
}

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
	logger   *slog.Logger
}

func NewKafkaProducer(topic string, logger *slog.Logger) (DataProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
	})
	if err != nil {
		logger.Error("Error creating kafka producer", "error", err)
		return nil, err
	}
	logger.Info("Kafka producer started on server ", "port", os.Getenv("KAFKA_BOOTSTRAP_SERVERS"))
	// Delivery report handler for produced messages
	//Start another goroutine to check if we have delivered the data
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logger.Error(fmt.Sprintf("Delivery failed: %v", ev.TopicPartition))
				} else {
					logger.Info(fmt.Sprintf("Delivered message to %v\n", ev.TopicPartition))
				}
			}
		}
	}()
	return &KafkaProducer{
		producer: p,
		topic:    topic,
		logger:   logger,
	}, nil
}

func (p *KafkaProducer) ProduceData(data entity.Company) error {
	b, err := json.Marshal(data)
	if err != nil {
		p.logger.Error("Error marshalling data", "error", err)
		return err
	}
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          b,
	}, nil)

}
