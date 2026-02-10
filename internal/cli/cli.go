package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Napageneral/tokcount/internal/count"
	"github.com/Napageneral/tokcount/internal/ignore"
	"github.com/Napageneral/tokcount/internal/output"
	"github.com/Napageneral/tokcount/internal/tokenizer"
	"github.com/spf13/cobra"
)

// NewRootCmd constructs the tokcount CLI command.
func NewRootCmd() *cobra.Command {
	var (
		outputFormat  string
		tokenizerName string
		ignoreFile    string
		showTree      bool
	)

	cmd := &cobra.Command{
		Use:   "tokcount [path]",
		Short: "Count tokens in a repository",
		Long:  "tokcount scans a repository, counts tokens, and prints an Intent Layer pricing estimate.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := "."
			if len(args) == 1 {
				target = args[0]
			}

			rootPath, err := filepath.Abs(target)
			if err != nil {
				return fmt.Errorf("resolve repository path: %w", err)
			}

			selectedTokenizer, err := tokenizer.New(tokenizerName)
			if err != nil {
				return err
			}

			ignoreSpec, err := ignore.LoadSpec(rootPath, ignoreFile)
			if err != nil {
				return err
			}

			result, err := count.Run(count.Options{
				Root:       rootPath,
				Tokenizer:  selectedTokenizer,
				IgnoreSpec: ignoreSpec,
			})
			if err != nil {
				return err
			}

			switch strings.ToLower(strings.TrimSpace(outputFormat)) {
			case "", "summary":
				fmt.Fprintln(cmd.OutOrStdout(), output.RenderSummary(result))
				if showTree {
					fmt.Fprintln(cmd.OutOrStdout(), output.RenderTree(result))
				}
			case "json":
				payload, err := output.RenderJSON(result)
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(payload))
			default:
				return fmt.Errorf("unsupported output format: %s (use: summary or json)", outputFormat)
			}

			return nil
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&outputFormat, "output", "summary", "Output format: summary | json")
	cmd.Flags().StringVar(&tokenizerName, "tokenizer", "estimate", "Tokenizer: estimate | openai | anthropic")
	cmd.Flags().StringVar(&ignoreFile, "ignore", "", "Custom ignore file path (gitignore syntax)")
	cmd.Flags().BoolVar(&showTree, "tree", false, "Show full directory tree breakdown (summary output only)")

	return cmd
}

// Execute runs the tokcount command and exits on failure.
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
