package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/faidfadjri/gostart/cmd/types"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var UsecaseCmd = &cobra.Command{
	Use:   "usecase [name]",
	Short: "Create a new usecase",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Convert name
		name := strings.ToLower(args[0])
		caser := cases.Title(language.English)
		serviceName := caser.String(name)

		// Template path
		tmplPath := "templates/usecase.tmpl"

		// Destination directory
		destDir := fmt.Sprintf("src/app/usecases/%s", name)
		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			log.Fatalf("❌ Failed to create usecases directory: %v", err)
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
		outputPath := filepath.Join(destDir, fmt.Sprintf("%s_usecase.go", name))
		err = os.WriteFile(outputPath, buf.Bytes(), 0644)
		if err != nil {
			log.Fatalf("❌ Failed to write file: %v", err)
		}

		fmt.Println("✅ Usecase created at:", outputPath)

		interfaceTmplPath := "templates/usecase_interface.tmpl"
		interfaceTmpl, err := template.ParseFiles(interfaceTmplPath)
		if err != nil {
			log.Fatalf("❌ Failed to parse interface template: %v", err)
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

		// Create or update usecases.go index file
		err = createOrUpdateUsecasesIndex(serviceName, name)
		if err != nil {
			log.Fatalf("❌ Failed to create/update usecases.go: %v", err)
		}

		fmt.Println("✅ Usecases index updated at: src/app/usecases/usecases.go")
	},
}

func createOrUpdateUsecasesIndex(serviceName, name string) error {
	// Get module name from go.mod
	moduleName, err := getModuleName()
	if err != nil {
		return fmt.Errorf("failed to get module name: %w", err)
	}

	usecasesPath := "src/app/usecases/usecases.go"

	// Check if usecases.go exists
	if _, err := os.Stat(usecasesPath); os.IsNotExist(err) {
		// Create new usecases.go
		return createNewUsecasesIndex(usecasesPath, moduleName, serviceName, name)
	}

	// Update existing usecases.go
	return updateExistingUsecasesIndex(usecasesPath, moduleName, serviceName, name)
}

func getModuleName() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimPrefix(line, "module "), nil
		}
	}

	return "", fmt.Errorf("module name not found in go.mod")
}

func createNewUsecasesIndex(path, moduleName, serviceName, name string) error {
	content := fmt.Sprintf(`package usecases

import "%s/src/app/usecases/%s"

var (
	New%sUsecase = %s.New%sUsecase
)
`, moduleName, name, serviceName, name, serviceName)

	// Format the Go code
	formattedContent, err := format.Source([]byte(content))
	if err != nil {
		// If formatting fails, use original content
		formattedContent = []byte(content)
	}

	return os.WriteFile(path, formattedContent, 0644)
}

func updateExistingUsecasesIndex(path, moduleName, serviceName, name string) error {
	// Read existing file
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	contentStr := string(content)

	// Check if this usecase already exists
	existsPattern := fmt.Sprintf(`New%sUsecase\s*=`, serviceName)
	exists, _ := regexp.MatchString(existsPattern, contentStr)
	if exists {
		// Usecase already exists, no need to update
		return nil
	}

	// Add new import
	newImport := fmt.Sprintf(`"%s/src/app/usecases/%s"`, moduleName, name)

	// Find import section and add new import
	importRegex := regexp.MustCompile(`(import\s+(?:\([^)]*\)|"[^"]*"))`)
	if strings.Contains(contentStr, "import (") {
		// Multi-line import
		importEndRegex := regexp.MustCompile(`(\s*)\)`)
		contentStr = importEndRegex.ReplaceAllString(contentStr, fmt.Sprintf("$1%s\n$1)", newImport))
	} else {
		// Single import, convert to multi-line
		contentStr = importRegex.ReplaceAllString(contentStr, fmt.Sprintf("import (\n\t$1\n\t%s\n)", newImport))
		// Clean up the old import format
		contentStr = strings.ReplaceAll(contentStr, `import (
	import `, "import (\n\t")
	}

	// Add new variable declaration
	newVar := fmt.Sprintf("\tNew%sUsecase = %s.New%sUsecase", serviceName, name, serviceName)

	// Find var section and add new variable
	varRegex := regexp.MustCompile(`(\s*)\)(\s*)$`)
	contentStr = varRegex.ReplaceAllString(contentStr, fmt.Sprintf("$1%s\n$1)$2", newVar))

	// Format the Go code
	formattedContent, err := format.Source([]byte(contentStr))
	if err != nil {
		// If formatting fails, use original content
		formattedContent = []byte(contentStr)
	}

	return os.WriteFile(path, formattedContent, 0644)
}
