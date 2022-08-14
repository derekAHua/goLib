package http

import "time"

type ServerConfig struct {
	Address      string        `yaml:"address"`
	ReadTimeout  time.Duration `yaml:"readTimeOut"`
	WriteTimeout time.Duration `yaml:"writeTimeOut"`
}
