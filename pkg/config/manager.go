package config

import (
	"fmt"
	"os"
)

type ConfigManager struct {
	requiredEnvSlice []string
	config           map[string]string
}

type ConfigError struct {
	envName string
}

func (ce *ConfigError) Error() string {
	return fmt.Sprintf(`missing env %s`, ce.envName)
}

func NewConfigManager(envSlice []string) *ConfigManager {
	return &ConfigManager{
		requiredEnvSlice: envSlice,
		config:           make(map[string]string),
	}
}

func (m *ConfigManager) CheckAndLoadConfig() error {
	for _, env := range m.requiredEnvSlice {

		v := os.Getenv(env)
		if v == "" {
			return &ConfigError{envName: env}
		}

		m.config[env] = v
	}
	return nil
}
