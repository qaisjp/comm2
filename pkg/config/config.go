package config

type Config struct {
	LogLevel string `default:"debug"`

	Postgres PostgresConfig

	JWTSecret string `required:"true"`
	Address   string `default:"0.0.0.0:8080"`
}

// PostgresConfig contains all configuration data for a PostgreSQL connection
type PostgresConfig struct {
	ConnectionString string `required:"true"`
	// User             string `required:"true"`
	// Database         string `required:"true"`
	// Password         string
}
