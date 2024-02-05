package config

import (
	"errors"
	"os"
	"time"

	"github.com/creasty/defaults"
	"github.com/goccy/go-yaml"
)

const (
	DevEnv  = "dev"
	ProdEnv = "prod"
)

type Config struct {
	Port int    `yaml:"port" default:"8081"`
	Env  string `yaml:"env" default:"prod"`
	DB   struct {
		DSN          string        `yaml:"dsn"`
		MaxOpenConns int           `yaml:"max_open_conns" default:"50"`
		MaxIdleConns int           `yaml:"max_idle_conns" default:"50"`
		MaxIdleTime  time.Duration `yaml:"max_idle_time" default:"1m"`
	}
}

func (c *Config) ReadFrom(filePath string) error {
	if c == nil {
		return errors.New("nil pointer config")
	}
	err := defaults.Set(c)
	if err != nil {
		return err
	}

	// _, err = os.Stat(filePath)

	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	err = yaml.NewDecoder(f).Decode(c)
	if err != nil {
		return err
	}

	return nil
}
