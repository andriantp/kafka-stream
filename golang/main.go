package main

import (
	"context"
	"fmt"
	"kafka-stream/kafka"
	"log"
	"math/rand/v2"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// SASL_PLAINTEXT
	setting := kafka.Setting{
		Broker:  "localhost:9094",
		Topic1:  "stream-go",
		Topic2:  "stream-go-agg",
		GroupID: "group-1",

		Username: "test",
		Password: "test-secret",
	}

	if len(os.Args) < 2 {
		panic("Usage: go run . [producer|consumer]")
	}

	repo, err := kafka.NewKafka(setting)
	if err != nil {
		log.Fatalf("NewKafka:%v", err)
	}

	switch os.Args[1] {
	case "producer-1":
		fmt.Println("🚀 Starting producer stream... Press Ctrl+C to stop")

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for i := 0; ; i++ {
			temp := rand.Float64()*10 + 20     // suhu acak 20–30°C
			humidity := rand.Float64()*30 + 50 // kelembapan acak 50–80%
			msg := fmt.Sprintf(
				`{"sensor":"A1","index":%d,"temp":%.2f,"humidity":%.2f,"time":"%s"}`,
				i, temp, humidity, time.Now().Format(time.RFC3339),
			)

			if err := repo.Producer1(ctx, msg); err != nil {
				log.Printf("❌ Producer error: %v", err)
				continue
			}

			fmt.Printf("✅ sent: %s\n", msg)
			<-ticker.C
		}

	case "consumer-1":
		if err := repo.Consumer1(ctx); err != nil {
			log.Fatalf("Consumer:%v", err)
		}

	case "consumer-2":
		if err := repo.Consumer2(ctx); err != nil {
			log.Fatalf("Consumer2:%v", err)
		}

	default:
		panic("Unknown command: use producer or consumer")
	}
}
