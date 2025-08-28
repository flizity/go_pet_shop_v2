package usecase

import (
	"context"
)

type MessageProducer interface {
	SendMessage(ctx context.Context, topic string, key, value []byte) error
}

type KafkaUsecase struct {
	producer MessageProducer
}

func NewKafkaUsecase(producer MessageProducer) *KafkaUsecase {
	return &KafkaUsecase{producer: producer}
}

func (u *KafkaUsecase) SendBusinessEvent(ctx context.Context, topic, key, value string) error {
	// Здесь может быть любая бизнес-логика
	return u.producer.SendMessage(ctx, topic, []byte(key), []byte(value))
}
