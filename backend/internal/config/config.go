package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr            string
	DSN             string
	AdminUser       string
	AdminPassword   string
	CookieSecure    bool
	SessionLifetime time.Duration
}

func Load() Config {
	host := env("DB_HOST", "127.0.0.1")
	port := env("DB_PORT", "3306")
	user := env("DB_USER", "root")
	password := env("DB_PASSWORD", "rootpassword")
	database := env("DB_NAME", "ntocean_feast")
	secure, _ := strconv.ParseBool(env("ADMIN_COOKIE_SECURE", "false"))

	return Config{
		Addr:            env("APP_ADDR", ":8081"),
		DSN:             fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&time_zone=%%27%%2B08%%3A00%%27&multiStatements=true", user, password, host, port, database),
		AdminUser:       env("ADMIN_USER", "admin"),
		AdminPassword:   env("ADMIN_PASSWORD", "admin"),
		CookieSecure:    secure,
		SessionLifetime: 12 * time.Hour,
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
