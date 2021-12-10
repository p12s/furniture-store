package config

import "github.com/kelseyhightower/envconfig"

// Config
type Config struct {
	DB     DB
	Server Server
	Auth   Auth
	Broker Broker
	Env    Env
}

// DB
type DB struct {
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

// Broker
type Broker struct {
	// TopicPrefix      string `envconfig:"BROKER_TOPIC_PREFIX" required:"true"`
	Brokers          string `envconfig:"BROKER_BROKERS" required:"true"`
	Username         string `envconfig:"BROKER_USERNAME" required:"true"`
	Password         string `envconfig:"BROKER_PASSWORD" required:"true"`
	TopicAccountBE   string `envconfig:"BROKER_TOPIC_ACCOUNT_BE" required:"true"`
	TopicAccountCUD  string `envconfig:"BROKER_TOPIC_ACCOUNT_CUD" required:"true"`
	TopicProductBE   string `envconfig:"BROKER_TOPIC_PRODUCT_BE" required:"true"`
	TopicProductCUD  string `envconfig:"BROKER_TOPIC_PRODUCT_CUD" required:"true"`
	TopicOrderBE     string `envconfig:"BROKER_TOPIC_ORDER_BE" required:"true"`
	TopicOrderCUD    string `envconfig:"BROKER_TOPIC_ORDER_CUD" required:"true"`
	TopicDeliveryBE  string `envconfig:"BROKER_TOPIC_DELIVERY_BE" required:"true"`
	TopicDeliveryCUD string `envconfig:"BROKER_TOPIC_DELIVERY_CUD" required:"true"`
	TopicBillingBE   string `envconfig:"BROKER_TOPIC_BILLING_BE" required:"true"`
	TopicBillingCUD  string `envconfig:"BROKER_TOPIC_BILLING_CUD" required:"true"`
	GroupId          string `envconfig:"BROKER_GROUP_ID" required:"true"`
}

// Env
type Env struct {
	Current string `envconfig:"ENV_CURRENT" required:"true"`
	Dev     string `envconfig:"ENV_DEV" required:"true"`
	Qa      string `envconfig:"ENV_QA" required:"true"`
	Prod    string `envconfig:"ENV_PROD" required:"true"`
}

// New - contructor
func New() (*Config, error) {
	cfg := new(Config)

	if err := envconfig.Process("db", &cfg.DB); err != nil {
		return nil, err
	}

	if err := envconfig.Process("server", &cfg.Server); err != nil {
		return nil, err
	}

	if err := envconfig.Process("auth", &cfg.Auth); err != nil {
		return nil, err
	}

	if err := envconfig.Process("broker", &cfg.Broker); err != nil {
		return nil, err
	}

	if err := envconfig.Process("env", &cfg.Env); err != nil {
		return nil, err
	}

	return cfg, nil
}
