package config

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

// DBConn represents the database connection configuration.
type DBConn struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     string `yaml:"port" mapstructure:"port"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
	DBName   string `yaml:"db_name" mapstructure:"database_name"`
}

type App struct {
	Port string `yaml:"port" mapstructure:"port"`
}

// Config represents the configuration of the application.
type Config struct {
	Conn DBConn `yaml:"database_connection" mapstructure:"database_connection"`
	App  App    `yaml:"app" mapstructure:"app"`
}

// splitFilename splits the filename into name and extension.
func splitFilename(filename string) (name string, ext string, err error) {
	lastPeriodIdx := strings.LastIndex(filename, ".")
	if lastPeriodIdx == -1 {
		return "", "", errors.New("file extension not found")
	}

	name = filename[:lastPeriodIdx]
	ext = filename[lastPeriodIdx+1:]
	return name, ext, nil
}

// ReadConfig reads the configuration from the given file.
func ReadConfig(filename string) (Config, error) {
	name, ext, err := splitFilename(filename)
	if err != nil {
		return Config{}, err
	}

	viper.SetConfigType(ext)
	viper.SetConfigName(name)
	viper.AddConfigPath("./configs/")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
