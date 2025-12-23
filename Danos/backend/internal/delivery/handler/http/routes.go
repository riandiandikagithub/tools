// ==================== internal/delivery/http/routes.go (UPDATED) ====================
package http

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, handler *Handler, wsManager *WebSocketManager) {
	// Health check
	app.Get("/health", handler.HealthCheck)

	// API v1 routes
	api := app.Group("/api/v1")

	// Monitoring endpoints
	monitoring := api.Group("/monitoring")
	{
		monitoring.Get("/metrics", handler.GetAllMetrics)
		monitoring.Get("/redis", handler.GetRedisMetrics)
		monitoring.Get("/kafka", handler.GetKafkaMetrics)
		monitoring.Get("/postgresql", handler.GetPostgreSQLMetrics)
		monitoring.Get("/mysql", handler.GetMySQLMetrics)
	}

	// Config endpoints
	config := api.Group("/config")
	{
		// Redis
		config.Get("/redis", handler.GetRedisConfig)
		config.Post("/redis", handler.SaveRedisConfig)

		// Kafka
		config.Get("/kafka", handler.GetKafkaConfig)
		config.Post("/kafka", handler.SaveKafkaConfig)

		// PostgreSQL
		config.Get("/postgresql", handler.GetPostgreSQLConfig)
		config.Post("/postgresql", handler.SavePostgreSQLConfig)

		// MySQL
		config.Get("/mysql", handler.GetMySQLConfig)
		config.Post("/mysql", handler.SaveMySQLConfig)
	}

	// Connection control endpoints
	connections := api.Group("/connections")
	{
		connections.Get("/status", handler.GetConnectionStatus)
		connections.Post("/connect", handler.ConnectService)
		connections.Post("/disconnect", handler.DisconnectService)
	}

	// WebSocket endpoint
	api.Get("/ws/metrics", WebSocketUpgrade, websocket.New(wsManager.HandleWebSocket))
}
