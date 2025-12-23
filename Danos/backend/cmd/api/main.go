// ==================== cmd/api/main.go (SIMPLIFIED) ====================
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Danos/backend/internal/delivery/handler/http"
	"github.com/Danos/backend/internal/infrastructure/config"
	"github.com/Danos/backend/internal/infrastructure/kafka"
	"github.com/Danos/backend/internal/infrastructure/mysql"
	"github.com/Danos/backend/internal/infrastructure/postgres"
	"github.com/Danos/backend/internal/infrastructure/redis"
	"github.com/Danos/backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize config loader (loads config files WITHOUT connecting)
	configPath := "./configs"
	configLoader := config.NewConfigLoader(configPath)

	// Load all configurations (just read files, no connection)
	log.Println("Loading configuration files...")
	if err := configLoader.LoadAll(); err != nil {
		log.Fatalf("Failed to load configurations: %v", err)
	}
	log.Println("✓ Configuration files loaded successfully")

	// Initialize managers (create instances, no connection yet)
	redisManager := redis.NewRedisManager()
	kafkaManager := kafka.NewKafkaManager()
	postgresManager := postgres.NewPostgresManager()
	mysqlManager := mysql.NewMySQLManager()

	// Connect to services (dengan error handling, tidak fatal)
	log.Println("\nAttempting to connect to services...")

	if err := redisManager.Initialize(configLoader.GetRedis()); err != nil {
		log.Printf("⚠ Redis connection failed: %v (service will start without Redis)", err)
	} else {
		log.Println("✓ Redis connected successfully")
	}

	if err := kafkaManager.Initialize(configLoader.GetKafka()); err != nil {
		log.Printf("⚠ Kafka connection failed: %v (service will start without Kafka)", err)
	} else {
		log.Println("✓ Kafka connected successfully")
	}

	if err := postgresManager.Initialize(configLoader.GetPostgreSQL()); err != nil {
		log.Printf("⚠ PostgreSQL connection failed: %v (service will start without PostgreSQL)", err)
	} else {
		log.Println("✓ PostgreSQL connected successfully")
	}

	if err := mysqlManager.Initialize(configLoader.GetMySQL()); err != nil {
		log.Printf("⚠ MySQL connection failed: %v (service will start without MySQL)", err)
	} else {
		log.Println("✓ MySQL connected successfully")
	}

	// Setup config change callbacks
	configLoader.OnChange(func() {
		log.Println("\n🔄 Configuration changed, attempting reconnection...")

		if err := redisManager.Reconnect(configLoader.GetRedis()); err != nil {
			log.Printf("⚠ Redis reconnection failed: %v", err)
		} else {
			log.Println("✓ Redis reconnected successfully")
		}

		if err := kafkaManager.Reconnect(configLoader.GetKafka()); err != nil {
			log.Printf("⚠ Kafka reconnection failed: %v", err)
		} else {
			log.Println("✓ Kafka reconnected successfully")
		}

		if err := postgresManager.Reconnect(configLoader.GetPostgreSQL()); err != nil {
			log.Printf("⚠ PostgreSQL reconnection failed: %v", err)
		} else {
			log.Println("✓ PostgreSQL reconnected successfully")
		}

		if err := mysqlManager.Reconnect(configLoader.GetMySQL()); err != nil {
			log.Printf("⚠ MySQL reconnection failed: %v", err)
		} else {
			log.Println("✓ MySQL reconnected successfully")
		}
	})

	// Start config watcher
	watcher, err := config.NewConfigWatcher(configLoader, configPath)
	if err != nil {
		log.Fatalf("Failed to create config watcher: %v", err)
	}

	if err := watcher.Start(); err != nil {
		log.Fatalf("Failed to start config watcher: %v", err)
	}
	defer watcher.Stop()

	// Initialize usecases
	monitoringUsecase := usecase.NewMonitoringUsecase(
		redisManager,
		kafkaManager,
		postgresManager,
		mysqlManager,
	)
	configUsecase := usecase.NewConfigUsecase(configPath)

	// Initialize HTTP handler (now with managers and configLoader)
	handler := http.NewHandler(
		monitoringUsecase,
		configUsecase,
		redisManager,
		kafkaManager,
		postgresManager,
		mysqlManager,
		configLoader,
	)

	// Initialize WebSocket manager
	wsManager := http.NewWebSocketManager(monitoringUsecase)
	wsManager.Start()
	defer wsManager.Stop()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Monitoring Backend API",
		ServerHeader: "Fiber",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	// Setup middleware
	http.SetupMiddleware(app)

	// Setup routes (all endpoints now in routes.go)
	http.SetupRoutes(app, handler, wsManager)

	log.Println("╔════════════════════════════════════════════════════════════════╗")
	log.Println("║      Monitoring Backend - Fiber Framework                      ║")
	log.Println("╠════════════════════════════════════════════════════════════════╣")
	log.Println("║  HTTP Server: http://localhost:8085                            ║")
	log.Println("║  WebSocket:   ws://localhost:8085/api/v1/ws/metrics            ║")
	log.Println("║  Health:      http://localhost:8085/health                     ║")
	log.Println("║  Status:      http://localhost:8085/api/v1/connections/status  ║")
	log.Println("║  Connect:     POST /api/v1/connections/connect                 ║")
	log.Println("║  Disconnect:  POST /api/v1/connections/disconnect              ║")
	log.Println("╚════════════════════════════════════════════════════════════════╝")
	log.Println("✓ Server started successfully")
	log.Println("✓ Watching for configuration changes...")

	// Start HTTP server in goroutine
	go func() {
		if err := app.Listen(":8085"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("\n🛑 Shutting down gracefully...")

	// Shutdown Fiber app
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	// Close connections
	redisManager.Close()
	kafkaManager.Close()
	postgresManager.Close()
	mysqlManager.Close()

	log.Println("✓ Shutdown complete")
}
