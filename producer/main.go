package main

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type CustomBalancer struct {
	// use the inbuilt round robin balancer
	baseBalancer kafka.RoundRobin
}

func (b *CustomBalancer) Balance(msg kafka.Message, partitions ...int) (partition int) {
	// only error logs are sent to first partition
	if string(msg.Key) == "error" {
		return 0
	}
	// everything else get round robined to other partition
	return b.baseBalancer.Balance(msg, partitions[1:]...)
}
func main() {
	// 1. create topic "logs"
	ctx := context.Background()
	topic := "logs"
	conn, err := kafka.DialContext(ctx, "tcp", "localhost:9092")
	if err != nil {
		panic(err)
	}
	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     3,
		ReplicationFactor: 1,
	})
	if err != nil {
		panic(err)
	}
	// 2 create writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    topic,
		Balancer: &CustomBalancer{},
	})
	err = writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte("error"),
			Value: []byte("some error happend in the application"),
		},
		kafka.Message{
			Key:   []byte("warn"),
			Value: []byte("this is a warning"),
		},
	)
	if err != nil {
		panic(err)
	}
	log.Println("written messages successfully")
	//
}
