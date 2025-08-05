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

//go:embed templates/repository.tmpl
var repositoryTemplate string

//go:embed templates/repository_interface.tmpl
var repositoryInterfaceTemplate string

var RepositoryCmd = &cobra.Command{
	Use:   "repository [name]",
	Short: "Create a new repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := strings.ToLower(args[0])
		caser := cases.Title(language.English)
		serviceName := caser.String(name)

		destDir := fmt.Sprintf("internal/infrastructure/repositories/%s", name)
		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			log.Fatalf("❌ Failed to create repositories directory: %v", err)
		}

		templateData := types.TemplateData{
			ServiceName:      serviceName,
			ServiceNameLower: name,
		}

		// Parse and write repository.tmpl
		tmpl, err := template.New("repository").Parse(repositoryTemplate)
		if err != nil {
			log.Fatalf("❌ Failed to parse embedded repository template: %v", err)
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, templateData)
		if err != nil {
			log.Fatalf("❌ Failed to execute repository template: %v", err)
		}
		outputPath := filepath.Join(destDir, fmt.Sprintf("%s_repository.go", name))
		err = os.WriteFile(outputPath, buf.Bytes(), 0644)
		if err != nil {
			log.Fatalf("❌ Failed to write repository file: %v", err)
		}
		fmt.Println("✅ Repository created at:", outputPath)

		// Parse and write repository_interface.tmpl
		interfaceTmpl, err := template.New("repository_interface").Parse(repositoryInterfaceTemplate)
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

		// Update repositories.go
		err = createOrUpdateRepositoriesIndex(serviceName, name)
		if err != nil {
			log.Fatalf("❌ Failed to create/update repositories.go: %v", err)
		}
		fmt.Println("✅ Repositories index updated at: internal/infrastructure/repositories/repositories.go")
	},
}

func createOrUpdateRepositoriesIndex(serviceName, name string) error {
	moduleName, err := getModuleName()
	if err != nil {
		return fmt.Errorf("failed to get module name: %w", err)
	}

	repositoriesPath := "internal/infrastructure/repositories/repositories.go"

	if _, err := os.Stat(repositoriesPath); os.IsNotExist(err) {
		return createNewRepositoriesIndex(repositoriesPath, moduleName, serviceName, name)
	}
	return updateExistingRepositoriesIndex(repositoriesPath, moduleName, serviceName, name)
}

func createNewRepositoriesIndex(path, moduleName, serviceName, name string) error {
	content := fmt.Sprintf(`package repositories

import "%s/internal/infrastructure/repositories/%s"

type %sRepository = %s.%sRepository

var (
	New%sRepository = %s.New%sRepository
)
`, moduleName, name, serviceName, name, serviceName, serviceName, name, serviceName)

	formattedContent, err := format.Source([]byte(content))
	if err != nil {
		formattedContent = []byte(content)
	}
	return os.WriteFile(path, formattedContent, 0644)
}

// RepositoryEntry represents a single repository entry
type RepositoryEntry struct {
	Name        string // e.g., "user"
	ServiceName string // e.g., "User"
	ModuleName  string
}

func updateExistingRepositoriesIndex(path, moduleName, serviceName, name string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	contentStr := string(content)

	// Skip if already exists
	if strings.Contains(contentStr, fmt.Sprintf("New%sRepository", serviceName)) {
		return nil
	}

	// Parse existing entries
	entries, err := parseExistingRepositoryEntries(contentStr, moduleName)
	if err != nil {
		return err
	}

	// Add new entry
	entries = append(entries, RepositoryEntry{
		Name:        name,
		ServiceName: serviceName,
		ModuleName:  moduleName,
	})

	// Sort entries by service name for consistency
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ServiceName < entries[j].ServiceName
	})

	// Generate clean content
	newContent := generateCleanRepositoriesFile(entries)

	// Format the content
	formatted, err := format.Source([]byte(newContent))
	if err != nil {
		log.Println("⚠️ Failed to format repositories.go, saving raw.")
		formatted = []byte(newContent)
	}

	return os.WriteFile(path, formatted, 0644)
}

func parseExistingRepositoryEntries(content, moduleName string) ([]RepositoryEntry, error) {
	var entries []RepositoryEntry

	// Extract imports
	importRegex := regexp.MustCompile(`"` + regexp.QuoteMeta(moduleName) + `/internal/infrastructure/repositories/([^"]+)"`)
	importMatches := importRegex.FindAllStringSubmatch(content, -1)

	// Extract type aliases
	typeRegex := regexp.MustCompile(`(\w+)Repository\s*=\s*(\w+)\.(\w+)Repository`)
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
				entries = append(entries, RepositoryEntry{
					Name:        packageName,
					ServiceName: serviceName,
					ModuleName:  moduleName,
				})
			}
		}
	}

	return entries, nil
}

func generateCleanRepositoriesFile(entries []RepositoryEntry) string {
	var buf strings.Builder

	buf.WriteString("package repositories\n\n")

	// Generate imports
	if len(entries) > 0 {
		buf.WriteString("import (\n")
		for _, entry := range entries {
			buf.WriteString(fmt.Sprintf("\t\"%s/internal/infrastructure/repositories/%s\"\n", entry.ModuleName, entry.Name))
		}
		buf.WriteString(")\n\n")
	}

	// Generate type aliases
	if len(entries) > 0 {
		buf.WriteString("type (\n")
		for _, entry := range entries {
			buf.WriteString(fmt.Sprintf("\t%sRepository = %s.%sRepository\n", entry.ServiceName, entry.Name, entry.ServiceName))
		}
		buf.WriteString(")\n\n")
	}

	// Generate var block
	if len(entries) > 0 {
		buf.WriteString("var (\n")
		for _, entry := range entries {
			buf.WriteString(fmt.Sprintf("\tNew%sRepository = %s.New%sRepository\n", entry.ServiceName, entry.Name, entry.ServiceName))
		}
		buf.WriteString(")\n")
	}

	return buf.String()
}
