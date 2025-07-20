package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"food-order-backend/internal/shared/config"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var globalHub *Hub
var once sync.Once

// GetHub returns the singleton websocket hub
func GetHub() *Hub {
	once.Do(func() {
		globalHub = NewHub()
		go globalHub.Run()
	})
	return globalHub
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *Client) writePump() {
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.conn.Close()
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	hub.register <- client
	go client.writePump()
	go client.readPump()
}

// Broadcast sends a message to all clients
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

// SubscribeAndBroadcastFromStream subscribes to Redis stream and broadcasts messages to all clients
func (h *Hub) SubscribeAndBroadcastFromStream(ctx context.Context) {
	// Create a stream consumer for WebSocket broadcasting
	consumer := config.NewStreamConsumer(
		config.OrderEventsStream,
		config.ConsumerGroup,
		config.ConsumerName,
		func(message redis.XMessage) error {
			// Convert message to JSON for broadcasting
			eventData := map[string]interface{}{
				"order_id":   message.Values["order_id"],
				"user_id":    message.Values["user_id"],
				"event_type": message.Values["event_type"],
				"timestamp":  message.Values["timestamp"],
				"data":       message.Values,
			}

			// Remove redundant fields from data
			if data, ok := eventData["data"].(map[string]interface{}); ok {
				delete(data, "order_id")
				delete(data, "user_id")
				delete(data, "event_type")
				delete(data, "timestamp")
			}

			// Marshal to JSON and broadcast
			if jsonData, err := json.Marshal(eventData); err == nil {
				h.Broadcast(jsonData)
			}

			return nil
		},
	)

	// Start the consumer
	if err := consumer.Start(ctx); err != nil {
		log.Printf("Failed to start stream consumer: %v", err)
	}
}
