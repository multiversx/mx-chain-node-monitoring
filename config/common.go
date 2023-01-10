package config

import (
	"github.com/multiversx/mx-chain-core-go/core"
)

// LoadConfig returns email configuration by reading the config file provided
func LoadConfig(filepath string) (*GeneralConfig, error) {
	cfg := &GeneralConfig{}
	err := core.LoadTomlFile(cfg, filepath)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
