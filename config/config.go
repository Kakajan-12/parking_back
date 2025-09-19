package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// getProjectPaths resolves paths relative to this file
func getProjectPaths() (baseDir string, projectDir string) {
	// Get current file path
	_, filename, _, _ := runtime.Caller(0)

	// Project dir: parent of parent (2 levels up)
	projectDir = filepath.Dir(filepath.Dir(filename))

	// Base dir: parent of project dir (3 levels up)
	baseDir = filepath.Dir(projectDir)

	return baseDir, projectDir
}

type Config struct {
	// Backend
	BackendHost string
	BackendPort int

	// Database
	DatabaseURL string

	// External services
	MacroscopURL      string
	MacroscopUsername string
	MacroscopPassword string

	// Storage paths
	StaticURL string
	StaticDir string
	MediaURL  string
	MediaDir  string

	// CORS
	CORSAllowOrigins []string

	// Initial superuser
	FirstSuperUsername     string
	FirstSuperuserPassword string

	// Security
	JWTPrivateKeyPath string
	JWTPublicKeyPath  string
	JwtAudience       string
	JwtIssuer         string
	// Expiration in minutes
	JwtExpiration time.Duration

	DefaultCurrency string
}

// Global config accessible anywhere
var AppConfig *Config

func LoadConfig() {
	// Load .env file
	_, projectDir := getProjectPaths() // e.g., src/backend
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = ".env.production"
	}
	var finalPath = filepath.Join(projectDir, envPath)
	err := godotenv.Load(finalPath)
	if err != nil {
		log.Fatalf("Failed to load environment file: %v", err)
	}

	backendPort, err := strconv.Atoi(getEnv("BACKEND_PORT", "8000"))
	if err != nil {
		log.Fatalf("Invalid BACKEND_PORT: %v", err)
	}

	privatePath := filepath.Join(projectDir, "security/private.pem")
	publicPath := filepath.Join(projectDir, "security/public.pem")

	AppConfig = &Config{
		// Backend
		BackendHost: getEnv("BACKEND_HOST", "127.0.0.1"),
		BackendPort: backendPort,

		// External services
		DatabaseURL: getEnv("DATABASE_URL", "host=127.0.0.1 user=change_this password=change_this dbname=change_this port=5432 sslmode=disable"),

		// Storage paths with defaults
		MacroscopURL:      getEnv("MACROSCOP_URL", "127.0.0.1:8080"),
		MacroscopUsername: getEnv("MACROSCOP_USERNAME", "admin"),
		MacroscopPassword: getEnv("MACROSCOP_PASSWORD", "admin"),

		StaticURL: getEnv("STATIC_URL", "/static"),
		StaticDir: getEnv("STATIC_DIR", "static"),
		MediaURL:  getEnv("MEDIA_URL", "/media"),
		MediaDir:  getEnv("MEDIA_DIR", "media"),

		// Security with defaults
		JWTPrivateKeyPath: getEnv("JWT_PRIVATE_KEY_PATH", privatePath),
		JWTPublicKeyPath:  getEnv("JWT_PUBLIC_KEY_PATH", publicPath),
		JwtAudience:       "client",
		JwtIssuer:         "backend",
		JwtExpiration:     24 * time.Hour,
		// CORS
		CORSAllowOrigins: parseCSV(getEnv("CORS_ALLOW_ORIGINS", "http://localhost:8000,http://127.0.0.1:8000")),

		// Initial superuser
		FirstSuperUsername:     getEnv("FIRST_SUPER_USERNAME", "admin"),
		FirstSuperuserPassword: getEnv("FIRST_SUPERUSER_PASSWORD", "password"),

		DefaultCurrency: "TMT",
	}
}

// Helper to get env variable or fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Parse comma-separated string to slice
func parseCSV(input string) []string {
	if input == "" {
		return []string{}
	}
	parts := strings.Split(input, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
