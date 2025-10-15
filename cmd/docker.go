package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var DockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Generate Dockerfile and docker-compose.yml",
}

func init() {
	DockerCmd.Run = func(cmd *cobra.Command, args []string) {
		generateDockerfile()
		generateDockerCompose()
	}
}

func generateDockerfile() {
	src := filepath.Join("cmd", "templates", "Dockerfile.tmpl")
	dst := "Dockerfile"
	if err := copyFile(src, dst); err != nil {
		fmt.Printf("Error generating Dockerfile: %v\n", err)
	} else {
		fmt.Println("Dockerfile generated successfully.")
	}
}

func generateDockerCompose() {
	src := filepath.Join("cmd", "templates", "docker_compose.tmpl")
	dst := "docker-compose.yaml"
	if err := copyFile(src, dst); err != nil {
		fmt.Printf("Error generating docker-compose.yaml: %v\n", err)
	} else {
		fmt.Println("docker-compose.yaml generated successfully.")
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Sync()
}
