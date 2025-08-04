package config

import (
	"log/slog"
	"os"
	"time"
)

const (
	DefaultTimeout = 3 * time.Second
)

const (
	prod = "prod"
	Mex  = "mex"

	regionAWS = "us-east-1"
)

type Config struct {
	// NP
	Env     string
	Country string

	// global
	TracingEnabled bool

	// features (ejemplo, podés quitar/ajustar según necesidades)
	GetStatusEnabled  bool
	GetItemEnabled    bool
	CreateItemEnabled bool

	// configs específicos de esta lambda
	DynamoTableName string

	// AWS
	AWSRegion string
}

func (c *Config) IsProd() bool { return c.Env == prod }
func (c *Config) IsLow() bool  { return !c.IsProd() }

func (c *Config) DynamoTable() string {
	return c.DynamoTableName
}

func (c *Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Bool("tracing_enabled", c.TracingEnabled),
		slog.Bool("public_get_status_enabled", c.GetStatusEnabled),
		slog.Bool("public_get_item_enabled", c.GetItemEnabled),
		slog.Bool("public_create_item_enabled", c.CreateItemEnabled),
		slog.String("aws_region", c.AWSRegion),
		slog.String("env", c.Env),
		slog.String("country", c.Country),
		slog.String("dynamo_table_name", c.DynamoTableName),
	)
}

func LoadConfig() *Config {
	return &Config{
		Env:               os.Getenv("NP_DIMENSION_ENVIRONMENT"),
		Country:           os.Getenv("NP_DIMENSION_COUNTRY"),
		TracingEnabled:    parseBool(os.Getenv("TRACING_ENABLED")),
		GetStatusEnabled:  parseBool(os.Getenv("GET_STATUS_ENABLED")),
		GetItemEnabled:    parseBool(os.Getenv("GET_ITEM_ENABLED")),
		CreateItemEnabled: parseBool(os.Getenv("CREATE_ITEM_ENABLED")),
		DynamoTableName:   os.Getenv("DYNAMODB_TABLE_NAME"),
		AWSRegion:         regionAWS,
	}
}

func parseBool(s string) bool {
	return s == "true"
}
