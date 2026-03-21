package main

import (
	"github.com/dynamder/synthdata/internal/cli"
	"github.com/dynamder/synthdata/internal/config"
	"github.com/spf13/cobra"
)

func main() {
	config.LoadConfig()
	//synthdata.InitLogger()
	rootCmd := &cobra.Command{
		Use:   "synthdata",
		Short: "Synthetic Dataset Generation Tool",
		Long:  `Generate synthetic datasets from description files using LLM.`,
	}
	rootCmd.AddCommand(cli.GenerateCmd)
	rootCmd.Execute()
}
