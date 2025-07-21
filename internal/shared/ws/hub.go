package ws

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
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
			log.Println("[WebSocket] Client registered. Total clients:", len(h.clients))
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Println("[WebSocket] Client unregistered. Total clients:", len(h.clients))
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			log.Printf("[WebSocket] Broadcasting message to %d clients", len(h.clients))
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
					log.Println("[WebSocket] Client send buffer full, removed client. Total clients:", len(h.clients))
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
		log.Println("[WebSocket] Client connection closed (readPump)")
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("[WebSocket] Read error:", err)
			break
		}
	}
}

func (c *Client) writePump() {
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("[WebSocket] Write error:", err)
			break
		}
	}
	c.conn.Close()
	log.Println("[WebSocket] Client connection closed (writePump)")
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	log.Println("[WebSocket] New client connection established")
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	hub.register <- client
	go client.writePump()
	go client.readPump()
}

// Broadcast sends a message to all clients
func (h *Hub) Broadcast(message []byte) {
	log.Printf("[WebSocket] Queuing message for broadcast: %s", string(message))
	h.broadcast <- message
}

// SubscribeAndBroadcastFromKafkaGo subscribes to Kafka topic using kafka-go and broadcasts messages to all clients
func (h *Hub) SubscribeAndBroadcastFromKafkaGo(ctx context.Context) {
	log.Println("[WebSocket] Starting kafka-go subscription for broadcasting...")

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if len(brokers) == 0 || brokers[0] == "" {
		brokers = []string{"localhost:9092"}
	}
	topic := os.Getenv("KAFKA_ORDER_EVENTS_TOPIC")
	if topic == "" {
		topic = "order_events"
	}
	groupID := os.Getenv("KAFKA_CONSUMER_GROUP")
	if groupID == "" {
		groupID = "websocket_group"
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	go func() {
		defer r.Close()
		for {
			m, err := r.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					log.Println("[WebSocket] Kafka-go context cancelled, stopping consumer.")
					return
				}
				log.Printf("[WebSocket] Kafka-go read error: %v", err)
				time.Sleep(time.Second)
				continue
			}
			log.Printf("[WebSocket] Received event from Kafka: %s", string(m.Value))
			h.Broadcast(m.Value)
		}
	}()

	log.Println("[WebSocket] kafka-go consumer started successfully.")
}
