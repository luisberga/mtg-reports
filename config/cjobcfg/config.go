package cjobcfg

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database        Database
	Job             Job
	ExchangeGateway ExchangeGateway
	LogLevel        string
}

type Database struct {
	User       string
	Password   string
	Host       string
	Port       string
	Database   string
	CommitSize int
}

type Job struct {
	Timeout time.Duration
}

type ExchangeGateway struct {
	Url string
}

func New() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	viper.SetDefault("conciliatejob.db.user", "root")
	viper.SetDefault("conciliatejob.db.password", "root")
	viper.SetDefault("conciliatejob.db.host", "localhost")
	viper.SetDefault("conciliatejob.db.port", "3306")
	viper.SetDefault("conciliatejob.db.database", "mydatabase")
	viper.SetDefault("conciliatejob.db.commitSize", 1000)

	viper.SetDefault("conciliatejob.timeout", "10s")

	viper.SetDefault("conciliatejob.log.level", "debug")

	user := viper.GetString("conciliatejob.db.user")
	password := viper.GetString("conciliatejob.db.password")
	host := viper.GetString("conciliatejob.db.host")
	dbPort := viper.GetString("conciliatejob.db.port")
	database := viper.GetString("conciliatejob.db.database")
	commitSize := viper.GetInt("conciliatejob.db.commitSize")

	exchangeUrl := viper.GetString("conciliatejob.exchange.url")

	timeoutStr := viper.GetString("conciliatejob.timeout")

	logLevel := viper.GetString("conciliatejob.log.level")

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return nil, fmt.Errorf("Error parsing duration, %w", err)
	}

	return &Config{
		Database: Database{
			User:       user,
			Password:   password,
			Host:       host,
			Port:       dbPort,
			Database:   database,
			CommitSize: commitSize,
		},
		Job: Job{
			Timeout: timeout,
		},
		ExchangeGateway: ExchangeGateway{
			Url: exchangeUrl,
		},
		LogLevel: logLevel,
	}, nil
}
