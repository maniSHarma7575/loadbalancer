package config

import (
	"github.com/spf13/viper"
)

type Server struct {
	Host       string `mapstructure:"host" json:"host" yaml:"host"`
	Port       int    `mapstructure:"port" json:"port" yaml:"port"`
	HealthPath string `mapstructure:"health_path" json:"health_path" yaml:"health_path"`
	AppName    string `mapstructure:"app_name" json:"app_name" yaml:"app_name"`
}

type StickySession struct {
	CookieKey  string `mapstructure:"cookie_name" json:"cookie_name" yaml:"cookie_name"`
	TTLSeconds int    `mapstructure:"ttl_seconds" json:"ttl_seconds" yaml:"ttl_seconds"`
}

type Tls struct {
	Enabled  bool   `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	CertFile string `mapstructure:"cert_file" json:"cert_file" yaml:"cert_file"`
	KeyFile  string `mapstructure:"key_file" json:"key_file" yaml:"key_file"`
}

type RouteAction struct {
	RouteTo string `mapstructure:"route_to" json:"route_to" yaml:"route_to"`
}

type RouteCondition struct {
	PathPrefix string            `mapstructure:"path_prefix" json:"path_prefix" yaml:"path_prefix"`
	Headers    map[string]string `mapstructure:"headers" json:"headers" yaml:"headers"`
	Method     string            `mapstructure:"method" json:"method" yaml:"method"`
}

type RoutingRule struct {
	Conditions []RouteCondition `mapstructure:"conditions" json:"conditions" yaml:"conditions"`
	Action     RouteAction      `mapstructure:"action" json:"action" yaml:"action"`
}
type Routing struct {
	Rules []RoutingRule `mapstructure:"rules" json:"rules" yaml:"rules"`
}

type Config struct {
	Port                       int           `mapstructure:"port" json:"port" yaml:"port"`
	LoadBalanceStrategy        string        `mapstructure:"load_balance_strategy" json:"load_balance_strategy" yaml:"load_balance_strategy"`
	HealthCheckIntervalSeconds int           `mapstructure:"health_check_interval_seconds" json:"health_check_interval_seconds" yaml:"health_check_interval_seconds"`
	Servers                    *[]Server     `mapstructure:"servers" json:"servers" yaml:"servers"`
	StickySession              StickySession `mapstructure:"sticky_session" json:"sticky_session" yaml:"sticky_session"`
	Tls                        Tls           `mapstructure:"tls" json:"tls" yaml:"tls"`
	Routing                    Routing       `mapstructure:"routing" json:"routing" yaml:"routing"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	var config Config

	err = viper.Unmarshal(&config)

	if err != nil {
		return nil, err
	}

	if config.HealthCheckIntervalSeconds == 0 {
		config.HealthCheckIntervalSeconds = 5
	}

	return &config, nil
}

func ConfigPaths() []string {
	return []string{
		"internal/config/config.json",
		"internal/config/config.yaml",
	}
}
