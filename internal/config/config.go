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
	Logging LoggingConfig
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

type LoggingConfig struct {
	RequestLogging bool
	IsProduction   bool
}

func Load() Config {
	var cfg Config

	flag.IntVar(&cfg.Port, "port", 8000, "server port")

	cfg.DB.flags()
	cfg.Limiter.flags()
	cfg.SMTP.flags()
	cfg.JWT.flags()
	cfg.Logging.flags()

	flag.Parse()
	return cfg
}

func (db *DBConfig) flags() {
	flag.StringVar(&db.DSN, "db-dsn", os.Getenv("DB_DSN"), "Postgres DSN")
	flag.IntVar(&db.MaxOpenConns, "db-max-open-conns", 25, "Postgres max open connections")
	flag.IntVar(&db.MaxIdleConns, "db-max-idle-conns", 25, "Postgres max idle connections")
	flag.StringVar(&db.MaxIdleTime, "db-max-idle-time", "15m", "Postgres max connection idle time")
}

func (l *LimiterConfig) flags() {
	flag.BoolVar(&l.Enabled, "limiter-enabled", true, "Enable rate limiter")
}

func (s *SMTPConfig) flags() {
	flag.StringVar(&s.Host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&s.Port, "smtp-port", 587, "SMTP port")
	flag.StringVar(&s.Username, "smtp-username", "17e7914803f235", "SMTP username")
	flag.StringVar(&s.Password, "smtp-password", "afef433d663144", "SMTP password")
	flag.StringVar(
		&s.Sender,
		"smtp-sender",
		"Reddit-Clone <no-reply@reddit-clone.com>",
		"SMTP sender",
	)
}

func (j *JWTConfig) flags() {
	flag.StringVar(&j.Secret, "jwt-secret", os.Getenv("JWT_SECRET"), "JWT Secret")
	j.RefreshTokenTTL = 30 * 24 * time.Hour
	j.AccessTokenTTL = 30 * time.Minute
}

func (l *LoggingConfig) flags() {
	flag.BoolVar(&l.RequestLogging, "log-requests", true, "enable HTTP request logging")
	flag.BoolVar(&l.IsProduction, "log-prod", false, "switch output format to production JSON")
}
