package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	App      AppConfig
	JWT      JWTConfig
	Upload   UploadConfig
}

type DatabaseConfig struct {
	Host      string
	Port      string
	User      string
	Password  string
	Name      string
	Charset   string
	ParseTime bool
	Loc       string
}

type AppConfig struct {
	Port string
	Env  string
}

type JWTConfig struct {
	Secret             string
	ExpireHours        int
	RefreshExpireHours int
}

type UploadConfig struct {
	Path        string
	MaxFileSize int64
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	expireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	refreshExpireHours, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRE_HOURS", "168"))
	maxFileSize, _ := strconv.ParseInt(getEnv("MAX_FILE_SIZE", "5242880"), 10, 64)
	parseTime, _ := strconv.ParseBool(getEnv("DB_PARSE_TIME", "true"))

	return &Config{
		Database: DatabaseConfig{
			Host:      getEnv("DB_HOST", "localhost"),
			Port:      getEnv("DB_PORT", "3306"),
			User:      getEnv("DB_USER", "root"),
			Password:  getEnv("DB_PASSWORD", ""),
			Name:      getEnv("DB_NAME", "go_commerce"),
			Charset:   getEnv("DB_CHARSET", "utf8mb4"),
			ParseTime: parseTime,
			Loc:       getEnv("DB_LOC", "Local"),
		},
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", "your-secret-key"),
			ExpireHours:        expireHours,
			RefreshExpireHours: refreshExpireHours,
		},
		Upload: UploadConfig{
			Path:        getEnv("UPLOAD_PATH", "./uploads"),
			MaxFileSize: maxFileSize,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}