// ==================== internal/usecase/config_usecase.go ====================
package usecase

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ConfigUsecase struct {
	configPath string
}

func NewConfigUsecase(configPath string) *ConfigUsecase {
	return &ConfigUsecase{
		configPath: configPath,
	}
}

func (u *ConfigUsecase) SaveRedisConfig(config string) error {
	return u.saveConfig("redis.yaml", config)
}

func (u *ConfigUsecase) SaveKafkaConfig(config string) error {
	return u.saveConfig("kafka.yaml", config)
}

func (u *ConfigUsecase) SavePostgreSQLConfig(config string) error {
	return u.saveConfig("postgresql.yaml", config)
}

func (u *ConfigUsecase) SaveMySQLConfig(config string) error {
	return u.saveConfig("mysql.yaml", config)
}

func (u *ConfigUsecase) saveConfig(filename, config string) error {
	path := filepath.Join(u.configPath, filename)

	// Validate YAML
	var temp interface{}
	if err := yaml.Unmarshal([]byte(config), &temp); err != nil {
		return fmt.Errorf("invalid YAML: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, []byte(config), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func (u *ConfigUsecase) GetRedisConfig() (string, error) {
	return u.readConfig("redis.yaml")
}

func (u *ConfigUsecase) GetKafkaConfig() (string, error) {
	return u.readConfig("kafka.yaml")
}

func (u *ConfigUsecase) GetPostgreSQLConfig() (string, error) {
	return u.readConfig("postgresql.yaml")
}

func (u *ConfigUsecase) GetMySQLConfig() (string, error) {
	return u.readConfig("mysql.yaml")
}

func (u *ConfigUsecase) readConfig(filename string) (string, error) {
	path := filepath.Join(u.configPath, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
