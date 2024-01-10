package trades

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"log"
	"net/url"
	"strconv"
	"strings"
)

var conn *websocket.Conn

const (
	subscribeId   = 1
	unSubscribeId = 2
)

func getConnection() (*websocket.Conn, error) {
	if conn != nil {
		return conn, nil
	}
	u := url.URL{Scheme: "wss", Host: "stream.binance.com:443", Path: "/ws"}
	log.Printf("connecting to %s", u.String())
	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("dial:", err)
	}
	return c, nil
}

type RequestParams struct {
	Id     int      `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}

func unsubscribeOnClose(conn *websocket.Conn, tradeTopics []string) error {
	message := RequestParams{
		Id:     unSubscribeId,
		Method: "UNSUBSCRIBE",
		Params: tradeTopics,
	}

	b, err := json.Marshal(message)

	if err != nil {
		log.Fatal("Failed to JSON Encode trade topics")
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, b)

	return nil
}

func SubscribeAndListen(topics []string) error {
	conn, err := getConnection()
	if err != nil {
		log.Fatalf("Failed to get connection %s", err.Error())
		return err
	}

	conn.SetPongHandler(func(appData string) error {
		fmt.Println("received pong:", appData)
		pingFrame := []byte{1, 2, 3, 4, 5}

		err := conn.WriteMessage(websocket.PingMessage, pingFrame)

		if err != nil {
			log.Println("Failed to send pong message")
			return err
		}

		return nil
	})

	tradeTopics := make([]string, 0, len(topics))
	for _, topic := range topics {
		tradeTopics = append(tradeTopics, fmt.Sprintf("%s@aggTrade", topic))
	}
	log.Printf("Subscribing to topics: %s", tradeTopics)

	message := RequestParams{
		Id:     subscribeId,
		Method: "SUBSCRIBE",
		Params: tradeTopics,
	}

	b, err := json.Marshal(message)

	if err != nil {
		log.Fatal("Failed to JSON Encode trade topics")
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, b)

	if err != nil {
		log.Fatal("Failed to write message")
		return err
	}

	defer conn.Close()
	defer unsubscribeOnClose(conn, tradeTopics)

	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return err
		}

		trade := Ticker{}

		err = json.Unmarshal(payload, &trade)
		if err != nil {
			log.Println(err)
			return err
		}

		log.Println(trade.Symbol, trade.Price, trade.Quantity)
		go func() { // <=== here
			convertAndPublishToKafka(trade)
		}()
	}
}

// add this function
func convertAndPublishToKafka(t Ticker) {
	bytes, err := json.Marshal(t)
	if err != nil {
		log.Println("Error marshalling Ticker data", err.Error())
	}

	Publish(kafka.Message{
		Key:   []byte(t.Symbol + "-" + strconv.Itoa(int(t.Time))),
		Value: bytes,
	}, "trades-"+strings.ToLower(t.Symbol))
}
