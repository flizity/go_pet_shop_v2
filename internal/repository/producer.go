package repository

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer sarama.AsyncProducer
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V2_8_0_0

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	kp := &KafkaProducer{producer: producer}
	go kp.handleResponses()
	return kp, nil
}

func (kp *KafkaProducer) handleResponses() {
	for {
		select {
		case msg := <-kp.producer.Successes():
			log.Printf("Message sent to topic %s partition %d offset %d", msg.Topic, msg.Partition, msg.Offset)
		case err := <-kp.producer.Errors():
			log.Printf("Failed to send message: %v", err)
		}
	}
}

func (kp *KafkaProducer) SendMessage(ctx context.Context, topic string, key, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	select {
	case kp.producer.Input() <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (kp *KafkaProducer) Close() error {
	return kp.producer.Close()
}
