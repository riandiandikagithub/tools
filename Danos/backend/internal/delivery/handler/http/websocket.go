// ==================== internal/delivery/http/websocket.go ====================
package http

import (
	"log"
	"sync"
	"time"

	"github.com/Danos/backend/internal/usecase"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type WebSocketManager struct {
	clients           map[*websocket.Conn]bool
	mu                sync.RWMutex
	monitoringUsecase *usecase.MonitoringUsecase
	broadcast         chan interface{}
	ticker            *time.Ticker
	stopChan          chan bool
}

func NewWebSocketManager(monitoringUsecase *usecase.MonitoringUsecase) *WebSocketManager {
	return &WebSocketManager{
		clients:           make(map[*websocket.Conn]bool),
		monitoringUsecase: monitoringUsecase,
		broadcast:         make(chan interface{}),
		ticker:            time.NewTicker(5 * time.Second),
		stopChan:          make(chan bool),
	}
}

// WebSocketUpgrade middleware to check if request is websocket upgrade
func WebSocketUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// HandleWebSocket handles websocket connections
func (wsm *WebSocketManager) HandleWebSocket(c *websocket.Conn) {
	// Register client
	wsm.mu.Lock()
	wsm.clients[c] = true
	wsm.mu.Unlock()

	log.Printf("New WebSocket client connected. Total clients: %d", len(wsm.clients))

	// Remove client on disconnect
	defer func() {
		wsm.mu.Lock()
		delete(wsm.clients, c)
		wsm.mu.Unlock()
		c.Close()
		log.Printf("WebSocket client disconnected. Total clients: %d", len(wsm.clients))
	}()

	// Send initial metrics
	metrics := wsm.monitoringUsecase.GetAllMetrics()
	if err := c.WriteJSON(metrics); err != nil {
		log.Printf("Error sending initial metrics: %v", err)
		return
	}

	// Keep connection alive and listen for messages
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}

// Start begins broadcasting metrics to all connected clients
func (wsm *WebSocketManager) Start() {
	go wsm.broadcastMetrics()
	log.Println("WebSocket manager started, broadcasting every 5 seconds")
}

func (wsm *WebSocketManager) broadcastMetrics() {
	for {
		select {
		case <-wsm.ticker.C:
			wsm.mu.RLock()
			if len(wsm.clients) == 0 {
				wsm.mu.RUnlock()
				continue
			}
			wsm.mu.RUnlock()

			// Get latest metrics
			metrics := wsm.monitoringUsecase.GetAllMetrics()

			// Broadcast to all clients
			wsm.mu.RLock()
			for client := range wsm.clients {
				err := client.WriteJSON(metrics)
				if err != nil {
					log.Printf("Error broadcasting to client: %v", err)
					wsm.mu.RUnlock()
					wsm.mu.Lock()
					delete(wsm.clients, client)
					client.Close()
					wsm.mu.Unlock()
					wsm.mu.RLock()
				}
			}
			wsm.mu.RUnlock()

		case <-wsm.stopChan:
			return
		}
	}
}

// Stop stops the websocket manager
func (wsm *WebSocketManager) Stop() {
	wsm.ticker.Stop()
	wsm.stopChan <- true

	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	for client := range wsm.clients {
		client.Close()
		delete(wsm.clients, client)
	}

	log.Println("WebSocket manager stopped")
}
