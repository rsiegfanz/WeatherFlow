package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/rsiegfanz/WeatherFlow/pkg/payload"
)

const url = "wss://thingsboard.bda-itnovum.com/api/ws"

type Client struct {
	Token string
}

func NewClient(token string) (*Client, error) {
	return &Client{
		Token: token,
	}, nil
}

func (c *Client) Close() {}

func (c *Client) Connect() error {

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("connection error: %v", err)
	}
	defer conn.Close()

	initPayload := payload.PrepareInitPayload(c.Token)

	jsonData, err := json.Marshal(initPayload)
	if err != nil {
		return fmt.Errorf("JSON conversion error: %v", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		return fmt.Errorf("send error: %v", err)
	}

	log.Println("Initial message sent")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})
	go c.readData(conn, done)

	for {
		select {
		case <-done:
			return nil
		case <-interrupt:
			log.Println("Closing connection")
			err := conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return fmt.Errorf("error closing: %v", err)
			}
			return nil
		}
	}
}

func (c *Client) readData(conn *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}
		log.Println("New message:")
		log.Println(string(message))
	}
}
