package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/faidfadjri/gostart/cmd/types"
	"github.com/spf13/cobra"
)

//go:embed templates/dockerfile.tmpl
var dockerfileTemplate string

//go:embed templates/docker_compose.tmpl
var dockerComposeTemplate string

var DockerCmd = &cobra.Command{
	Use:   "docker [name]",
	Short: "Generate Dockerfile and docker-compose.yml",
	Run: func(cmd *cobra.Command, args []string) {
		generateDockerfile()
		serviceName := "app"
		if len(args) > 0 && args[0] != "" {
			serviceName = args[0]
		}
		generateDockerCompose(serviceName)
	},
}

func generateDockerfile() {
	tmpl, err := template.New("Dockerfile").Parse(dockerfileTemplate)
	if err != nil {
		log.Fatalf("❌ Failed to parse Dockerfile template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		log.Fatalf("❌ Failed to execute Dockerfile template: %v", err)
	}

	dst := "Dockerfile"
	if err := os.WriteFile(dst, buf.Bytes(), 0644); err != nil {
		log.Fatalf("❌ Error generating Dockerfile: %v", err)
	}
	fmt.Println("✅ Dockerfile generated successfully.")
}

func generateDockerCompose(serviceName string) {
	tmpl, err := template.New("docker-compose").Parse(dockerComposeTemplate)
	if err != nil {
		log.Fatalf("❌ Failed to parse docker-compose template: %v", err)
	}

	data := types.TemplateData{
		ServiceName:      serviceName,
		ServiceNameLower: strings.ToLower(serviceName),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatalf("❌ Failed to execute docker-compose template: %v", err)
	}

	dst := "docker-compose.yaml"
	if err := os.WriteFile(dst, buf.Bytes(), 0644); err != nil {
		log.Fatalf("❌ Error generating docker-compose.yaml: %v", err)
	}
	fmt.Println("✅ docker-compose.yaml generated successfully.")
}
