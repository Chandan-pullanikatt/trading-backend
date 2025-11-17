package controllers

import (
	// "encoding/json" // <- REMOVED: This import was unused
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/chandan/trading-dashboard/mockdata"
	"github.com/chandan/trading-dashboard/models"
	"github.com/chandan/trading-dashboard/websocket" // This is your internal package
	"github.com/gin-gonic/gin"
	gw "github.com/gorilla/websocket" //RENAMED: Aliased to 'gw' to avoid conflict
)

// CHANGED: Now uses 'gw.Upgrader'
var upgrader = gw.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var orderLock sync.Mutex

// GetPrices returns current stock prices
func GetPrices(c *gin.Context) {
	stocks := make([]*models.Stock, 0, len(mockdata.Stocks))
	for _, s := range mockdata.Stocks {
		stocks = append(stocks, s)
	}
	c.JSON(200, stocks)
}

// WSHandler upgrades connection and streams live updates pushed by hub
func WSHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket upgrade failed:", err)
		return
	}

	client := &websocket.Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	// register client
	websocketHubRegister(client)

	// start writer and reader for this client
	go wsWriter(client)
	// we don't need reader for this simple feed, but keep to detect close
	go wsReader(client)
}

func websocketHubRegister(c *websocket.Client) {
	websocket.RegisterClient(c)
}

func wsWriter(client *websocket.Client) {
	for msg := range client.Send {
		// CHANGED: Now uses 'gw.TextMessage'
		if err := client.Conn.WriteMessage(gw.TextMessage, msg); err != nil {
			break
		}
	}
	// cleanup
	websocket.UnregisterClient(client)
}

func wsReader(client *websocket.Client) {
	defer websocket.UnregisterClient(client)
	client.Conn.SetReadLimit(512)
	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	for {
		_, _, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}
		// ignore incoming messages for now
	}
}

// POST /orders
func CreateOrder(c *gin.Context) {
	var payload models.Order
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid payload", "detail": err.Error()})
		return
	}
	if payload.Side != "buy" && payload.Side != "sell" {
		c.JSON(400, gin.H{"error": "side must be 'buy' or 'sell'"})
		return
	}
	if payload.Quantity <= 0 || payload.Price <= 0 {
		c.JSON(400, gin.H{"error": "quantity and price must be > 0"})
		return
	}

	orderLock.Lock()
	defer orderLock.Unlock()
	payload.ID = nextOrderID()
	payload.Time = time.Now().Unix()
	
	// CHANGED: Capitalized 'OrdersMutex'
	models.OrdersMutex.Lock()
	models.Orders = append(models.Orders, payload)
	// CHANGED: Capitalized 'OrdersMutex'
	models.OrdersMutex.Unlock()

	c.JSON(201, payload)
}

// GET /orders
func GetOrders(c *gin.Context) {
	// CHANGED: Capitalized 'OrdersMutex'
	models.OrdersMutex.Lock()
	// CHANGED: Capitalized 'OrdersMutex'
	defer models.OrdersMutex.Unlock()
	c.JSON(200, models.Orders)
}

// helpers for order ID
var (
	orderIDMu sync.Mutex
	orderID   int64 = 1
)

func nextOrderID() int64 {
	orderIDMu.Lock()
	defer orderIDMu.Unlock()
	id := orderID
	orderID++
	return id
}