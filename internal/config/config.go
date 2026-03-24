package config

import (
	"flag"
)

type Config struct {
	Port int
	DB   struct {
		DSN          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
}

func Load() Config {
	var cfg Config

	flag.IntVar(&cfg.Port, "port", 8000, "server port")
	flag.StringVar(
		&cfg.DB.DSN,
		"db-dsn",
		"postgres://reddit:1234@localhost/reddit?sslmode=disable",
		"PostgreSQL DSN",
	)
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(
		&cfg.DB.MaxIdleTime,
		"db-mDB-idle-time",
		"15m",
		"PostgreSQL max connection idle time",
	)

	flag.Parse()
	return cfg
}
