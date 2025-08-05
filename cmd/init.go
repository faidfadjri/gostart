package cmd

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/faidfadjri/gostart/cmd/types"
	"github.com/spf13/cobra"
)

//go:embed templates/*
var templateFS embed.FS

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init Folder Structure",
	Long:  "Initialize the standard Go project folder structure with necessary directories and files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Initializing Go project structure...")

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

		for _, folder := range folders {
			if err := os.MkdirAll(folder, os.ModePerm); err != nil {
				log.Fatalf("❌ Failed to create directory %s: %v", folder, err)
			}
			fmt.Printf("✅ Created directory: %s\n", folder)
		}

		moduleName, err := getModuleName()
		if err != nil {
			log.Printf("⚠️ Could not get module name: %v. Using placeholder.", err)
			moduleName = "your-project-name"
		}

		data := types.TemplateData{
			ServiceName:      "",
			ServiceNameLower: "",
			ModuleName:       moduleName,
		}

		templatedFiles := map[string]string{
			"src/cmd/main.go":                             "templates/main.tmpl",
			".air.toml":                                   "templates/air.tmpl",
			"src/infrastructure/databases/db.go":          "templates/db.tmpl",
			".gitignore":                                  "templates/gitignore.tmpl",
			"README.md":                                   "templates/readme.tmpl",
			".env.example":                                "templates/env.tmpl",
			"src/app/config/config.go":                    "templates/config.tmpl",
			"src/interface/response/response.go":          "templates/response.tmpl",
			"src/interface/request/request.go":            "templates/request.tmpl",
			"src/infrastructure/databases/models/user.go": "templates/models.tmpl",
		}

		for filePath, tmplPath := range templatedFiles {
			if err := generateFileFromEmbed(filePath, tmplPath, data); err != nil {
				log.Fatalf("❌ Failed to generate %s: %v", filePath, err)
			}
			fmt.Printf("✅ Generated file: %s\n", filePath)
		}

		fmt.Println("\n🎉 Project structure initialized successfully!")
		fmt.Println("\nNext steps:")
		fmt.Println("1. Copy .env.example to .env and configure your environment")
		fmt.Println("2. Run: go mod init your-project-name")
		fmt.Println("3. Run: go mod tidy")
		fmt.Println("4. Run: go run src/main.go")
		fmt.Println("\nGenerate components with:")
		fmt.Println("  gostart create usecase <name>")
		fmt.Println("  gostart create repository <name>")
	},
}

func generateFileFromEmbed(outputPath, tmplPath string, data types.TemplateData) error {
	tmplBytes, err := templateFS.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	tmpl, err := template.New(filepath.Base(tmplPath)).Parse(string(tmplBytes))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
