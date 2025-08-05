package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/faidfadjri/gostart/cmd/types"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed templates/usecase.tmpl
var usecaseTemplate string

//go:embed templates/usecase_interface.tmpl
var interfaceTemplate string

var UsecaseCmd = &cobra.Command{
	Use:   "usecase [name]",
	Short: "Create a new usecase",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := strings.ToLower(args[0])
		caser := cases.Title(language.English)
		serviceName := caser.String(name)

		destDir := fmt.Sprintf("src/app/usecases/%s", name)
		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			log.Fatalf("❌ Failed to create usecases directory: %v", err)
		}

		moduleName, _ := getModuleName()

		templateData := types.TemplateData{
			ServiceName:      serviceName,
			ServiceNameLower: name,
			ModuleName:       moduleName,
		}

		// Parse and write usecase.tmpl
		tmpl, err := template.New("usecase").Parse(usecaseTemplate)
		if err != nil {
			log.Fatalf("❌ Failed to parse embedded usecase template: %v", err)
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, templateData)
		if err != nil {
			log.Fatalf("❌ Failed to execute usecase template: %v", err)
		}
		outputPath := filepath.Join(destDir, fmt.Sprintf("%s_usecase.go", name))
		err = os.WriteFile(outputPath, buf.Bytes(), 0644)
		if err != nil {
			log.Fatalf("❌ Failed to write usecase file: %v", err)
		}
		fmt.Println("✅ Usecase created at:", outputPath)

		// Parse and write usecase_interface.tmpl
		interfaceTmpl, err := template.New("usecase_interface").Parse(interfaceTemplate)
		if err != nil {
			log.Fatalf("❌ Failed to parse embedded interface template: %v", err)
		}
		var interfaceBuf bytes.Buffer
		err = interfaceTmpl.Execute(&interfaceBuf, templateData)
		if err != nil {
			log.Fatalf("❌ Failed to execute interface template: %v", err)
		}
		interfacePath := filepath.Join(destDir, "interface.go")
		err = os.WriteFile(interfacePath, interfaceBuf.Bytes(), 0644)
		if err != nil {
			log.Fatalf("❌ Failed to write interface.go: %v", err)
		}
		fmt.Println("✅ Interface created at:", interfacePath)

		// Update usecases.go
		err = createOrUpdateUsecasesIndex(serviceName, name)
		if err != nil {
			log.Fatalf("❌ Failed to create/update usecases.go: %v", err)
		}
		fmt.Println("✅ Usecases index updated at: src/app/usecases/usecases.go")
	},
}

func createOrUpdateUsecasesIndex(serviceName, name string) error {
	moduleName, err := getModuleName()
	if err != nil {
		return fmt.Errorf("failed to get module name: %w", err)
	}

	usecasesPath := "src/app/usecases/usecases.go"

	if _, err := os.Stat(usecasesPath); os.IsNotExist(err) {
		return createNewUsecasesIndex(usecasesPath, moduleName, serviceName, name)
	}
	return updateExistingUsecasesIndex(usecasesPath, moduleName, serviceName, name)
}

func createNewUsecasesIndex(path, moduleName, serviceName, name string) error {
	content := fmt.Sprintf(`package usecases

import "%s/src/app/usecases/%s"

type %sUsecase = %s.%sUsecase

var (
	New%sUsecase = %s.New%sUsecase
)
`, moduleName, name, serviceName, name, serviceName, serviceName, name, serviceName)

	formattedContent, err := format.Source([]byte(content))
	if err != nil {
		formattedContent = []byte(content)
	}
	return os.WriteFile(path, formattedContent, 0644)
}

// UsecaseEntry represents a single usecase entry
type UsecaseEntry struct {
	Name        string // e.g., "user"
	ServiceName string // e.g., "User"
	ModuleName  string
}

func updateExistingUsecasesIndex(path, moduleName, serviceName, name string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	contentStr := string(content)

	// Skip if already exists
	if strings.Contains(contentStr, fmt.Sprintf("New%sUsecase", serviceName)) {
		return nil
	}

	// Parse existing entries
	entries, err := parseExistingEntries(contentStr, moduleName)
	if err != nil {
		return err
	}

	// Add new entry
	entries = append(entries, UsecaseEntry{
		Name:        name,
		ServiceName: serviceName,
		ModuleName:  moduleName,
	})

	// Sort entries by service name for consistency
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ServiceName < entries[j].ServiceName
	})

	// Generate clean content
	newContent := generateCleanUsecasesFile(entries)

	// Format the content
	formatted, err := format.Source([]byte(newContent))
	if err != nil {
		log.Println("⚠️ Failed to format usecases.go, saving raw.")
		formatted = []byte(newContent)
	}

	return os.WriteFile(path, formatted, 0644)
}

func parseExistingEntries(content, moduleName string) ([]UsecaseEntry, error) {
	var entries []UsecaseEntry

	// Extract imports
	importRegex := regexp.MustCompile(`"` + regexp.QuoteMeta(moduleName) + `/src/app/usecases/([^"]+)"`)
	importMatches := importRegex.FindAllStringSubmatch(content, -1)

	// Extract type aliases
	typeRegex := regexp.MustCompile(`(\w+)Usecase\s*=\s*(\w+)\.(\w+)Usecase`)
	typeMatches := typeRegex.FindAllStringSubmatch(content, -1)

	// Create a map for easier lookup
	typeMap := make(map[string]string) // serviceName -> name
	for _, match := range typeMatches {
		if len(match) >= 4 {
			serviceName := match[1] // e.g., "User"
			packageName := match[2] // e.g., "user"
			typeMap[serviceName] = packageName
		}
	}

	// Build entries based on imports
	for _, match := range importMatches {
		if len(match) >= 2 {
			packageName := match[1] // e.g., "user"
			// Convert to service name (capitalize first letter)
			caser := cases.Title(language.English)
			serviceName := caser.String(packageName)

			// Verify this entry exists in types
			if _, exists := typeMap[serviceName]; exists {
				entries = append(entries, UsecaseEntry{
					Name:        packageName,
					ServiceName: serviceName,
					ModuleName:  moduleName,
				})
			}
		}
	}

	return entries, nil
}

func generateCleanUsecasesFile(entries []UsecaseEntry) string {
	var buf strings.Builder

	buf.WriteString("package usecases\n\n")

	// Generate imports
	if len(entries) > 0 {
		buf.WriteString("import (\n")
		for _, entry := range entries {
			buf.WriteString(fmt.Sprintf("\t\"%s/src/app/usecases/%s\"\n", entry.ModuleName, entry.Name))
		}
		buf.WriteString(")\n\n")
	}

	// Generate type aliases
	if len(entries) > 0 {
		buf.WriteString("type (\n")
		for _, entry := range entries {
			buf.WriteString(fmt.Sprintf("\t%sUsecase = %s.%sUsecase\n", entry.ServiceName, entry.Name, entry.ServiceName))
		}
		buf.WriteString(")\n\n")
	}

	// Generate var block
	if len(entries) > 0 {
		buf.WriteString("var (\n")
		for _, entry := range entries {
			buf.WriteString(fmt.Sprintf("\tNew%sUsecase = %s.New%sUsecase\n", entry.ServiceName, entry.Name, entry.ServiceName))
		}
		buf.WriteString(")\n")
	}

	return buf.String()
}
