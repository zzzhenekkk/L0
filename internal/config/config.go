package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Nats   NatsConfig
	DB     DBConfig
	Server ServerConfig
}

type NatsConfig struct {
	Cluster        string
	URL            string
	Client         string
	ClientNotifier string
	Subject        string
	DurableName    string
}

type DBConfig struct {
	DBName   string
	Host     string
	User     string
	Password string
	Port     int
	SSLMode  string
}

type ServerConfig struct {
	Host string
	Port int
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode config, %s", err)
		return nil, err
	}

	return &config, nil
}
