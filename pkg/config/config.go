package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	Host     string
	Port     int
	Database string
}

func (c Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func Read() (cfg Config, err error) {
	_, err = toml.DecodeFile("config.toml", &cfg)
	if err != nil {
		return Config{}, err
	}

	if cfg.Host == "" {
		cfg.Host = "0.0.0.0"
		log.Println("The config host was empty! Server default host was set (0.0.0.0)")
	}
	if cfg.Port <= 0 {
		cfg.Port = 8000
		log.Println("The config port was empty! Server default port was set (8000)")
	}
	if cfg.Database == "" {
		cfg.Port = 8000
		log.Println("The config database path was empty! Server default path was set (db.sqlite3)")
	}

	return cfg, nil
}
