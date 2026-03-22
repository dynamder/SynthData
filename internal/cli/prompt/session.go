package prompt

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type PromptSession struct {
	Command       *cobra.Command
	ProvidedArgs  map[string]interface{}
	CollectedArgs map[string]interface{}
	Metadata      []ArgumentMetadata
}

func NewPromptSession(cmd *cobra.Command, providedArgs map[string]interface{}) (*PromptSession, error) {
	metadata, err := ExtractArgumentMetadata(cmd)
	if err != nil {
		return nil, err
	}

	return &PromptSession{
		Command:       cmd,
		ProvidedArgs:  providedArgs,
		CollectedArgs: make(map[string]interface{}),
		Metadata:      metadata,
	}, nil
}

func (s *PromptSession) Run() (map[string]interface{}, error) {
	fmt.Print(ShowWelcome())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)

	go func() {
		<-sigChan
		fmt.Println("\n\nExiting interactive mode. No changes made.")
		os.Exit(0)
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for i, arg := range s.Metadata {
		if s.isArgProvided(arg.Name) {
			continue
		}

		if !arg.IsRequired && arg.DefaultValue == nil {
			continue
		}

		for {
			prompt := BuildPrompt(arg, i, len(s.Metadata))
			fmt.Print(prompt)

			if !scanner.Scan() {
				return nil, fmt.Errorf("input closed")
			}

			input := scanner.Text()

			if input == "" {
				if arg.IsRequired {
					fmt.Print(ShowError("This field is required. Please enter a value.\n"))
					continue
				}

				if arg.DefaultValue != nil {
					s.CollectedArgs[arg.Name] = arg.DefaultValue
					fmt.Print(ShowSuccess(fmt.Sprintf("Using default value: %v\n", arg.DefaultValue)))
					break
				}

				break
			}

			result := ValidateInput(input, arg)

			if result.Value == "help" {
				fmt.Print(ShowHelp(arg))
				continue
			}

			if !result.IsValid {
				fmt.Print(ShowError(result.ErrorMessage))
				continue
			}

			s.CollectedArgs[arg.Name] = result.Value
			break
		}
	}

	for k, v := range s.ProvidedArgs {
		s.CollectedArgs[k] = v
	}

	return s.CollectedArgs, nil
}

func (s *PromptSession) isArgProvided(name string) bool {
	_, provided := s.ProvidedArgs[name]
	if provided {
		return true
	}

	flag := s.Command.Flags().Lookup(name)
	if flag != nil && flag.Value.String() != "" && flag.Value.String() != flag.DefValue {
		return true
	}

	return false
}

func GetFlagValues(cmd *cobra.Command, args map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if val, ok := args[flag.Name]; ok {
			result[flag.Name] = val
		} else if flag.Value.String() != "" && flag.Value.String() != flag.DefValue {
			result[flag.Name] = flag.Value.String()
		}
	})

	return result
}

func (s *PromptSession) ShowSummary() {
	fmt.Print(ShowCollectedArgs(s.CollectedArgs))
}
