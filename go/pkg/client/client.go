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

const wsURL = "wss://thingsboard.bda-itnovum.com/api/ws"

type Client struct {
	token string
	conn  *websocket.Conn
}

func New(token string) *Client {
	return &Client{token: token}
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		c.conn.Close()
	}
}

func (c *Client) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	c.conn = conn

	if err := c.sendInitPayload(); err != nil {
		return err
	}

	log.Println("Connected and subscribed")

	done := make(chan struct{})
	go c.readMessages(done)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case <-done:
		return nil
	case <-interrupt:
		log.Println("Interrupted, closing connection")
		return nil
	}
}

func (c *Client) sendInitPayload() error {
	initPayload := payload.PrepareInitPayload(c.token)

	jsonData, err := json.Marshal(initPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		return fmt.Errorf("failed to send payload: %w", err)
	}

	return nil
}

func (c *Client) readMessages(done chan struct{}) {
	defer close(done)
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println("Connection closed normally")
			} else {
				log.Printf("Read error: %v", err)
			}
			return
		}
		log.Println(string(message))
	}
}
