package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jetbuild/runner/internal/flow"
)

type Config struct {
	Flow flow.Flow `json:"flow"`
}

func (c *Config) Load() error {
	e, ok := os.LookupEnv("config")
	if !ok {
		return fmt.Errorf("environment variable 'config' does not exist")
	}

	if err := json.Unmarshal([]byte(e), c); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}
