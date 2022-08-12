package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type Config struct {
	User     string `envconfig:"db_user" required:"true"`
	Password string `envconfig:"db_pass" required:"true"`
	Host     string `envconfig:"db_host" required:"true"`
	Port     uint   `envconfig:"db_port" required:"true"`
	DBName   string `envconfig:"db_schema" required:"true"`
}

func NewConfig() *Config {
	c := Config{}
	if err := envconfig.Process("", &c); err != nil {
		zap.L().Fatal("error getting config", zap.Error(err))
	}
	return &c
}

// DBConnectionString constructs URL to connect to DB
func (c *Config) DBConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s search_path=brownbag sslmode=disable", c.Host, c.Port, c.User, c.DBName, c.Password)
}
