package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

var hub = &Hub{
	clients:    make(map[*Client]bool),
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

func RunHub() {
	for {
		select {
		case client := <-hub.register:
			hub.mu.Lock()
			hub.clients[client] = true
			hub.mu.Unlock()
			log.Println("client registered, count:", len(hub.clients))
		case client := <-hub.unregister:
			hub.mu.Lock()
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.Send)
				_ = client.Conn.Close()
			}
			hub.mu.Unlock()
			log.Println("client unregistered, count:", len(hub.clients))
		case message := <-hub.broadcast:
			hub.mu.Lock()
			for c := range hub.clients {
				select {
				case c.Send <- message:
				default:
					close(c.Send)
					delete(hub.clients, c)
				}
			}
			hub.mu.Unlock()
		}
	}
}

func BroadcastPriceUpdate(data []byte) {
	hub.broadcast <- data
}

func RegisterClient(c *Client) {
	hub.register <- c
}

func UnregisterClient(c *Client) {
	hub.unregister <- c
}
