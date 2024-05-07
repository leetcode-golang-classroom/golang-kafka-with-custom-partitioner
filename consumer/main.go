package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/segmentio/kafka-go"
)

func main() {
	ctx := context.Background()
	topic := "logs"

	errReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     topic,
		Partition: 0, // read only error logs
	})
	msgReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     topic,
		Partition: 1,
	})
	go func() {
		for {
			msg, err := errReader.ReadMessage(ctx)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(msg.Key))
			fmt.Println(string(msg.Value))
		}
	}()
	go func() {
		for {
			msg, err := msgReader.ReadMessage(ctx)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(msg.Key))
			fmt.Println(string(msg.Value))
		}
	}()
	// wait for ctrl + c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	log.Printf("receive signal %v, closing connection and shuting down", os.Interrupt)
}
