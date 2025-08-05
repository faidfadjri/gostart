package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init Folder Structure",
	Long:  "Initialize the standard Go project folder structure with necessary directories and files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Initializing Go project structure...")

		// Define the folder structure
		folders := []string{
			"src/app/controllers",
			"src/app/usecases",
			"src/app/config",
			"src/infrastructure/middlewares",
			"src/infrastructure/databases/models",
			"src/infrastructure/repositories",
			"src/infrastructure/services",
			"src/interface/handlers",
			"src/interface/request",
			"src/interface/response",
		}

		// Create directories
		for _, folder := range folders {
			if err := os.MkdirAll(folder, os.ModePerm); err != nil {
				log.Fatalf("❌ Failed to create directory %s: %v", folder, err)
			}
			fmt.Printf("✅ Created directory: %s\n", folder)
		}

		// Create initial files
		files := map[string]string{
			"src/main.go": `package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	
	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Hello, World!",
		})
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
`,
			"src/config/config.go": `package config

import (
	"os"
)

type Config struct {
	Port     string
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Username: getEnv("DB_USERNAME", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
`,
			"src/pkg/responses/response.go": `package responses

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        ` + "`json:\"success\"`" + `
	Message string      ` + "`json:\"message\"`" + `
	Data    interface{} ` + "`json:\"data,omitempty\"`" + `
	Error   interface{} ` + "`json:\"error,omitempty\"`" + `
}

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func Success(w http.ResponseWriter, message string, data interface{}) {
	writeJSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, statusCode int, message string, err interface{}) {
	writeJSON(w, statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func BadRequest(w http.ResponseWriter, message string, err interface{}) {
	Error(w, http.StatusBadRequest, message, err)
}

func InternalServerError(w http.ResponseWriter, message string, err interface{}) {
	Error(w, http.StatusInternalServerError, message, err)
}

func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message, nil)
}

func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, message, nil)
}
`,
			"src/pkg/errors/errors.go": `package errors

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrUnauthorized  = errors.New("unauthorized access")
	ErrBadRequest    = errors.New("bad request")
	ErrInternalError = errors.New("internal server error")
	ErrValidation    = errors.New("validation error")
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
`,
			".env.example": `# Server Configuration
PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=your_username
DB_PASSWORD=your_password
DB_NAME=your_database

# JWT Configuration
JWT_SECRET=your_jwt_secret_key

# Other Configuration
APP_ENV=development
`,
			".gitignore": `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with "go test -c"
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
vendor/

# Go workspace file
go.work

# Environment files
.env
.env.local

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Logs
*.log

# Binary
main
app

# Temporary files
tmp/
temp/
`,
			"README.md": `# Go Project

This project was generated using gostart.

## Project Structure

` + "```" + `
src/
├── app/
│   ├── controllers/      # HTTP handlers
│   ├── middlewares/      # HTTP middlewares
│   ├── models/          # Data models
│   ├── repositories/    # Data access layer
│   ├── usecases/       # Business logic
│   └── services/       # External services
├── config/             # Configuration
├── database/          # Database related files
│   ├── migrations/    # Database migrations
│   └── seeders/      # Database seeders
└── pkg/              # Shared packages
    ├── utils/        # Utility functions
    ├── validators/   # Input validators
    ├── errors/      # Custom errors
    └── responses/   # HTTP response helpers
` + "```" + `

## Getting Started

1. Copy ` + "`.env.example`" + ` to ` + "`.env`" + ` and configure your environment variables
2. Install dependencies: ` + "`go mod tidy`" + `
3. Run the application: ` + "`go run src/main.go`" + `

## Commands

Generate new components using gostart:

` + "```bash" + `
# Generate usecase
gostart create usecase user

# Generate repository  
gostart create repository user
` + "```" + `

## License

This project is licensed under the MIT License.
`,
		}

		// Create files
		for filePath, content := range files {
			// Create directory if it doesn't exist
			dir := filepath.Dir(filePath)
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				log.Fatalf("❌ Failed to create directory for %s: %v", filePath, err)
			}

			// Write file
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				log.Fatalf("❌ Failed to create file %s: %v", filePath, err)
			}
			fmt.Printf("✅ Created file: %s\n", filePath)
		}

		fmt.Println()
		fmt.Println("🎉 Project structure initialized successfully!")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("1. Copy .env.example to .env and configure your environment")
		fmt.Println("2. Run: go mod init your-project-name")
		fmt.Println("3. Run: go mod tidy")
		fmt.Println("4. Run: go run src/main.go")
		fmt.Println()
		fmt.Println("Generate components with:")
		fmt.Println("  gostart create usecase <name>")
		fmt.Println("  gostart create repository <name>")
	},
}
