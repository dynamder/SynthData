package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/anomalyco/synthdata/internal/config"
	"github.com/anomalyco/synthdata/internal/services/batch"
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
	batchSize       int
	concurrency     int
	maxRetries      int
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
	GenerateCmd.Flags().IntVarP(&batchSize, "batch-size", "", 100, "Number of records per batch (for large scale generation)")
	GenerateCmd.Flags().IntVarP(&concurrency, "concurrency", "", 5, "Maximum parallel LLM calls")
	GenerateCmd.Flags().IntVarP(&maxRetries, "max-retries", "", 3, "Maximum retry attempts for failed records")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	if descriptionFile == "" || outputFile == "" {
		return fmt.Errorf("both --description and --output are required")
	}

	if scale <= 0 {
		return fmt.Errorf("scale must be greater than 0")
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

	if batchSize > 0 && scale > batchSize {
		return runBatchGenerate(client, scale)
	}

	descFile.Count = scale
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

func runBatchGenerate(client llm.Client, targetCount int) error {
	batchCfg := config.GetConfig().GetBatchConfig()
	if batchSize > 0 {
		batchCfg.BatchSize = batchSize
	}
	if concurrency > 0 {
		batchCfg.Concurrency = concurrency
	}
	if maxRetries > 0 {
		batchCfg.MaxRetries = maxRetries
	}

	svc := batch.NewService(client, batchCfg.Concurrency, batchCfg.MaxRetries)
	req := batch.GenerationRequest{
		DescriptionFile: descriptionFile,
		TargetCount:     targetCount,
		BatchSize:       batchCfg.BatchSize,
		Concurrency:     batchCfg.Concurrency,
		MaxRetries:      batchCfg.MaxRetries,
		Format:          dataFormat,
		Output:          outputFile,
	}

	session, err := svc.Generate(context.Background(), req)
	if err != nil {
		return err
	}

	if session.FailedRecords > 0 && session.TotalRecords > 0 {
		os.Exit(3)
	}

	return nil
}
