package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type Payload struct {
	Sensor   string  `json:"sensor"`
	Index    uint    `json:"index"`
	Temp     float32 `json:"temp"`
	Humidity float32 `json:"humidity"`
	Time     string  `json:"time"`
}

func (r *repository) Consumer1(ctx context.Context) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{r.setting.Broker},
		Topic:   r.setting.Topic1,
		GroupID: r.setting.GroupID,
		Dialer: &kafka.Dialer{
			Timeout:       10 * time.Second,
			SASLMechanism: r.sasl,
		},
		// Logger: log.New(os.Stdout, "kafka-reader: ", 0),
	})
	defer reader.Close()

	fmt.Println("📥 Consumer started, waiting for messages...")

	msgCh := make(chan kafka.Message, 100)
	errCh := make(chan error, 1)
	wg := sync.WaitGroup{}

	// goroutine baca message
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					fmt.Println("🛑 Consumer stopped gracefully")
					return
				}
				errCh <- err
				return
			}
			msgCh <- m
		}
	}()

	// goroutine agregasi window
	wg.Add(1)
	go func() {
		defer wg.Done()
		window := time.NewTicker(5 * time.Minute)
		defer window.Stop()

		var values []float32
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-errCh:
				fmt.Println("❌ error:", err)
				return
			case m := <-msgCh:
				fmt.Printf("✅ Offset[%d] message received:%s\n", m.Offset, string(m.Value))
				var data Payload
				if err := json.Unmarshal(m.Value, &data); err != nil {
					continue
				}
				values = append(values, data.Temp)

			case <-window.C:
				if len(values) > 0 {
					avg := average(values)
					fmt.Printf("🧮 5-min window average: %.2f\n", avg)
					// direct ke consumer-2 (topic stream-go-agg)
					if err := r.Producer2(ctx, fmt.Sprintf("%.2f", avg)); err != nil {
						fmt.Printf("❌ failed to send to topic-2: %v\n", err)
					}
					values = nil
				}
			}
		}
	}()

	wg.Wait()
	return nil
}

func average(values []float32) float32 {
	if len(values) == 0 {
		return 0
	}
	var sum float32
	for _, v := range values {
		sum += v
	}
	return sum / float32(len(values))
}

func (r *repository) Consumer2(ctx context.Context) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{r.setting.Broker},
		Topic:   r.setting.Topic2,
		GroupID: r.setting.GroupID,
		Dialer: &kafka.Dialer{
			Timeout:       10 * time.Second,
			SASLMechanism: r.sasl,
		},
		// Logger: log.New(os.Stdout, "kafka-reader: ", 0),
	})
	defer reader.Close()

	fmt.Println("📥 Consumer-2 started, waiting for messages...")

	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}
		fmt.Printf("✅ message received: %s\n", string(m.Value))
	}
}
