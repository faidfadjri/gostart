package cmd

import "github.com/spf13/cobra"

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource (e.g. handler, repository, usecase)",
}

func init() {
	CreateCmd.AddCommand(HandlerCmd)
	CreateCmd.AddCommand(UsecaseCmd)
	CreateCmd.AddCommand(RepositoryCmd)
	CreateCmd.AddCommand(FeatureCmd)
}
