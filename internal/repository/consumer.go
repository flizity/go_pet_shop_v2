package repository

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
)

func StartConsumer(brokers []string, topic string, group string, ctx context.Context) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_8_0_0

	consumer, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer consumer.Close()

	handler := &consumerGroupHandler{}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		for {
			if err := consumer.Consume(ctx, []string{topic}, handler); err != nil {
				log.Printf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	select {
	case <-signals:
		log.Println("Consumer stopped by signal")
	case <-ctx.Done():
		log.Println("Consumer stopped by context")
	}
}

type consumerGroupHandler struct{}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)
		session.MarkMessage(msg, "")
	}
	return nil
}
