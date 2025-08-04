package cmd

import (
	"bufio"
	"bytes"
	_ "embed"
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

		templateData := types.TemplateData{
			ServiceName:      serviceName,
			ServiceNameLower: name,
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

var (
	New%sUsecase = %s.New%sUsecase
)
`, moduleName, name, serviceName, name, serviceName)

	formattedContent, err := format.Source([]byte(content))
	if err != nil {
		formattedContent = []byte(content)
	}
	return os.WriteFile(path, formattedContent, 0644)
}

func updateExistingUsecasesIndex(path, moduleName, serviceName, name string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	contentStr := string(content)

	existsPattern := fmt.Sprintf(`New%sUsecase\s*=`, serviceName)
	exists, _ := regexp.MatchString(existsPattern, contentStr)
	if exists {
		return nil
	}

	newImport := fmt.Sprintf(`"%s/src/app/usecases/%s"`, moduleName, name)

	if strings.Contains(contentStr, "import (") {
		contentStr = regexp.MustCompile(`(?m)^import \(`).ReplaceAllString(contentStr, "import (\n\t"+newImport)
	} else {
		contentStr = regexp.MustCompile(`(?m)^import\s+"[^"]+"`).ReplaceAllString(contentStr, "import (\n\t$0\n\t"+newImport+"\n)")
	}

	newVar := fmt.Sprintf("\tNew%sUsecase = %s.New%sUsecase", serviceName, name, serviceName)
	contentStr = regexp.MustCompile(`(?m)var \(([^)]+)\)`).ReplaceAllString(contentStr, "var (\n$1\n"+newVar+"\n)")

	formattedContent, err := format.Source([]byte(contentStr))
	if err != nil {
		formattedContent = []byte(contentStr)
	}
	return os.WriteFile(path, formattedContent, 0644)
}
