package config

import (
	"fmt"
	"os"
	"user_service/pkg/logging"

	"github.com/spf13/viper"
)

type PostgresConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Db       string `mapstructure:"db"`
	Host     string
	Port     string
}

type Config struct {
	Postgres  PostgresConfig `mapstructure:"postgres"`
	Port      string         `mapstructure:"port"`
	JwtSecret string         `mapstructure:"jwt_secret"`
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("USER_SERVICE")

	cfg := &Config{}

	cfg.Postgres.User = getEnv("POSTGRES_USER", "")
	cfg.Postgres.Password = getEnv("POSTGRES_PASSWORD", "")
	cfg.Postgres.Db = getEnv("POSTGRES_DB", "")
	cfg.Postgres.Host = getEnv("POSTGRES_HOST", "")
	cfg.Postgres.Port = getEnv("POSTGRES_PORT", "")
	cfg.Port = getEnv("PORT", "")
	cfg.JwtSecret = getEnv("JWT_SECRET", "")

	err := validateConfig(cfg)
	if err != nil {
		return &Config{}, err
	}

	logging.Instance.Info(fmt.Sprintf("Loaded config: %+v", cfg))

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func validateConfig(cfg *Config) error {

	requiredFields := map[string]string{
		"POSTGRES_USER":     cfg.Postgres.User,
		"POSTGRES_PASSWORD": cfg.Postgres.Password,
		"POSTGRES_DB":       cfg.Postgres.Db,
		"PORT":              cfg.Port,
		"JWT_SECRET":        cfg.JwtSecret,
	}

	for field, value := range requiredFields {
		if value == "" {
			return fmt.Errorf("required environment variable %s is missing", field)
		}
	}

	return nil
}
