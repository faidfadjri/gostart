package cmd

import (
	"bytes"
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

var HandlerCmd = &cobra.Command{
	Use:   "handler [name]",
	Short: "Create a new handler",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Convert name
		name := strings.ToLower(args[0])
		caser := cases.Title(language.English)
		serviceName := caser.String(name)

		// Template path
		tmplPath := "templates/handler.tmpl"

		// Destination directory
		destDir := fmt.Sprintf("src/interface/handlers/%s", name)
		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			log.Fatalf("❌ Failed to create handlers directory: %v", err)
		}

		// Read and parse template
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			log.Fatalf("❌ Failed to parse template: %v", err)
		}

		var buf bytes.Buffer

		templateData := types.TemplateData{
			ServiceName:      serviceName,
			ServiceNameLower: name,
		}

		err = tmpl.Execute(&buf, templateData)
		if err != nil {
			log.Fatalf("❌ Failed to execute template: %v", err)
		}

		// Write to file
		outputPath := filepath.Join(destDir, fmt.Sprintf("%s_handler.go", name))
		err = os.WriteFile(outputPath, buf.Bytes(), 0644)
		if err != nil {
			log.Fatalf("❌ Failed to write file: %v", err)
		}

		fmt.Println("✅ Handler created at:", outputPath)
	},
}
