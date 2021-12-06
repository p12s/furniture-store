package config

import "github.com/kelseyhightower/envconfig"

// Config
type Config struct {
	Db           Db
	Server       Server
	Auth         Auth
	Cloudkarafka Cloudkarafka
}

// Db
type Db struct {
	Driver string `envconfig:"DB_DRIVER" required:"true"`
}

// Server
type Server struct {
	Port int `envconfig:"SERVER_PORT" required:"true"`
}

// Auth
type Auth struct {
	Salt       string `envconfig:"AUTH_SALT" required:"true"`
	TokenTTL   int    `envconfig:"AUTH_TOKEN_TTL" required:"true"`
	SigningKey string `envconfig:"AUTH_SIGNING_KEY" required:"true"`
}

// Cloudkarafka
type Cloudkarafka struct {
	TopicPrefix string `envconfig:"CLOUDKARAFKA_TOPIC_PREFIX" required:"true"`
	Brokers     string `envconfig:"CLOUDKARAFKA_BROKERS" required:"true"`
	Username    string `envconfig:"CLOUDKARAFKA_USERNAME" required:"true"`
	Password    string `envconfig:"CLOUDKARAFKA_PASSWORD" required:"true"`
}

// New - contructor
func New() (*Config, error) {
	cfg := new(Config)

	if err := envconfig.Process("db", &cfg.Db); err != nil {
		return nil, err
	}

	if err := envconfig.Process("server", &cfg.Server); err != nil {
		return nil, err
	}

	if err := envconfig.Process("auth", &cfg.Auth); err != nil {
		return nil, err
	}

	if err := envconfig.Process("cloudkarafka", &cfg.Cloudkarafka); err != nil {
		return nil, err
	}

	return cfg, nil
}
