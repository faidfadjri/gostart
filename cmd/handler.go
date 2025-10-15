package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/faidfadjri/gostart/cmd/types"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed templates/handler.tmpl
var handlerTmpl string

var HandlerCmd = &cobra.Command{
	Use:   "handler [name]",
	Short: "Create a new handler",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Convert name to lowercase and Title case
		name := strings.ToLower(args[0])
		caser := cases.Title(language.English)

		// Handle nested path like "task/comment"
		parts := strings.Split(name, "/")
		last := parts[len(parts)-1]
		serviceName := caser.String(last)

		// Destination directory
		destDir := filepath.Join("internal/interface/handlers")
		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			log.Fatalf("❌ Failed to create handler directory: %v", err)
		}

		// Parse embedded template
		tmpl, err := template.New("handler").Parse(handlerTmpl)
		if err != nil {
			log.Fatalf("❌ Failed to parse embedded template: %v", err)
		}

		// Get module name from go.mod
		moduleName, err := getModuleName()
		if err != nil {
			log.Fatalf("❌ Failed to get module name from go.mod: %v", err)
		}

		// Prepare data
		templateData := types.TemplateData{
			ServiceName:      serviceName,
			ServiceNameLower: last,
			ModuleName:       moduleName,
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, templateData); err != nil {
			log.Fatalf("❌ Failed to execute template: %v", err)
		}

		// Output file path
		outputPath := filepath.Join(destDir, fmt.Sprintf("%s_handler.go", last))
		if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
			log.Fatalf("❌ Failed to write handler file: %v", err)
		}

		fmt.Println("✅ Handler created at:", outputPath)
	},
}
