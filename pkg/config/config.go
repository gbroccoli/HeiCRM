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

	Jwt struct {
		Alg       string `yaml:"alg"`
		SecretKey string `yaml:"secret_key"`
	} `yaml:"jwt"`

	Cookie struct {
		Domain string `yaml:"domain"` // e.g., ".yourdomain.com" for cross-subdomain cookies
	} `yaml:"cookie"`

	Database struct {
		Host           string `yaml:"host"`
		Port           string `yaml:"port"`
		User           string `yaml:"user"`
		Pass           string `yaml:"password"`
		Name           string `yaml:"name"`
		SSLMode        string `yaml:"sslmode"`
		MigrationsPath string `yaml:"migrations_path"`
	} `yaml:"database"`

	AccessTokenTTL  string `yaml:"access_token_ttl"`  // используем как строку (e.g. "15m")
	RefreshTokenTTL string `yaml:"refresh_token_ttl"` // (e.g. "720h")

	NATS struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"nats"`

	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
	} `yaml:"redis"`

	Serves struct {
		Auth         string `yaml:"auth"`
		Users        string `yaml:"users"`
		Tickets      string `yaml:"tickets"`
		Notification string `yaml:"notification"`
	} `yaml:"serves"`
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
