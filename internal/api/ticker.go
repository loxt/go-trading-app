package api

import (
	"context"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/loxt/go-trading-app/internal/core"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
)

func GetAllTickers(c *fiber.Ctx) error {
	_ = c.Status(200).JSON(core.GetAllTickers())
	return nil
}

func ListenTicker(conn *websocket.Conn) {
	// get current ticker in conn params
	currTicker := conn.Params("ticker")
	log.Println("Current ticker: ", currTicker)

	if !core.IsTickerAllowed(currTicker) {
		conn.WriteMessage(websocket.CloseUnsupportedData, []byte("Ticker not allowed"))
		log.Println("Ticker not allowed: ", currTicker)
		return
	}

	topic := fmt.Sprintf("trades-%s", strings.ToLower(currTicker))
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{fmt.Sprintf("%s:%s", core.KAFKA_HOST, core.KAFKA_PORT)},
		Topic:   topic,
	})

	reader.SetOffset(-1)

	conn.SetCloseHandler(func(code int, text string) error {
		reader.Close()
		log.Printf("received connection closed request. Closing connection...")
		return nil
	})

	defer reader.Close()
	defer conn.Close()

	go func() {
		code, wsMessage, err := conn.NextReader()
		if err != nil {
			log.Println("error reading last message from ws connection. Exiting....")
			conn.Close()
			return
		}
		log.Printf("received close request from client. Closing connection... code: %d, message: %s", code, wsMessage)
	}()

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error: ", err)
			return
		}

		fmt.Println("Reading....", string(message.Value))

		conn.WriteMessage(websocket.TextMessage, message.Value)
	}

}
