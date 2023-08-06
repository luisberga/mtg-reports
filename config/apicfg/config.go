package apicfg

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Database Database
	Api      Api
	LogLevel string
}

type Database struct {
	User       string
	Password   string
	Host       string
	Port       string
	Database   string
	CommitSize int
}

type Api struct {
	Port string
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

	viper.SetDefault("api.db.user", "root")
	viper.SetDefault("api.db.password", "root")
	viper.SetDefault("api.db.host", "localhost")
	viper.SetDefault("api.db.port", "3306")
	viper.SetDefault("api.db.database", "mydatabase")
	viper.SetDefault("api.db.commitSize", 1000)

	viper.SetDefault("api.port", "8088")

	viper.SetDefault("api.log.level", "debug")

	user := viper.GetString("api.db.user")
	password := viper.GetString("api.db.password")
	host := viper.GetString("api.db.host")
	dbPort := viper.GetString("api.db.port")
	database := viper.GetString("api.db.database")
	commitSize := viper.GetInt("api.db.commitSize")

	apiPort := viper.GetString("api.port")

	logLevel := viper.GetString("api.log.level")

	return &Config{
		Database: Database{
			User:       user,
			Password:   password,
			Host:       host,
			Port:       dbPort,
			Database:   database,
			CommitSize: commitSize,
		},
		Api: Api{
			Port: apiPort,
		},
		LogLevel: logLevel,
	}, nil
}
