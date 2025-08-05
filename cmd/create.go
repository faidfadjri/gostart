package cmd

import "github.com/spf13/cobra"

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource (e.g. service, module)",
}

func init() {
	CreateCmd.AddCommand(HandlerCmd)
	CreateCmd.AddCommand(UsecaseCmd)
}
