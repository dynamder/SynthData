package main

import (
	"github.com/anomalyco/synthdata/internal/cli"
	"github.com/anomalyco/synthdata/internal/config"
	"github.com/spf13/cobra"
)

func main() {
	config.LoadConfig()
	rootCmd := &cobra.Command{
		Use:   "synthdata",
		Short: "Synthetic Dataset Generation Tool",
		Long:  `Generate synthetic datasets from description files using LLM.`,
	}
	rootCmd.AddCommand(cli.GenerateCmd)
	rootCmd.Execute()
}
