package config

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type ConfigWatcher struct {
	watcher      *fsnotify.Watcher
	configLoader *ConfigLoader
	configPath   string
}

func NewConfigWatcher(configLoader *ConfigLoader, configPath string) (*ConfigWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &ConfigWatcher{
		watcher:      watcher,
		configLoader: configLoader,
		configPath:   configPath,
	}, nil
}

func (w *ConfigWatcher) Start() error {
	if err := w.watcher.Add(w.configPath); err != nil {
		return err
	}

	go w.watch()
	log.Println("Config watcher started, watching:", w.configPath)
	return nil
}

func (w *ConfigWatcher) watch() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("Config file changed:", event.Name)
				w.handleConfigChange(event.Name)
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Println("Config watcher error:", err)
		}
	}
}

func (w *ConfigWatcher) handleConfigChange(filename string) {
	base := filepath.Base(filename)

	switch base {
	case "redis.yaml":
		if err := w.configLoader.LoadRedis(); err != nil {
			log.Printf("Failed to reload redis config: %v", err)
			return
		}
		log.Println("Redis config reloaded successfully")

	case "kafka.yaml":
		if err := w.configLoader.LoadKafka(); err != nil {
			log.Printf("Failed to reload kafka config: %v", err)
			return
		}
		log.Println("Kafka config reloaded successfully")

	case "postgresql.yaml":
		if err := w.configLoader.LoadPostgreSQL(); err != nil {
			log.Printf("Failed to reload postgresql config: %v", err)
			return
		}
		log.Println("PostgreSQL config reloaded successfully")

	case "mysql.yaml":
		if err := w.configLoader.LoadMySQL(); err != nil {
			log.Printf("Failed to reload mysql config: %v", err)
			return
		}
		log.Println("MySQL config reloaded successfully")
	}

	w.configLoader.notifyChanges()
}

func (w *ConfigWatcher) Stop() error {
	return w.watcher.Close()
}
