package usecase

import (
	"log"
	"sync"

	"github.com/Danos/backend/internal/domain"
	"github.com/Danos/backend/internal/infrastructure/kafka"
	"github.com/Danos/backend/internal/infrastructure/mysql"
	"github.com/Danos/backend/internal/infrastructure/postgres"
	"github.com/Danos/backend/internal/infrastructure/redis"
	"github.com/Danos/backend/internal/repository"
)

type MonitoringUsecase struct {
	redisManager    *redis.RedisManager
	kafkaManager    *kafka.KafkaManager
	postgresManager *postgres.PostgresManager
	mysqlManager    *mysql.MySQLManager
}

func NewMonitoringUsecase(
	redisManager *redis.RedisManager,
	kafkaManager *kafka.KafkaManager,
	postgresManager *postgres.PostgresManager,
	mysqlManager *mysql.MySQLManager,
) *MonitoringUsecase {
	return &MonitoringUsecase{
		redisManager:    redisManager,
		kafkaManager:    kafkaManager,
		postgresManager: postgresManager,
		mysqlManager:    mysqlManager,
	}
}

func (u *MonitoringUsecase) GetAllMetrics() domain.MetricsResponse {
	var wg sync.WaitGroup
	response := domain.MetricsResponse{}

	// Get Redis metrics
	wg.Add(5)
	go func() {
		defer wg.Done()
		metrics, err := u.redisManager.GetMetrics("single")
		if err != nil {
			log.Printf("Error getting Redis metrics: %v", err)
			return
		}
		response.Redis = metrics
	}()
	go func() {
		defer wg.Done()
		metrics, err := u.redisManager.GetMetrics("cluster")
		if err != nil {
			log.Printf("Error getting Redis metrics: %v", err)
			return
		}
		response.Redis = append(response.Redis, metrics...)
	}()

	// Get Kafka metrics

	go func() {
		defer wg.Done()
		metrics, err := u.kafkaManager.GetMetrics()
		if err != nil {
			log.Printf("Error getting Kafka metrics: %v", err)
			return
		}
		response.Kafka = metrics
	}()

	// Get PostgreSQL metrics

	go func() {
		defer wg.Done()
		metrics := u.getPostgreSQLMetrics()
		response.PostgreSQL = metrics
	}()

	// Get MySQL metrics

	go func() {
		defer wg.Done()
		metrics := u.getMySQLMetrics()
		response.MySQL = metrics
	}()

	wg.Wait()
	return response
}

func (u *MonitoringUsecase) getPostgreSQLMetrics() []domain.PostgreSQLMetrics {
	clients := u.postgresManager.GetAllClients()
	metrics := make([]domain.PostgreSQLMetrics, 0)

	for name, db := range clients {
		repo := repository.NewPostgresRepository(db)
		metric, err := repo.GetMetrics(name, "")
		if err != nil {
			log.Printf("Error getting PostgreSQL metrics for %s: %v", name, err)
			continue
		}
		metrics = append(metrics, metric)
	}

	return metrics
}

func (u *MonitoringUsecase) getMySQLMetrics() []domain.MySQLMetrics {
	clients := u.mysqlManager.GetAllClients()
	metrics := make([]domain.MySQLMetrics, 0)

	for name, db := range clients {
		repo := repository.NewMySQLRepository(db)
		metric, err := repo.GetMetrics(name, "")
		if err != nil {
			log.Printf("Error getting MySQL metrics for %s: %v", name, err)
			continue
		}
		metrics = append(metrics, metric)
	}

	return metrics
}
