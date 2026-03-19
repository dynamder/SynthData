package cli

import (
	"fmt"
	"os"

	"github.com/anomalyco/synthdata/internal/config"
	"github.com/anomalyco/synthdata/internal/services/generator"
	"github.com/anomalyco/synthdata/internal/services/llm"
	"github.com/anomalyco/synthdata/internal/services/parser"
	"github.com/anomalyco/synthdata/internal/services/validator"
	"github.com/spf13/cobra"
)

var (
	descriptionFile string
	outputFile      string
	dataFormat      string
	scale           int
	configFile      string
	force           bool
)

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate synthetic dataset",
	Long:  `Generate synthetic dataset from description file`,
	RunE:  runGenerate,
}

func init() {
	GenerateCmd.Flags().StringVarP(&descriptionFile, "description", "d", "", "Path to description file (required)")
	GenerateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (required)")
	GenerateCmd.Flags().StringVarP(&dataFormat, "format", "f", "json", "Output format: json, csv")
	GenerateCmd.Flags().IntVarP(&scale, "scale", "s", 10, "Number of records to generate")
	GenerateCmd.Flags().StringVarP(&configFile, "config", "c", "", "Config file path")
	GenerateCmd.Flags().BoolVarP(&force, "force", "", false, "Overwrite existing output file")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	if descriptionFile == "" || outputFile == "" {
		return fmt.Errorf("both --description and --output are required")
	}

	if !force {
		if _, err := os.Stat(outputFile); err == nil {
			return fmt.Errorf("output file already exists. Use --force to overwrite")
		}
	}

	descFile, err := parser.ParseDescriptionFile(descriptionFile)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if err := validator.ValidateDescription(descFile); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if scale > 0 {
		descFile.Count = scale
	}

	cfg := config.GetConfig()
	var client llm.Client
	if configFile != "" {
		client = llm.NewOpenAIClientWithConfig(
			cfg.LLM.APIKey,
			cfg.LLM.BaseURL,
			cfg.LLM.Model,
		)
	} else {
		client = llm.NewOpenAIClient()
	}

	gen := generator.New(client)
	records, err := gen.Generate(descFile)
	if err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	formatter, err := generator.GetFormatter(dataFormat)
	if err != nil {
		return fmt.Errorf("formatter error: %w", err)
	}

	output, err := formatter.Format(records)
	if err != nil {
		return fmt.Errorf("formatting failed: %w", err)
	}

	if err := os.WriteFile(outputFile, output, 0644); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	fmt.Printf("Successfully generated %d records to %s\n", len(records), outputFile)
	return nil
}
