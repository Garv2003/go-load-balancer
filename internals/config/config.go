package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port                string   `json:"port" yaml:"port"`
	Strategy            string   `json:"strategy" yaml:"strategy"`
	HealthCheckInterval int      `json:"healthCheckInterval" yaml:"healthCheckInterval"`
	Servers             []string `json:"servers" yaml:"servers"`
	mu                  sync.RWMutex
}

func (c *Config) ReloadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	ext := filepath.Ext(path)
	newCfg := Config{}

	switch ext {
	case ".json":
		if err := json.NewDecoder(file).Decode(&newCfg); err != nil {
			return err
		}
	case ".yaml", ".yml":
		if err := yaml.NewDecoder(file).Decode(&newCfg); err != nil {
			return err
		}
	default:
		return errors.New("unsupported config file format")
	}

	c.mu.Lock()
	c.Port = newCfg.Port
	c.Strategy = newCfg.Strategy
	c.HealthCheckInterval = newCfg.HealthCheckInterval
	c.Servers = newCfg.Servers
	c.mu.Unlock()

	return nil
}

func (c *Config) SafeGet() Config {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return Config{
		Port:                c.Port,
		Strategy:            c.Strategy,
		HealthCheckInterval: c.HealthCheckInterval,
		Servers:             append([]string{}, c.Servers...), // deep copy
	}
}
