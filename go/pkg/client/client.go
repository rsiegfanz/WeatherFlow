package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rsiegfanz/WeatherFlow/pkg/payload"
)

const wsURL = "wss://thingsboard.bda-itnovum.com/api/ws"

type Client struct {
	token     string
	conn      *websocket.Conn
	msgLogger *log.Logger

	mu          sync.RWMutex
	entityNames map[string]string // entityId -> displayName
}

func New(token string, msgLogger *log.Logger) *Client {
	return &Client{
		token:       token,
		msgLogger:   msgLogger,
		entityNames: make(map[string]string),
	}
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
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println("Connection closed normally")
			} else {
				log.Printf("Read error: %v", err)
			}
			return
		}

		enriched := c.enrichMessage(raw)
		c.msgLogger.Println(string(enriched))
	}
}

type wsMessage struct {
	CmdID  int             `json:"cmdId"`
	Data   *entityDataPage `json:"data"`
	Update json.RawMessage `json:"update"`
}

type entityDataPage struct {
	Data []entityData `json:"data"`
}

type entityData struct {
	EntityID entityID                          `json:"entityId"`
	Latest   map[string]map[string]tsValue     `json:"latest"`
}

type entityID struct {
	EntityType string `json:"entityType"`
	ID         string `json:"id"`
}

type tsValue struct {
	Ts    int64  `json:"ts"`
	Value string `json:"value"`
}

func (c *Client) enrichMessage(raw []byte) []byte {
	var msg wsMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return raw
	}

	// Learn entity names from initial data responses
	if msg.Data != nil {
		c.learnEntityNames(msg.Data.Data)
	}

	// Enrich update messages with entity names
	if msg.Update != nil {
		var updates []entityData
		if err := json.Unmarshal(msg.Update, &updates); err == nil {
			c.learnEntityNames(updates)
			return c.buildEnrichedJSON(raw, updates)
		}
	}

	return raw
}

func (c *Client) learnEntityNames(entities []entityData) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, e := range entities {
		// From ATTRIBUTE.displayName
		if attrs, ok := e.Latest["ATTRIBUTE"]; ok {
			if dn, ok := attrs["displayName"]; ok && dn.Value != "" {
				c.entityNames[e.EntityID.ID] = dn.Value
			}
		}
		// From ENTITY_FIELD.label as fallback
		if _, exists := c.entityNames[e.EntityID.ID]; !exists {
			if fields, ok := e.Latest["ENTITY_FIELD"]; ok {
				if label, ok := fields["label"]; ok && label.Value != "" {
					c.entityNames[e.EntityID.ID] = label.Value
				}
			}
		}
	}
}

func (c *Client) buildEnrichedJSON(original []byte, updates []entityData) []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Parse as generic map to preserve all fields
	var full map[string]json.RawMessage
	if err := json.Unmarshal(original, &full); err != nil {
		return original
	}

	// Build entityNames lookup for this message
	names := make(map[string]string)
	for _, u := range updates {
		if name, ok := c.entityNames[u.EntityID.ID]; ok {
			names[u.EntityID.ID] = name
		}
	}

	if len(names) > 0 {
		namesJSON, _ := json.Marshal(names)
		full["_entityNames"] = namesJSON
	}

	enriched, err := json.Marshal(full)
	if err != nil {
		return original
	}
	return enriched
}
