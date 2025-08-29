package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type KafkaSender interface {
	SendMessage(ctx context.Context, topic string, key, value []byte) error
}

func SendToKafka(log *slog.Logger, producer KafkaSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.SendToKafka"
		log = log.With(slog.String("fn", fn))

		var req struct {
			Topic   string `json:"topic"`
			Key     string `json:"key"`
			Message string `json:"message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		err := producer.SendMessage(r.Context(), req.Topic, []byte(req.Key), []byte(req.Message))
		if err != nil {
			log.Error("failed to send to kafka", slog.Any("error", err))
			http.Error(w, "failed to send", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
