package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var global *Config

type Config struct {
	Env         string `yaml:"env"`
	ServiceName string `yaml:"service_name"`

	Database struct {
		Host           string `yaml:"host"`
		Port           string `yaml:"port"`
		User           string `yaml:"user"`
		Pass           string `yaml:"password"`
		Name           string `yaml:"name"`
		SSLMode        string `yaml:"sslmode"`
		MigrationsPath string `yaml:"migrations_path"`
	} `yaml:"database"`

	JWT struct {
		Issuer      string `yaml:"issuer"`
		Audience    string `yaml:"audience"`
		Ed25519Seed string `yaml:"ed25519_seed"`
		KeyID       string `yaml:"key_id"`
	} `yaml:"jwt"`

	AccessTokenTTL  string `yaml:"access_token_ttl"`  // используем как строку (e.g. "15m")
	RefreshTokenTTL string `yaml:"refresh_token_ttl"` // (e.g. "720h")

	NATS struct {
		URL string `yaml:"url"`
	} `yaml:"nats"`

	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	var config Config
	if err := yaml.Unmarshal(b, &config); err != nil {
		return nil, fmt.Errorf("unmarshal %s: %w", path, err)
	}

	return &config, nil
}

func MustLoad(path string) {
	config, err := Load(path)
	if err != nil {
		panic(err)
	}

	global = config
}

func G() *Config {
	if global == nil {
		panic("global config is nil")
	}

	return global
}
