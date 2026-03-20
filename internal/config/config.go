package config

import (
	"flag"
)

type Config struct {
	Port int
}

func Load() Config {
	var cfg Config

	flag.IntVar(&cfg.Port, "port", 8000, "server port")

	flag.Parse()
	return cfg
}
