package config

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Server ServerConfig
	JWT    JWTConfig
	DB     DBConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

type DBConfig struct {
	Driver   string // "postgres" or "sqlite"
	Host     string // for postgres
	Port     string // for postgres
	User     string // for postgres
	Password string // for postgres
	DBName   string // for postgres
	SQLitePath string // for sqlite
}

func Load() (*Config, error) {
	// Try to load .env, but don't fail if it doesn't exist
	_ = godotenv.Load()

	dbDriver := getEnv("DB_DRIVER", "sqlite")
	
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		},
		DB: DBConfig{
			Driver: dbDriver,
		},
	}

	if dbDriver == "postgres" {
		cfg.DB.Host = getEnv("DB_HOST", "localhost")
		cfg.DB.Port = getEnv("DB_PORT", "5432")
		cfg.DB.User = getEnv("DB_USER", "postgres")
		cfg.DB.Password = getEnv("DB_PASSWORD", "")
		cfg.DB.DBName = getEnv("DB_NAME", "kotoba")
	} else {
		// SQLite
		sqlitePath := getEnv("SQLITE_PATH", "./kotoba.db")
		// Ensure directory exists
		dir := filepath.Dir(sqlitePath)
		if dir != "." {
			_ = os.MkdirAll(dir, 0755)
		}
		cfg.DB.SQLitePath = sqlitePath
	}

	return cfg, nil
}

func (c *Config) GetDatabaseDSN() string {
	if c.DB.Driver == "postgres" {
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.DBName)
	}
	// SQLite
	return c.DB.SQLitePath
}

func (c *Config) GetDB() (*sql.DB, error) {
	if c.DB.Driver == "postgres" {
		return sql.Open("postgres", c.GetDatabaseDSN())
	}
	return sql.Open("sqlite3", c.GetDatabaseDSN())
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	value := 0
	_, err := fmt.Sscanf(valueStr, "%d", &value)
	if err != nil {
		return defaultValue
	}
	return value
}
