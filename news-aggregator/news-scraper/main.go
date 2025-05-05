package main

import (
	"context"
	"fmt"
	"time"
	"sync"
	"github.com/segmentio/kafka-go"
)

var kafkaWriter = kafka.NewWriter(kafka.WriterConfig{
	Brokers: []string{"localhost:9092"},
	Topic:   "news-topic",
	Balancer: &kafka.LeastBytes{},
})

func scrapeSource(name string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Giả lập nội dung bài viết
	news := fmt.Sprintf(`{"source": "%s", "title": "Bài viết từ %s", "content": "Nội dung mẫu", "publishedAt": "%s"}`,
		name, name, time.Now().Format(time.RFC3339))

	// Gửi vào Kafka
	err := kafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Value: []byte(news),
		},
	)

	if err != nil {
		fmt.Printf("❌ Lỗi gửi Kafka: %v\n", err)
	} else {
		fmt.Printf("✅ Gửi tin từ %s thành công!\n", name)
	}
}

func main() {
	defer kafkaWriter.Close()

	var wg sync.WaitGroup
	sources := []string{"vnexpress.net", "tuoitre.vn", "thanhnien.vn"}

	for _, source := range sources {
		wg.Add(1)
		go scrapeSource(source, &wg)
	}

	wg.Wait()
	fmt.Println("✅ Hoàn tất thu thập tin!")
}
