package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Encryption EncryptionConfig `mapstructure:"encryption"`
	Logging    LoggingConfig    `mapstructure:"logging"`
	Telemetry  TelemetryConfig  `mapstructure:"telemetry"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	TLSEnabled   bool          `mapstructure:"tls_enabled"`
	TLSCertPath  string        `mapstructure:"tls_cert_path"`
	TLSKeyPath   string        `mapstructure:"tls_key_path"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	DSN               string        `mapstructure:"dsn"`
	MaxConns          int32         `mapstructure:"max_conns"`
	MinConns          int32         `mapstructure:"min_conns"`
	MaxConnLifetime   time.Duration `mapstructure:"max_conn_lifetime"`
	HealthCheckPeriod time.Duration `mapstructure:"health_check_period"`
}

type RedisConfig struct {
	Addr       string `mapstructure:"addr"`
	Password   string `mapstructure:"password"`
	TLSEnabled bool   `mapstructure:"tls_enabled"`
	TLSCAPath  string `mapstructure:"tls_ca_path"`
	PoolSize   int    `mapstructure:"pool_size"`
}

type EncryptionConfig struct {
	Mode         string `mapstructure:"mode"`
	LocalKeyPath string `mapstructure:"local_key_path"`
	KMSKeyID     string `mapstructure:"kms_key_id"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

type TelemetryConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	ServiceName string `mapstructure:"service_name"`
	Exporter    string `mapstructure:"exporter"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetEnvPrefix("BASALT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	bindEnvKeys(v)

	setDefaults(v)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshal: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func bindEnvKeys(v *viper.Viper) {
	keys := []string{
		"server.port",
		"server.tls_enabled",
		"server.tls_cert_path",
		"server.tls_key_path",
		"server.read_timeout",
		"server.write_timeout",
		"database.dsn",
		"database.max_conns",
		"database.min_conns",
		"database.max_conn_lifetime",
		"database.health_check_period",
		"redis.addr",
		"redis.password",
		"redis.tls_enabled",
		"redis.tls_ca_path",
		"redis.pool_size",
		"encryption.mode",
		"encryption.local_key_path",
		"encryption.kms_key_id",
		"logging.level",
		"logging.format",
		"logging.output",
		"telemetry.enabled",
		"telemetry.service_name",
		"telemetry.exporter",
	}
	for _, key := range keys {
		_ = v.BindEnv(key)
	}
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.tls_enabled", false)
	v.SetDefault("server.read_timeout", 30*time.Second)
	v.SetDefault("server.write_timeout", 30*time.Second)

	v.SetDefault("database.max_conns", int32(25))
	v.SetDefault("database.min_conns", int32(5))
	v.SetDefault("database.max_conn_lifetime", 30*time.Minute)
	v.SetDefault("database.health_check_period", 30*time.Second)

	v.SetDefault("redis.tls_enabled", false)
	v.SetDefault("redis.pool_size", 10)

	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
	v.SetDefault("logging.output", "stdout")

	v.SetDefault("telemetry.enabled", true)
	v.SetDefault("telemetry.exporter", "stdout")
}

func (c *Config) Validate() error {
	var errs []string

	if strings.TrimSpace(c.Database.DSN) == "" {
		errs = append(errs, "database.dsn is required (set BASALT_DATABASE_DSN)")
	}

	if strings.TrimSpace(c.Redis.Addr) == "" {
		errs = append(errs, "redis.addr is required (set BASALT_REDIS_ADDR)")
	}

	if strings.TrimSpace(c.Encryption.Mode) == "" {
		errs = append(errs, "encryption.mode is required (set BASALT_ENCRYPTION_MODE to local or kms)")
	} else {
		switch strings.ToLower(c.Encryption.Mode) {
		case "local":
			if strings.TrimSpace(c.Encryption.LocalKeyPath) == "" {
				errs = append(errs, "encryption.local_key_path is required when encryption.mode is local")
			}
		case "kms":
			if strings.TrimSpace(c.Encryption.KMSKeyID) == "" {
				errs = append(errs, "encryption.kms_key_id is required when encryption.mode is kms")
			}
		default:
			errs = append(errs, "encryption.mode must be local or kms")
		}
	}

	if c.Server.TLSEnabled {
		if strings.TrimSpace(c.Server.TLSCertPath) == "" {
			errs = append(errs, "server.tls_cert_path is required when server.tls_enabled is true")
		}
		if strings.TrimSpace(c.Server.TLSKeyPath) == "" {
			errs = append(errs, "server.tls_key_path is required when server.tls_enabled is true")
		}
	}

	if c.Redis.TLSEnabled && strings.TrimSpace(c.Redis.TLSCAPath) == "" {
		errs = append(errs, "redis.tls_ca_path is required when redis.tls_enabled is true")
	}

	if len(errs) > 0 {
		return fmt.Errorf("config validation failed:\n  - %s", strings.Join(errs, "\n  - "))
	}

	return nil
}
