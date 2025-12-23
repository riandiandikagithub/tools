// ==================== internal/delivery/http/handler.go (FIXED) ====================
package http

import (
	"github.com/Danos/backend/internal/infrastructure/config"
	"github.com/Danos/backend/internal/infrastructure/kafka"
	"github.com/Danos/backend/internal/infrastructure/mysql"
	"github.com/Danos/backend/internal/infrastructure/postgres"
	"github.com/Danos/backend/internal/infrastructure/redis"
	"github.com/Danos/backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	monitoringUsecase *usecase.MonitoringUsecase
	configUsecase     *usecase.ConfigUsecase
	// Add managers for connection control
	redisManager    *redis.RedisManager
	kafkaManager    *kafka.KafkaManager
	postgresManager *postgres.PostgresManager
	mysqlManager    *mysql.MySQLManager
	configLoader    *config.ConfigLoader // Use concrete type instead of interface{}
}

func NewHandler(
	monitoringUsecase *usecase.MonitoringUsecase,
	configUsecase *usecase.ConfigUsecase,
	redisManager *redis.RedisManager,
	kafkaManager *kafka.KafkaManager,
	postgresManager *postgres.PostgresManager,
	mysqlManager *mysql.MySQLManager,
	configLoader *config.ConfigLoader, // Concrete type
) *Handler {
	return &Handler{
		monitoringUsecase: monitoringUsecase,
		configUsecase:     configUsecase,
		redisManager:      redisManager,
		kafkaManager:      kafkaManager,
		postgresManager:   postgresManager,
		mysqlManager:      mysqlManager,
		configLoader:      configLoader,
	}
}

// ==================== Response Helpers ====================
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func successResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Response{
		Success: true,
		Data:    data,
	})
}

func successMessageResponse(c *fiber.Ctx, message string) error {
	return c.JSON(Response{
		Success: true,
		Message: message,
	})
}

func errorResponse(c *fiber.Ctx, status int, err string) error {
	return c.Status(status).JSON(Response{
		Success: false,
		Error:   err,
	})
}

// ==================== Monitoring Endpoints ====================

func (h *Handler) GetAllMetrics(c *fiber.Ctx) error {
	metrics := h.monitoringUsecase.GetAllMetrics()
	return successResponse(c, metrics)
}

func (h *Handler) GetRedisMetrics(c *fiber.Ctx) error {
	metrics := h.monitoringUsecase.GetAllMetrics()
	return successResponse(c, fiber.Map{
		"redis": metrics.Redis,
	})
}

func (h *Handler) GetClusterOverview(c *fiber.Ctx) error {
	metrics, err := h.monitoringUsecase.GetClusterOverview(c.Context())
	if err != nil {
		return errorResponse(c, fiber.Map{
			"code": fiber.ErrInternalServerError.Code,
			"msg":  fiber.ErrInternalServerError.Message,
		})
	}
	return successResponse(c, fiber.Map{
		"redis": metrics,
	})
}
func (h *Handler) GetKafkaMetrics(c *fiber.Ctx) error {
	metrics := h.monitoringUsecase.GetAllMetrics()
	return successResponse(c, fiber.Map{
		"kafka": metrics.Kafka,
	})
}

func (h *Handler) GetPostgreSQLMetrics(c *fiber.Ctx) error {
	metrics := h.monitoringUsecase.GetAllMetrics()
	return successResponse(c, fiber.Map{
		"postgresql": metrics.PostgreSQL,
	})
}

func (h *Handler) GetMySQLMetrics(c *fiber.Ctx) error {
	metrics := h.monitoringUsecase.GetAllMetrics()
	return successResponse(c, fiber.Map{
		"mysql": metrics.MySQL,
	})
}

// ==================== Config Endpoints ====================

type SaveConfigRequest struct {
	Config string `json:"config" validate:"required"`
}

func (h *Handler) SaveRedisConfig(c *fiber.Ctx) error {
	var req SaveConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	if req.Config == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Config field is required")
	}

	if err := h.configUsecase.SaveRedisConfig(req.Config); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return successMessageResponse(c, "Redis configuration saved successfully")
}

func (h *Handler) SaveKafkaConfig(c *fiber.Ctx) error {
	var req SaveConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	if req.Config == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Config field is required")
	}

	if err := h.configUsecase.SaveKafkaConfig(req.Config); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return successMessageResponse(c, "Kafka configuration saved successfully")
}

func (h *Handler) SavePostgreSQLConfig(c *fiber.Ctx) error {
	var req SaveConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	if req.Config == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Config field is required")
	}

	if err := h.configUsecase.SavePostgreSQLConfig(req.Config); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return successMessageResponse(c, "PostgreSQL configuration saved successfully")
}

func (h *Handler) SaveMySQLConfig(c *fiber.Ctx) error {
	var req SaveConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	if req.Config == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Config field is required")
	}

	if err := h.configUsecase.SaveMySQLConfig(req.Config); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return successMessageResponse(c, "MySQL configuration saved successfully")
}

func (h *Handler) GetRedisConfig(c *fiber.Ctx) error {
	config, err := h.configUsecase.GetRedisConfig()
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return successResponse(c, config)
}

func (h *Handler) GetKafkaConfig(c *fiber.Ctx) error {
	config, err := h.configUsecase.GetKafkaConfig()
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return successResponse(c, config)
}

func (h *Handler) GetPostgreSQLConfig(c *fiber.Ctx) error {
	config, err := h.configUsecase.GetPostgreSQLConfig()
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return successResponse(c, config)
}

func (h *Handler) GetMySQLConfig(c *fiber.Ctx) error {
	config, err := h.configUsecase.GetMySQLConfig()
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return successResponse(c, config)
}

// ==================== Connection Control Endpoints ====================

type ConnectRequest struct {
	Service string `json:"service" validate:"required"` // redis, kafka, postgresql, mysql, all
}

// ConnectService handles manual service connection
func (h *Handler) ConnectService(c *fiber.Ctx) error {
	var req ConnectRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	if req.Service == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Service field is required")
	}

	results := make(map[string]interface{})

	switch req.Service {
	case "redis":
		if err := h.redisManager.Initialize(h.configLoader.GetRedis()); err != nil {
			results["redis"] = map[string]interface{}{"status": "failed", "error": err.Error()}
		} else {
			results["redis"] = map[string]interface{}{"status": "connected"}
		}

	case "kafka":
		if err := h.kafkaManager.Initialize(h.configLoader.GetKafka()); err != nil {
			results["kafka"] = map[string]interface{}{"status": "failed", "error": err.Error()}
		} else {
			results["kafka"] = map[string]interface{}{"status": "connected"}
		}

	case "postgresql":
		if err := h.postgresManager.Initialize(h.configLoader.GetPostgreSQL()); err != nil {
			results["postgresql"] = map[string]interface{}{"status": "failed", "error": err.Error()}
		} else {
			results["postgresql"] = map[string]interface{}{"status": "connected"}
		}

	case "mysql":
		if err := h.mysqlManager.Initialize(h.configLoader.GetMySQL()); err != nil {
			results["mysql"] = map[string]interface{}{"status": "failed", "error": err.Error()}
		} else {
			results["mysql"] = map[string]interface{}{"status": "connected"}
		}

	case "all":
		// Connect all services
		if err := h.redisManager.Initialize(h.configLoader.GetRedis()); err != nil {
			results["redis"] = map[string]interface{}{"status": "failed", "error": err.Error()}
		} else {
			results["redis"] = map[string]interface{}{"status": "connected"}
		}

		if err := h.kafkaManager.Initialize(h.configLoader.GetKafka()); err != nil {
			results["kafka"] = map[string]interface{}{"status": "failed", "error": err.Error()}
		} else {
			results["kafka"] = map[string]interface{}{"status": "connected"}
		}

		if err := h.postgresManager.Initialize(h.configLoader.GetPostgreSQL()); err != nil {
			results["postgresql"] = map[string]interface{}{"status": "failed", "error": err.Error()}
		} else {
			results["postgresql"] = map[string]interface{}{"status": "connected"}
		}

		if err := h.mysqlManager.Initialize(h.configLoader.GetMySQL()); err != nil {
			results["mysql"] = map[string]interface{}{"status": "failed", "error": err.Error()}
		} else {
			results["mysql"] = map[string]interface{}{"status": "connected"}
		}

	default:
		return errorResponse(c, fiber.StatusBadRequest, "Invalid service name. Use: redis, kafka, postgresql, mysql, or all")
	}

	return successResponse(c, results)
}

// GetConnectionStatus returns current connection status of all services
func (h *Handler) GetConnectionStatus(c *fiber.Ctx) error {
	status := make(map[string]interface{})
	mode := c.Query("mode", "single") // default: auto
	// Check Redis
	if h.redisManager.IsConnected(mode) {
		status["redis"] = map[string]interface{}{
			"status": "connected",
			"nodes":  h.redisManager.GetConnectionInfo(mode),
		}
	} else {
		status["redis"] = map[string]interface{}{
			"status": "disconnected",
		}
	}

	// Check Kafka
	if h.kafkaManager.IsConnected() {
		status["kafka"] = map[string]interface{}{
			"status":  "connected",
			"brokers": h.kafkaManager.GetConnectionInfo(),
		}
	} else {
		status["kafka"] = map[string]interface{}{
			"status": "disconnected",
		}
	}

	// Check PostgreSQL
	pgStatus := h.postgresManager.GetConnectionStatus()
	status["postgresql"] = pgStatus

	// Check MySQL
	mysqlStatus := h.mysqlManager.GetConnectionStatus()
	status["mysql"] = mysqlStatus

	return successResponse(c, status)
}

// DisconnectService handles manual service disconnection
func (h *Handler) DisconnectService(c *fiber.Ctx) error {
	var req ConnectRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	results := make(map[string]interface{})

	switch req.Service {
	case "redis":
		if err := h.redisManager.Close(); err != nil {
			results["redis"] = map[string]interface{}{"status": "error", "error": err.Error()}
		} else {
			results["redis"] = map[string]interface{}{"status": "disconnected"}
		}

	case "kafka":
		if err := h.kafkaManager.Close(); err != nil {
			results["kafka"] = map[string]interface{}{"status": "error", "error": err.Error()}
		} else {
			results["kafka"] = map[string]interface{}{"status": "disconnected"}
		}

	case "postgresql":
		if err := h.postgresManager.Close(); err != nil {
			results["postgresql"] = map[string]interface{}{"status": "error", "error": err.Error()}
		} else {
			results["postgresql"] = map[string]interface{}{"status": "disconnected"}
		}

	case "mysql":
		if err := h.mysqlManager.Close(); err != nil {
			results["mysql"] = map[string]interface{}{"status": "error", "error": err.Error()}
		} else {
			results["mysql"] = map[string]interface{}{"status": "disconnected"}
		}

	case "all":
		h.redisManager.Close()
		h.kafkaManager.Close()
		h.postgresManager.Close()
		h.mysqlManager.Close()
		results["all"] = map[string]interface{}{"status": "disconnected"}

	default:
		return errorResponse(c, fiber.StatusBadRequest, "Invalid service name")
	}

	return successResponse(c, results)
}

// ==================== Health Check ====================

func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Monitoring service is running",
	})
}
