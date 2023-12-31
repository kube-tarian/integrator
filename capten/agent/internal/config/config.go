package config

import (
	"github.com/kelseyhightower/envconfig"
)

type SericeConfig struct {
	Host                      string `envconfig:"HOST" default:"0.0.0.0"`
	Port                      int    `envconfig:"PORT" default:"9091"`
	Mode                      string `envconfig:"MODE" default:"production"`
	AuthEnabled               bool   `envconfig:"AUTH_ENABLED" default:"false"`
	CrossplaneSyncJobEnabled  bool   `envconfig:"CROSSPLANE_SYNC_JOB_ENABLED" default:"true"`
	CrossplaneSyncJobInterval string `envconfig:"CROSSPLANE_SYNC_JOB_INTERVAL" default:"@every 1h"`
	DomainName                string `envconfig:"DOMAIN_NAME" default:".example.com"`
}

func GetServiceConfig() (*SericeConfig, error) {
	cfg := &SericeConfig{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
