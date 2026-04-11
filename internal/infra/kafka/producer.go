package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	kafkaContracts "github.com/jva44ka/ozon-simulator-go-products/api_internal/kafka"
	segkafka "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *segkafka.Writer
}

func NewProducer(brokers []string, topic string, writeTimeout time.Duration) *Producer {
	return &Producer{
		writer: &segkafka.Writer{
			Addr:         segkafka.TCP(brokers...),
			Topic:        topic,
			WriteTimeout: writeTimeout,
		},
	}
}

func (p *Producer) PublishProductChangedBatch(ctx context.Context, events []kafkaContracts.ProductEventMessage) error {
	messages := make([]segkafka.Message, 0, len(events))
	for _, message := range events {
		data, err := json.Marshal(message.Body)
		if err != nil {
			return fmt.Errorf("Producer.PublishProductChangedBatch: marshal sku=%d: %w", message.Key, err)
		}
		messages = append(messages, segkafka.Message{
			Key:   []byte(message.Key),
			Value: data,
		})
	}
	return p.writer.WriteMessages(ctx, messages...)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
