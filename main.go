package main

import (
	"log"

	"github.com/faidfadjri/gostart/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gostart",
	Short: "gostart is a mini Go framework generator",
}

func main() {
	rootCmd.AddCommand(cmd.CreateCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
