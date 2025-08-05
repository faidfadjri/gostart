package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var FeatureCmd = &cobra.Command{
	Use:   "feature [name]",
	Short: "Create usecase, repository, and handler, and inject to bootstrap.go",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := strings.ToLower(args[0])
		pascal := cases.Title(language.English).String(name)
		log.Println("🚀 Generating feature:", name)

		if UsecaseCmd.Run != nil {
			UsecaseCmd.Run(UsecaseCmd, args)
		}
		if RepositoryCmd.Run != nil {
			RepositoryCmd.Run(RepositoryCmd, args)
		}
		if HandlerCmd.Run != nil {
			HandlerCmd.Run(HandlerCmd, args)
		}

		// Bootstrap path
		bootstrapPath := filepath.Join("src", "app", "bootstrap", "bootstrap.go")
		if _, err := os.Stat(bootstrapPath); os.IsNotExist(err) {
			err := os.MkdirAll(filepath.Dir(bootstrapPath), os.ModePerm)
			if err != nil {
				log.Fatalf("❌ Failed to create bootstrap directory: %v", err)
			}

			initialContent := `package bootstrap

import (
	"gorm.io/gorm"
)

type Dependencies struct {
	DB *gorm.DB
}

func InitDependencies() *Dependencies {
	db, err := database.ConnectDB()

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Repositories

	// Usecases

	// Handlers

	return &Dependencies{
		DB: db,
	}
}
`
			if err := os.WriteFile(bootstrapPath, []byte(initialContent), 0644); err != nil {
				log.Fatalf("❌ Failed to create bootstrap.go: %v", err)
			}
			log.Println("📦 Created new bootstrap.go")
		}

		if err := injectToBootstrap(name, pascal); err != nil {
			log.Printf("❌ Failed to inject to bootstrap: %v", err)
		} else {
			log.Println("✅ Injected to bootstrap.go")
		}
	},
}

func injectToBootstrap(name, pascal string) error {
	module, _ := getModuleName()
	path := "src/app/bootstrap/bootstrap.go"
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(contentBytes)

	imports := map[string]string{
		`"log"`: "log",
		`"` + module + `/src/infrastructure/database"`:     "database",
		`"` + module + `/src/infrastructure/repositories"`: "repositories",
		`"` + module + `/src/interface/handler"`:           "handler",
		`"` + module + `/src/app/usecases"`:                "usecases",
	}

	// Inject missing imports
	for full, pkg := range imports {
		if !strings.Contains(content, pkg) {
			content = injectImport(content, full)
		}
	}

	// Inject Repository
	repoLine := pascal + "Repo := repositories.New" + pascal + "Repository(db)"
	if !strings.Contains(content, repoLine) {
		content = injectAfter(content, "// Repositories", "\n\t"+repoLine)
	}

	// Inject Usecase
	usecaseLine := name + "Usecase := usecases.New" + pascal + "Usecase(" + name + "Repo)"
	if !strings.Contains(content, usecaseLine) {
		content = injectAfter(content, "// Usecases", "\n\t"+usecaseLine)
	}

	// Inject Handler
	handlerLine := name + "Handler := handler.New" + pascal + "Handler(" + name + "Usecase)"
	if !strings.Contains(content, handlerLine) {
		content = injectAfter(content, "// Handlers", "\n\t"+handlerLine)
	}

	// Inject return
	returnLine := pascal + "Handler: " + name + "Handler,"
	if !strings.Contains(content, returnLine) {
		content = injectAfter(content, "return &Dependencies{", "\n\t\t"+returnLine)
	}

	// Inject Dependencies struct field
	structLine := pascal + "Handler *handler." + pascal + "Handler"
	if !strings.Contains(content, structLine) {
		content = injectAfter(content, "type Dependencies struct {", "\n\t"+structLine)
	}

	return os.WriteFile(path, []byte(content), 0644)
}

func injectImport(content, importLine string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "import (") {
			// Cari akhir dari block import
			for j := i + 1; j < len(lines); j++ {
				if strings.HasPrefix(lines[j], ")") {
					lines = append(lines[:j], append([]string{"\t" + importLine}, lines[j:]...)...)
					return strings.Join(lines, "\n")
				}
			}
		}
	}
	return content
}

func injectAfter(content, marker, toInject string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, marker) {
			indent := detectIndentation(lines, i+1)
			lines = append(lines[:i+1], append([]string{indent + strings.TrimLeft(toInject, "\n\t")}, lines[i+1:]...)...)
			break
		}
	}
	return strings.Join(lines, "\n")
}

func detectIndentation(lines []string, start int) string {
	for i := start; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) != "" {
			return line[:len(line)-len(strings.TrimLeft(line, "\t "))]
		}
	}
	return "\t"
}
