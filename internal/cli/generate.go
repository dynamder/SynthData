package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	synthdata "github.com/dynamder/synthdata/internal"
	"github.com/dynamder/synthdata/internal/cli/prompt"
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
	interactive     bool
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
	GenerateCmd.Flags().StringVarP(&configFile, "config", "c", "configs/default.toml", "Config file path")
	GenerateCmd.Flags().BoolVarP(&force, "force", "", false, "Overwrite existing output file")
	GenerateCmd.Flags().IntVarP(&batchSize, "batch-size", "", 10, "Number of records per batch (for large scale generation)")
	GenerateCmd.Flags().IntVarP(&concurrency, "concurrency", "", 5, "Maximum parallel LLM calls")
	GenerateCmd.Flags().IntVarP(&maxRetries, "max-retries", "", 3, "Maximum retry attempts for failed records")
	GenerateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	GenerateCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Enable interactive mode")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	err := runGenerateInner(cmd)
	if err != nil {
		if isLLMError(err) {
			fmt.Fprintf(os.Stderr, "\nLLM Error: %v\n", err)
			if strings.Contains(err.Error(), "Detail:") {
				fmt.Fprintln(os.Stderr, "")
				fmt.Fprintln(os.Stderr, "Please check your LLM configuration:")
				fmt.Fprintln(os.Stderr, "  - API Base URL: check the 'base_url' in your config file")
				fmt.Fprintln(os.Stderr, "  - API Key: ensure OPENAI_API_KEY is set or 'api_key' is in your config")
				fmt.Fprintln(os.Stderr, "  - Model Name: verify the 'model' setting in your config")
			}
			os.Exit(1)
		}
		printArgError(err)
		return nil
	}
	return nil
}

func printArgError(err error) {
	errStr := err.Error()
	if strings.Contains(errStr, "both --description and --output") {
		fmt.Fprintf(os.Stderr, "error: required arguments `-d, --description` and `-o, --output` are not set\n")
	} else if strings.Contains(errStr, "--description") {
		fmt.Fprintf(os.Stderr, "error: required argument `-d, --description` is not set\n")
	} else if strings.Contains(errStr, "--output") {
		fmt.Fprintf(os.Stderr, "error: required argument `-o, --output` is not set\n")
	} else if strings.Contains(errStr, "scale") {
		fmt.Fprintf(os.Stderr, "error: option `--scale` must be greater than 0\n")
	} else if strings.Contains(errStr, "output file already exists") {
		fmt.Fprintf(os.Stderr, "error: %s\n", errStr)
		fmt.Fprintln(os.Stderr, "Use `--force` to overwrite")
	} else {
		fmt.Fprintf(os.Stderr, "error: %s\n", errStr)
	}
}

func runGenerateInner(cmd *cobra.Command) error {
	if interactive {
		return runInteractiveGenerate(cmd)
	}

	if descriptionFile == "" || outputFile == "" {
		return fmt.Errorf("both --description and --output are required")
	}

	return runGenerateWithArgs(cmd)
}

func isLLMError(err error) bool {
	errStr := strings.ToLower(err.Error())
	llmKeywords := []string{
		"llm",
		"api",
		"openai",
		"authentication",
		"unauthorized",
		"rate limit",
		"quota",
		"timeout",
		"connection",
		"network",
		"request failed",
		"client",
		"context deadline",
		"generation failed",
		"all records failed",
	}
	for _, keyword := range llmKeywords {
		if strings.Contains(errStr, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func runInteractiveGenerate(cmd *cobra.Command) error {
	providedArgs := make(map[string]interface{})

	if descriptionFile != "" {
		providedArgs["description"] = descriptionFile
	}
	if outputFile != "" {
		providedArgs["output"] = outputFile
	}
	if logDir != "" && logDir != "log/" {
		providedArgs["log"] = logDir
	}
	if dataFormat != "" && dataFormat != "json" {
		providedArgs["format"] = dataFormat
	}
	if scale != 0 && scale != 10 {
		providedArgs["scale"] = scale
	}
	if configFile != "" && configFile != "configs/default.toml" {
		providedArgs["config"] = configFile
	}
	if force {
		providedArgs["force"] = force
	}
	if batchSize != 0 && batchSize != 10 {
		providedArgs["batch-size"] = batchSize
	}
	if concurrency != 0 && concurrency != 5 {
		providedArgs["concurrency"] = concurrency
	}
	if maxRetries != 0 && maxRetries != 3 {
		providedArgs["max-retries"] = maxRetries
	}
	if verbose {
		providedArgs["verbose"] = verbose
	}

	session, err := prompt.NewPromptSession(cmd, providedArgs)
	if err != nil {
		return fmt.Errorf("failed to initialize interactive session: %w", err)
	}

	collectedArgs, err := session.Run()
	if err != nil {
		return fmt.Errorf("interactive session failed: %w", err)
	}

	if descriptionFileVal, ok := collectedArgs["description"].(string); ok {
		descriptionFile = descriptionFileVal
	}
	if outputFileVal, ok := collectedArgs["output"].(string); ok {
		outputFile = outputFileVal
	}
	if logDirVal, ok := collectedArgs["log"].(string); ok {
		logDir = logDirVal
	}
	if dataFormatVal, ok := collectedArgs["format"].(string); ok {
		dataFormat = dataFormatVal
	}
	if scaleVal, ok := collectedArgs["scale"].(int); ok {
		scale = scaleVal
	}
	if configFileVal, ok := collectedArgs["config"].(string); ok {
		configFile = configFileVal
	}
	if forceVal, ok := collectedArgs["force"].(bool); ok {
		force = forceVal
	}
	if batchSizeVal, ok := collectedArgs["batch-size"].(int); ok {
		batchSize = batchSizeVal
	}
	if concurrencyVal, ok := collectedArgs["concurrency"].(int); ok {
		concurrency = concurrencyVal
	}
	if maxRetriesVal, ok := collectedArgs["max-retries"].(int); ok {
		maxRetries = maxRetriesVal
	}
	if verboseVal, ok := collectedArgs["verbose"].(bool); ok {
		verbose = verboseVal
	}

	session.ShowSummary()

	return runGenerateWithArgs(cmd)
}

func runGenerateWithArgs(cmd *cobra.Command) error {
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
