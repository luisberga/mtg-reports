package rjobcfg

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database Database
	Job      Job
	Email    Email
	LogLevel string
}

type Database struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

type Job struct {
	Timeout time.Duration
}

type Email struct {
	Host     string
	Username string
	Password string
	To       string
	Port     string
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

	viper.SetDefault("reportjob.db.user", "root")
	viper.SetDefault("reportjob.db.password", "root")
	viper.SetDefault("reportjob.db.host", "localhost")
	viper.SetDefault("reportjob.db.port", "3306")
	viper.SetDefault("reportjob.db.database", "mydatabase")

	viper.SetDefault("reportjob.timeout", "10s")

	viper.SetDefault("reportjob.log.level", "debug")

	user := viper.GetString("reportjob.db.user")
	password := viper.GetString("reportjob.db.password")
	host := viper.GetString("reportjob.db.host")
	dbPort := viper.GetString("reportjob.db.port")
	database := viper.GetString("reportjob.db.database")

	emailHost := viper.GetString("reportjob.email.host")
	emailUser := viper.GetString("reportjob.email.username")
	emailPassword := viper.GetString("reportjob.email.password")
	emailTo := viper.GetString("reportjob.email.to")
	emailPort := viper.GetString("reportjob.email.port")

	timeoutStr := viper.GetString("reportjob.timeout")

	logLevel := viper.GetString("reportjob.log.level")

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return nil, fmt.Errorf("Error parsing duration, %w", err)
	}

	return &Config{
		Database: Database{
			User:     user,
			Password: password,
			Host:     host,
			Port:     dbPort,
			Database: database,
		},
		Job: Job{
			Timeout: timeout,
		},
		LogLevel: logLevel,
		Email: Email{
			Host:     emailHost,
			Username: emailUser,
			Password: emailPassword,
			To:       emailTo,
			Port:     emailPort,
		},
	}, nil
}
