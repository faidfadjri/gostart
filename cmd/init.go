package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/faidfadjri/gostart/cmd/types"
	"github.com/spf13/cobra"
)

//go:embed templates/main.tmpl
var mainTemplate string

//go:embed templates/db.tmpl
var dbTemplate string

//go:embed templates/gitignore.tmpl
var gitignoreTemplate string

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

		// Get module name for templates
		moduleName, err := getModuleName()
		if err != nil {
			log.Printf("⚠️ Could not get module name: %v. Using placeholder.", err)
			moduleName = "your-project-name"
		}

		templateData := types.TemplateData{
			ServiceName:      "",
			ServiceNameLower: "",
			ModuleName:       moduleName,
		}

		// Generate templated files
		templatedFiles := map[string]string{
			"src/main.go":                        mainTemplate,
			"src/infrastructure/databases/db.go": dbTemplate,
			".gitignore":                         gitignoreTemplate,
		}

		for filePath, tmplContent := range templatedFiles {
			if err := generateTemplatedFile(filePath, tmplContent, templateData); err != nil {
				log.Fatalf("❌ Failed to generate %s: %v", filePath, err)
			}
			fmt.Printf("✅ Generated file: %s\n", filePath)
		}

		// Create static files
		staticFiles := map[string]string{
			"src/app/config/config.go": `package config

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
	SSLMode  string
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
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
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
			"src/interface/response/response.go": `package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        ` + "`json:\"success\"`" + `
	Message string      ` + "`json:\"message\"`" + `
	Data    interface{} ` + "`json:\"data,omitempty\"`" + `
	Errors  any         ` + "`json:\"error,omitempty\"`" + `
}

func NewSuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(message string, err any) APIResponse {
	var errorData any
	switch e := err.(type) {
	case error:
		errorData = e.Error()
	case map[string]string, map[string]any:
		errorData = e
	default:
		errorData = e
	}
	
	return APIResponse{
		Success: false,
		Message: message,
		Errors:  errorData,
	}
}

func JSONResponse(w http.ResponseWriter, statusCode int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

// Convenience functions for common responses
func Success(w http.ResponseWriter, message string, data interface{}) {
	JSONResponse(w, http.StatusOK, NewSuccessResponse(message, data))
}

func Created(w http.ResponseWriter, message string, data interface{}) {
	JSONResponse(w, http.StatusCreated, NewSuccessResponse(message, data))
}

func BadRequest(w http.ResponseWriter, message string, err any) {
	JSONResponse(w, http.StatusBadRequest, NewErrorResponse(message, err))
}

func Unauthorized(w http.ResponseWriter, message string, err any) {
	JSONResponse(w, http.StatusUnauthorized, NewErrorResponse(message, err))
}

func Forbidden(w http.ResponseWriter, message string, err any) {
	JSONResponse(w, http.StatusForbidden, NewErrorResponse(message, err))
}

func NotFound(w http.ResponseWriter, message string, err any) {
	JSONResponse(w, http.StatusNotFound, NewErrorResponse(message, err))
}

func InternalServerError(w http.ResponseWriter, message string, err any) {
	JSONResponse(w, http.StatusInternalServerError, NewErrorResponse(message, err))
}
`,
			"src/interface/request/request.go": `package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ParseJSON parses JSON request body into the provided struct
func ParseJSON(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}
	defer r.Body.Close()
	
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}
	
	return nil
}

// GetURLParam gets URL parameter from chi router
func GetURLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// GetURLParamInt gets URL parameter as integer
func GetURLParamInt(r *http.Request, key string) (int, error) {
	param := chi.URLParam(r, key)
	if param == "" {
		return 0, fmt.Errorf("parameter %s is required", key)
	}
	
	value, err := strconv.Atoi(param)
	if err != nil {
		return 0, fmt.Errorf("parameter %s must be a valid integer", key)
	}
	
	return value, nil
}

// GetQueryParam gets query parameter from request
func GetQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// GetQueryParamInt gets query parameter as integer
func GetQueryParamInt(r *http.Request, key string) (int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return 0, fmt.Errorf("query parameter %s is required", key)
	}
	
	value, err := strconv.Atoi(param)
	if err != nil {
		return 0, fmt.Errorf("query parameter %s must be a valid integer", key)
	}
	
	return value, nil
}

// GetQueryParamWithDefault gets query parameter with default value
func GetQueryParamWithDefault(r *http.Request, key, defaultValue string) string {
	if value := r.URL.Query().Get(key); value != "" {
		return value
	}
	return defaultValue
}
`,
			"src/infrastructure/databases/models/base.go": `package models

import "time"

type BaseModel struct {
	ID        int       ` + "`json:\"id\" db:\"id\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\" db:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\" db:\"updated_at\"`" + `
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
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your_jwt_secret_key

# Other Configuration
APP_ENV=development
`,
			"README.md": `# Go Project

This project was generated using gostart.

## Project Structure

` + "```" + `
src/
├── app/
│   ├── controllers/         # Business logic controllers
│   ├── usecases/           # Application use cases
│   └── config/             # Application configuration
├── infrastructure/
│   ├── middlewares/        # HTTP middlewares
│   ├── databases/          # Database connections and models
│   │   └── models/         # Database models
│   ├── repositories/       # Data access layer
│   └── services/          # External services
└── interface/
    ├── handlers/           # HTTP request handlers
    ├── request/           # Request DTOs and parsers
    └── response/          # Response DTOs and helpers
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

		// Create static files
		for filePath, content := range staticFiles {
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

func generateTemplatedFile(filePath, tmplContent string, data types.TemplateData) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", filePath, err)
	}

	// Parse template
	tmpl, err := template.New(filepath.Base(filePath)).Parse(tmplContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
