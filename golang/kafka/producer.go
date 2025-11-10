package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func (r *repository) Producer1(ctx context.Context, msg string) error {
	w := kafka.Writer{
		Addr:  kafka.TCP(r.setting.Broker),
		Topic: r.setting.Topic1,
		Transport: &kafka.Transport{
			SASL: r.sasl, // SASL PLAIN
		},
	}
	defer w.Close()

	if err := w.WriteMessages(ctx,
		kafka.Message{
			Value: []byte(msg),
		},
	); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return nil
}

func (r *repository) Producer2(ctx context.Context, msg string) error {
	w := kafka.Writer{
		Addr:  kafka.TCP(r.setting.Broker),
		Topic: r.setting.Topic2,
		Transport: &kafka.Transport{
			SASL: r.sasl, // SASL PLAIN
		},
	}
	defer w.Close()

	if err := w.WriteMessages(ctx,
		kafka.Message{
			Value: []byte(msg),
		},
	); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return nil
}
