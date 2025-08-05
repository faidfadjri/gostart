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
	Short: "Initialize Go project structure",
	Long:  "Scaffold a standard Go project folder structure with essential directories and boilerplate files.",
	Run:   runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	fmt.Println("🚀 Initializing Go project structure...")

	createFolders()
	moduleName := resolveModuleName()

	data := types.TemplateData{
		ServiceName:      "",
		ServiceNameLower: "",
		ModuleName:       moduleName,
	}

	generateTemplateFiles(data)
	printNextSteps()
}

func createFolders() {
	dirs := []string{
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

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatalf("❌ Failed to create directory %s: %v", dir, err)
		}
		fmt.Printf("📁 Created directory: %s\n", dir)
	}
}

func resolveModuleName() string {
	moduleName, err := getModuleName()
	if err != nil {
		log.Printf("⚠️ Could not determine module name: %v. Using placeholder.", err)
		return "your-project-name"
	}
	return moduleName
}

func generateTemplateFiles(data types.TemplateData) {
	files := map[string]string{
		"src/cmd/main.go":                             "templates/main.tmpl",
		"src/app/bootstrap/bootstrap.go":              "templates/bootstrap.tmpl",
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

	for outPath, tmplPath := range files {
		if err := renderTemplate(outPath, tmplPath, data); err != nil {
			log.Fatalf("❌ Failed to generate %s: %v", outPath, err)
		}
		fmt.Printf("✅ Generated: %s\n", outPath)
	}
}

func renderTemplate(outputPath, templatePath string, data types.TemplateData) error {
	tmplBytes, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(tmplBytes))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", outputPath, err)
	}

	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", outputPath, err)
	}

	return nil
}

func printNextSteps() {
	fmt.Println("\n🎉 Project structure initialized successfully!")
	fmt.Println("📌 Next steps:")
	fmt.Println("1. Copy `.env.example` to `.env` and configure it")
	fmt.Println("2. Run: `go mod init your-project-name`")
	fmt.Println("3. Run: `go mod tidy`")
	fmt.Println("4. Run: `go run src/main.go`")
	fmt.Println("🛠 Generate components with:")
	fmt.Println("  gostart create usecase <name>")
	fmt.Println("  gostart create repository <name>")
}
