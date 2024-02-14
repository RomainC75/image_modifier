package trades

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

var conn *websocket.Conn

const (
	subscribeId   = 1
	unSubscribeId = 1
)

func getConnection() (*websocket.Conn, error) {
	if conn != nil {
		return conn, nil
	}
	url := url.URL{Scheme: "wss", Host: "stream.binance.com:443", Path: "/ws"}
	log.Printf("connecting to : %s\n", url)

	c, resp, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Printf("handshake failed with code : %s \n", resp.StatusCode)
		log.Fatal("dial: ", err)
	}

	return c, nil
}

func SubscribeAndListen(topics []string) error {
	conn, err := getConnection()
	if err != nil {
		log.Fatal("Failed to connect %s", err.Error())
	}

	// binance regular ping
	conn.SetPongHandler(func(appData string) error {
		fmt.Println("Received pong:", appData)
		pingFrame := []byte{1, 2, 3, 4, 5}
		err := conn.WriteMessage(websocket.PingMessage, pingFrame)
		if err != nil {
			fmt.Println(err)
		}
		return nil
	})

	tradeTopics := make([]string, 0, len(topics))
	for _, t := range topics {
		tradeTopics = append(tradeTopics, t+"@aggTrade")
	}
	log.Println("Listen to trades for :", tradeTopics)

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
		log.Fatal("Failed to subscribe to topics" + err.Error())
	}

	defer conn.Close()
	// defer unsubscirbeOnClose(tradeTopics)

	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return err
		}

		trade := Ticker{}

		err = json.Unmarshal(payload, &trade)
		if err != nil {
			fmt.Println(err)
			return err
		}
		log.Println(trade.Symbol, trade.Price, trade.Quantity)
	}

}
