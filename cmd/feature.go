package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var FeatureCmd = &cobra.Command{
	Use:   "feature [name]",
	Short: "Create usecase, repository, and handler",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := strings.ToLower(args[0])
		log.Println("🚀 Generating feature:", name)

		if UsecaseCmd.Run != nil {
			UsecaseCmd.Run(UsecaseCmd, args)
		} else {
			log.Println("⚠️ UsecaseCmd.Run is nil")
		}

		if RepositoryCmd.Run != nil {
			RepositoryCmd.Run(RepositoryCmd, args)
		} else {
			log.Println("⚠️ RepositoryCmd.Run is nil")
		}

		if HandlerCmd.Run != nil {
			HandlerCmd.Run(HandlerCmd, args)
		} else {
			log.Println("⚠️ HandlerCmd.Run is nil")
		}

		log.Println("✅ Feature generated successfully!")
	},
}
