package config

import (
	"flag"
	"os"
	"time"
)

type Config struct {
	Port    int
	DB      DBConfig
	Limiter LimiterConfig
	SMTP    SMTPConfig
	JWT     JWTConfig
}

type DBConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type JWTConfig struct {
	Secret          string
	RefreshTokenTTL time.Duration
	AccessTokenTTL  time.Duration
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
}

type LimiterConfig struct {
	Enabled bool
}

func Load() Config {
	var cfg Config

	flag.IntVar(&cfg.Port, "port", 8000, "server port")
	flag.StringVar(&cfg.DB.DSN, "db-dsn", os.Getenv("JWT_SECRET"), "PostgreSQL DSN")
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(
		&cfg.DB.MaxIdleTime,
		"db-mDB-idle-time",
		"15m",
		"PostgreSQL max connection idle time",
	)
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.SMTP.Host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.SMTP.Port, "smtp-port", 587, "SMTP port")
	flag.StringVar(&cfg.SMTP.Username, "smtp-username", "17e7914803f235", "SMTP username")
	flag.StringVar(&cfg.SMTP.Password, "smtp-password", "afef433d663144", "SMTP password")
	flag.StringVar(
		&cfg.SMTP.Sender,
		"smtp-sender",
		"Reddit-Clone <no-reply@reddit-clone.com>",
		"SMTP sender",
	)

	flag.StringVar(&cfg.JWT.Secret, "jwt-secret", os.Getenv("JWT_SECRET"), "JWT Secret")
	cfg.JWT.RefreshTokenTTL = 30 * 24 * time.Hour
	cfg.JWT.AccessTokenTTL = 30 * time.Minute

	flag.Parse()
	return cfg
}
