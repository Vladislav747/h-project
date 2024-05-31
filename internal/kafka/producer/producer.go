package producer

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"h-project/internal/entity"
	"log/slog"
)

type DataProducer interface {
	ProduceData(entity.Company) error
}

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaProducer(topic string, logger *slog.Logger) (DataProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:19092",
	})
	if err != nil {
		logger.Error("Error creating kafka producer", "error", err)
		return nil, err
	}
	// Delivery report handler for produced messages
	//Start another goroutine to check if we have delivered the data
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	return &KafkaProducer{
		producer: p,
		topic:    topic,
	}, nil
}

func (p *KafkaProducer) ProduceData(data entity.Company) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          b,
	}, nil)

}
