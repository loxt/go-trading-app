package trades

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

var (
	HOST string
	PORT string
)

func LoadHostAndPort(host string, port string) {
	HOST = host
	PORT = port
}

//var conn *kafka.Conn
//
//func getConnection() *kafka.Conn {
//	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
//	if err != nil {
//		log.Fatal("failed to dial leader:", err)
//	}
//}

func Publish(message kafka.Message, topic string) error {
	messages := []kafka.Message{
		message,
	}

	w := &kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		Balancer:               &kafka.LeastBytes{},
	}

	defer w.Close()

	err := w.WriteMessages(context.Background(), messages...)

	if err != nil {
		log.Println("error writing messages to kafka: ", err)
		return err
	}

	log.Println("published message to kafka successfully on topic: ", topic)

	return nil
}
