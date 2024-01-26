package config

import (
	"errors"
	"os"

	"github.com/goccy/go-yaml"
)

const (
	DevEnv  = "dev"
	ProdEnv = "prod"
)

type Config struct {
	Port int `yaml:"port"`
	Env  string
	DB   struct {
		DSN          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
}

func (c *Config) ReadFrom(filePath string) error {
	if c == nil {
		return errors.New("nil pointer config")
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	err = yaml.NewDecoder(f).Decode(c)
	if err != nil {
		return err
	}

	return nil
}
