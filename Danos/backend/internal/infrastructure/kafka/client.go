package kafka

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Danos/backend/internal/domain"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaManager struct {
	mu          sync.RWMutex
	adminClient *kafka.AdminClient
	consumer    *kafka.Consumer
	config      *domain.KafkaConfig
	ctx         context.Context
	cancelFunc  context.CancelFunc
	connected   bool // Track connection state
}

func NewKafkaManager() *KafkaManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &KafkaManager{
		ctx:        ctx,
		cancelFunc: cancel,
		connected:  false,
	}
}

func (m *KafkaManager) Initialize(config *domain.KafkaConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config = config

	// Build Kafka configuration with aggressive timeouts
	kafkaConfig := m.buildKafkaConfig()

	// Create Admin Client
	adminClient, err := kafka.NewAdminClient(&kafkaConfig)
	if err != nil {
		return fmt.Errorf("failed to create kafka admin client: %w", err)
	}

	// CRITICAL: Test connection immediately with timeout
	// This prevents background retries
	// ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	// defer cancel()

	// Try to get metadata - this will fail fast if Kafka is not available
	metadata, err := adminClient.GetMetadata(nil, false, 3000) // 3 second timeout
	if err != nil {
		adminClient.Close() // Close immediately on failure
		return fmt.Errorf("failed to connect to kafka brokers: %w", err)
	}

	// Verify we have at least one broker
	if len(metadata.Brokers) == 0 {
		adminClient.Close()
		return fmt.Errorf("no kafka brokers available")
	}

	// Only create consumer if connection test passed
	consumerConfig := m.buildKafkaConfig()
	consumerConfig["group.id"] = "monitoring-consumer"
	consumerConfig["auto.offset.reset"] = "latest"
	consumerConfig["enable.auto.commit"] = false // Disable auto commit

	consumer, err := kafka.NewConsumer(&consumerConfig)
	if err != nil {
		adminClient.Close()
		return fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	// Set as connected only after everything succeeds
	m.adminClient = adminClient
	m.consumer = consumer
	m.connected = true

	log.Println("Connected to Kafka cluster using Confluent client")
	return nil
}

func (m *KafkaManager) buildKafkaConfig() kafka.ConfigMap {
	kafkaConfig := kafka.ConfigMap{
		"bootstrap.servers": strings.Join(m.config.Brokers, ","),

		// CRITICAL: Aggressive timeouts to fail fast
		"socket.timeout.ms":                  3000, // 3 seconds
		"metadata.request.timeout.ms":        3000, // 3 seconds
		"request.timeout.ms":                 3000, // 3 seconds
		"socket.connection.setup.timeout.ms": 3000, // 3 seconds
		"connections.max.idle.ms":            3000, // Close idle connections
		"reconnect.backoff.ms":               1000, // 1 second
		"reconnect.backoff.max.ms":           1000, // Max 1 second backoff

		// Disable features that cause retries
		"api.version.request":                 false, // Don't request API version
		"broker.version.fallback":             "2.0.0",
		"log.connection.close":                false, // Don't log connection closes
		"socket.keepalive.enable":             false, // Disable keepalive
		"enable.ssl.certificate.verification": false, // Skip SSL verification

		// Reduce retries to minimum
		"message.send.max.retries": 0,
		"message.timeout.ms":       3000,
		"retry.backoff.ms":         100,
		"retries":                  0,

		// Logging - suppress errors
		"log_level": 3, // 0=emerg, 1=alert, 2=crit, 3=err, 4=warning, 5=notice, 6=info, 7=debug
		"debug":     ",",
	}

	// Add security configuration
	if m.config.Security.Protocol == "SASL_SSL" {
		kafkaConfig["security.protocol"] = "SASL_SSL"
		kafkaConfig["sasl.mechanism"] = m.config.Security.SASLMechanism
		kafkaConfig["sasl.username"] = m.config.Security.Username
		kafkaConfig["sasl.password"] = m.config.Security.Password
		kafkaConfig["ssl.ca.location"] = "/etc/ssl/certs/ca-certificates.crt"
	} else if m.config.Security.Protocol == "SASL_PLAINTEXT" {
		kafkaConfig["security.protocol"] = "SASL_PLAINTEXT"
		kafkaConfig["sasl.mechanism"] = m.config.Security.SASLMechanism
		kafkaConfig["sasl.username"] = m.config.Security.Username
		kafkaConfig["sasl.password"] = m.config.Security.Password
	} else if m.config.Security.Protocol == "SSL" {
		kafkaConfig["security.protocol"] = "SSL"
		kafkaConfig["ssl.ca.location"] = "/etc/ssl/certs/ca-certificates.crt"
	}

	return kafkaConfig
}

func (m *KafkaManager) Reconnect(config *domain.KafkaConfig) error {
	// Close existing connections first
	m.Close()

	// Wait a bit for cleanup
	time.Sleep(100 * time.Millisecond)

	// Try to reconnect
	log.Println("Reconnecting Kafka clients...")
	return m.Initialize(config)
}

func (m *KafkaManager) GetMetrics() ([]domain.KafkaMetrics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.connected || m.adminClient == nil {
		return nil, fmt.Errorf("kafka client not connected")
	}

	metrics := make([]domain.KafkaMetrics, 0)

	// Get cluster metadata with timeout
	metadata, err := m.adminClient.GetMetadata(nil, false, 3000)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster metadata: %w", err)
	}

	// Create metrics for each broker
	for _, broker := range metadata.Brokers {
		metric := domain.KafkaMetrics{
			BrokerID:  int(broker.ID),
			Host:      broker.Host,
			Status:    "online",
			Timestamp: time.Now(),
		}

		// Get topics metrics
		topicMetrics, totalPartitions := m.getTopicsMetrics(metadata)
		metric.Topics = topicMetrics
		metric.TotalPartitions = totalPartitions

		// Get consumer groups metrics
		consumerMetrics, err := m.getConsumerGroupsMetrics()
		if err != nil {
			log.Printf("Error getting consumer groups metrics: %v", err)
		} else {
			metric.ConsumerGroups = consumerMetrics
		}

		// Count under-replicated and offline partitions
		metric.UnderReplicated = m.countUnderReplicatedPartitions(metadata)
		metric.OfflinePartitions = m.countOfflinePartitions(metadata)

		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func (m *KafkaManager) getTopicsMetrics(metadata *kafka.Metadata) ([]domain.KafkaTopicMetrics, int) {
	metrics := make([]domain.KafkaTopicMetrics, 0)
	totalPartitions := 0

	for topicName, topic := range metadata.Topics {
		// Only monitor configured topics
		if !m.shouldMonitorTopic(topicName) {
			continue
		}

		if topic.Error.Code() != kafka.ErrNoError {
			log.Printf("Topic %s has error: %v", topicName, topic.Error)
			continue
		}

		partitionCount := len(topic.Partitions)
		totalPartitions += partitionCount

		// Get replication factor from first partition
		replicationFactor := 0
		if partitionCount > 0 {
			replicationFactor = len(topic.Partitions[0].Replicas)
		}

		metric := domain.KafkaTopicMetrics{
			Name:              topicName,
			Partitions:        partitionCount,
			ReplicationFactor: replicationFactor,
			MessagesPerSec:    0,
			BytesInPerSec:     0,
			BytesOutPerSec:    0,
		}

		metrics = append(metrics, metric)
	}

	return metrics, totalPartitions
}

func (m *KafkaManager) getConsumerGroupsMetrics() ([]domain.KafkaConsumerMetrics, error) {
	metrics := make([]domain.KafkaConsumerMetrics, 0)

	// List consumer groups with short timeout
	ctx, cancel := context.WithTimeout(m.ctx, 3*time.Second)
	defer cancel()

	result, err := m.adminClient.ListConsumerGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list consumer groups: %w", err)
	}

	for _, group := range result.Valid {
		// Only monitor configured consumer groups
		if !m.shouldMonitorGroup(group.GroupID) {
			continue
		}

		// Describe consumer group
		groupDesc, err := m.describeConsumerGroup(group.GroupID)
		if err != nil {
			log.Printf("Error describing consumer group %s: %v", group.GroupID, err)
			continue
		}

		// Get consumer lag
		topicLags, totalLag := m.getConsumerLag(group.GroupID)

		metric := domain.KafkaConsumerMetrics{
			GroupID:   group.GroupID,
			State:     groupDesc.State,
			Members:   groupDesc.Members,
			Lag:       totalLag,
			TopicLags: topicLags,
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func (m *KafkaManager) describeConsumerGroup(groupID string) (*ConsumerGroupDescription, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 3*time.Second)
	defer cancel()

	res, err := m.adminClient.DescribeConsumerGroups(ctx, []string{groupID})
	if err != nil {
		return nil, err
	}

	if len(res.ConsumerGroupDescriptions) == 0 {
		return nil, fmt.Errorf("consumer group %s not found", groupID)
	}

	desc := res.ConsumerGroupDescriptions[0]

	return &ConsumerGroupDescription{
		GroupID: desc.GroupID,
		State:   desc.State.String(),
		Members: len(desc.Members),
	}, nil
}

type ConsumerGroupDescription struct {
	GroupID string
	State   string
	Members int
}

func (m *KafkaManager) getConsumerLag(groupID string) ([]domain.KafkaTopicLag, int64) {
	lags := make([]domain.KafkaTopicLag, 0)
	var totalLag int64

	// Get committed offsets for the consumer group
	ctx, cancel := context.WithTimeout(m.ctx, 3*time.Second)
	defer cancel()

	result, err := m.adminClient.ListConsumerGroupOffsets(
		ctx,
		[]kafka.ConsumerGroupTopicPartitions{
			{
				Group: groupID,
			},
		},
	)
	if err != nil {
		log.Printf("Error listing consumer group offsets: %v", err)
		return lags, totalLag
	}

	if len(result.ConsumerGroupsTopicPartitions) == 0 {
		return lags, totalLag
	}

	partitions := result.ConsumerGroupsTopicPartitions[0].Partitions

	for _, partition := range partitions {
		if partition.Error != nil {
			continue
		}

		// Get high water mark with timeout
		low, high, err := m.consumer.QueryWatermarkOffsets(
			*partition.Topic,
			partition.Partition,
			3000, // 3 second timeout
		)
		if err != nil {
			log.Printf("Error querying watermark offsets: %v", err)
			continue
		}

		_ = low

		committedOffset := int64(partition.Offset)
		lag := high - committedOffset
		if lag < 0 {
			lag = 0
		}

		topicLag := domain.KafkaTopicLag{
			Topic:     *partition.Topic,
			Partition: partition.Partition,
			Lag:       lag,
		}

		lags = append(lags, topicLag)
		totalLag += lag
	}

	return lags, totalLag
}

func (m *KafkaManager) countUnderReplicatedPartitions(metadata *kafka.Metadata) int {
	count := 0
	for _, topic := range metadata.Topics {
		for _, partition := range topic.Partitions {
			if len(partition.Replicas) > len(partition.Isrs) {
				count++
			}
		}
	}
	return count
}

func (m *KafkaManager) countOfflinePartitions(metadata *kafka.Metadata) int {
	count := 0
	for _, topic := range metadata.Topics {
		for _, partition := range topic.Partitions {
			if partition.Leader == -1 {
				count++
			}
		}
	}
	return count
}

func (m *KafkaManager) shouldMonitorTopic(topic string) bool {
	// Skip internal topics
	if strings.HasPrefix(topic, "__") {
		return false
	}

	// If no topics configured, monitor all
	if len(m.config.Monitoring.Topics) == 0 {
		return true
	}

	for _, t := range m.config.Monitoring.Topics {
		if t == topic {
			return true
		}
	}
	return false
}

func (m *KafkaManager) shouldMonitorGroup(group string) bool {
	// If no groups configured, monitor all
	if len(m.config.Monitoring.ConsumerGroups) == 0 {
		return true
	}

	for _, g := range m.config.Monitoring.ConsumerGroups {
		if g == group {
			return true
		}
	}
	return false
}

func (m *KafkaManager) IsConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.connected && m.adminClient != nil
}

func (m *KafkaManager) GetConnectionInfo() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.config == nil {
		return []string{}
	}
	return m.config.Brokers
}

func (m *KafkaManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.connected = false

	if m.cancelFunc != nil {
		m.cancelFunc()
	}

	if m.adminClient != nil {
		m.adminClient.Close()
		m.adminClient = nil
		log.Println("Kafka admin client closed")
	}

	if m.consumer != nil {
		if err := m.consumer.Close(); err != nil {
			log.Printf("Error closing Kafka consumer: %v", err)
		}
		m.consumer = nil
		log.Println("Kafka consumer closed")
	}

	return nil
}

// GetClusterID returns the Kafka cluster ID
func (m *KafkaManager) GetClusterID() (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.connected || m.adminClient == nil {
		return "", fmt.Errorf("kafka admin client not connected")
	}

	metadata, err := m.adminClient.GetMetadata(nil, false, 3000)
	if err != nil {
		return "", err
	}

	return metadata.OriginatingBroker.Host, nil
}

// CreateTopic creates a new Kafka topic
func (m *KafkaManager) CreateTopic(topicName string, numPartitions int, replicationFactor int) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.connected || m.adminClient == nil {
		return fmt.Errorf("kafka admin client not connected")
	}

	ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
	defer cancel()

	topicSpec := kafka.TopicSpecification{
		Topic:             topicName,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}

	results, err := m.adminClient.CreateTopics(
		ctx,
		[]kafka.TopicSpecification{topicSpec},
	)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}

	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError {
			return fmt.Errorf("failed to create topic %s: %v", result.Topic, result.Error)
		}
	}

	log.Printf("Topic %s created successfully", topicName)
	return nil
}

// DeleteTopic deletes a Kafka topic
func (m *KafkaManager) DeleteTopic(topicName string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.connected || m.adminClient == nil {
		return fmt.Errorf("kafka admin client not connected")
	}

	ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
	defer cancel()

	results, err := m.adminClient.DeleteTopics(
		ctx,
		[]string{topicName},
	)
	if err != nil {
		return fmt.Errorf("failed to delete topic: %w", err)
	}

	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError {
			return fmt.Errorf("failed to delete topic %s: %v", result.Topic, result.Error)
		}
	}

	log.Printf("Topic %s deleted successfully", topicName)
	return nil
}

// GetTopicConfig retrieves configuration for a specific topic
func (m *KafkaManager) GetTopicConfig(topicName string) (map[string]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.connected || m.adminClient == nil {
		return nil, fmt.Errorf("kafka admin client not connected")
	}

	ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	defer cancel()

	resource := kafka.ConfigResource{
		Type: kafka.ResourceTopic,
		Name: topicName,
	}

	results, err := m.adminClient.DescribeConfigs(ctx, []kafka.ConfigResource{resource})
	if err != nil {
		return nil, fmt.Errorf("failed to describe topic config: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("topic %s not found", topicName)
	}

	configMap := make(map[string]string)
	for _, entry := range results[0].Config {
		configMap[entry.Name] = entry.Value
	}

	return configMap, nil
}
