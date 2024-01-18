package config

import "github.com/spf13/viper"

type Server struct {
	Host       string `mapstructure:"host" json:"host" yaml:"host"`
	Port       int    `mapstructure:"port" json:"port" yaml:"port"`
	HealthPath string `mapstructure:"health_path" json:"health_path" yaml:"health_path"`
}

type Config struct {
	Port                       int       `mapstructure:"port" json:"port" yaml:"port"`
	LoadBalanceStrategy        string    `mapstructure:"load_balance_strategy" json:"load_balance_strategy" yaml:"load_balance_strategy"`
	HealthCheckIntervalSeconds int       `mapstructure:"health_check_interval_seconds" json:"health_check_interval_seconds" yaml:"health_check_interval_seconds"`
	Servers                    *[]Server `mapstructure:"servers" json:"servers" yaml:"servers"`
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
