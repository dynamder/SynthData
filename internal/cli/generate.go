package cli

import (
	"context"
	"fmt"
	"os"

	synthdata "github.com/dynamder/synthdata/internal"
	"github.com/dynamder/synthdata/internal/config"
	"github.com/dynamder/synthdata/internal/services/batch"
	"github.com/dynamder/synthdata/internal/services/generator"
	"github.com/dynamder/synthdata/internal/services/llm"
	"github.com/dynamder/synthdata/internal/services/parser"
	"github.com/dynamder/synthdata/internal/services/validator"
	"github.com/spf13/cobra"
)

var (
	descriptionFile string
	outputFile      string
	logDir          string
	dataFormat      string
	scale           int
	configFile      string
	force           bool
	batchSize       int
	concurrency     int
	maxRetries      int
	verbose         bool
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
	GenerateCmd.Flags().StringVarP(&logDir, "log", "l", "log/", "Directory to hold logs")
	GenerateCmd.Flags().StringVarP(&dataFormat, "format", "f", "json", "Output format: json, csv")
	GenerateCmd.Flags().IntVarP(&scale, "scale", "s", 10, "Number of records to generate")
	GenerateCmd.Flags().StringVarP(&configFile, "config", "c", "", "Config file path")
	GenerateCmd.Flags().BoolVarP(&force, "force", "", false, "Overwrite existing output file")
	GenerateCmd.Flags().IntVarP(&batchSize, "batch-size", "", 10, "Number of records per batch (for large scale generation)")
	GenerateCmd.Flags().IntVarP(&concurrency, "concurrency", "", 5, "Maximum parallel LLM calls")
	GenerateCmd.Flags().IntVarP(&maxRetries, "max-retries", "", 3, "Maximum retry attempts for failed records")
	GenerateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
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

	if err := synthdata.InitLogger(logDir); err != nil {
		return fmt.Errorf("log error: %w", err)
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
		return runBatchGenerate(client, scale, verbose)
	}

	descFile.Count = scale
	gen := generator.New(client)

	//TODO: the generator should write to the output file via a stream.
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

func runBatchGenerate(client llm.Client, targetCount int, verbose bool) error {
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

	svc := batch.NewService(client, batchCfg.Concurrency, batchCfg.MaxRetries, verbose)
	req := batch.GenerationRequest{
		DescriptionFile: descriptionFile,
		TargetCount:     targetCount,
		BatchSize:       batchCfg.BatchSize,
		Concurrency:     batchCfg.Concurrency,
		MaxRetries:      batchCfg.MaxRetries,
		Format:          dataFormat,
		Output:          outputFile,
		Verbose:         verbose,
	}

	_, err := svc.Generate(context.Background(), req)
	if err != nil {
		return err
	}

	return nil
}
