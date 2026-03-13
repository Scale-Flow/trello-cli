package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/brettmcdowell/trello-cli/internal/contract"
	"github.com/spf13/cobra"
)

var (
	prettyFlag  bool
	verboseFlag bool
)

var rootCmd = &cobra.Command{
	Use:   "trello",
	Short: "Trello CLI — a machine-friendly Trello interface",
	Long:  "A cross-platform CLI for Trello, designed for coding agents and terminal users. All commands return JSON.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		level := slog.LevelWarn
		if verboseFlag {
			level = slog.LevelDebug
		}
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
		slog.SetDefault(logger)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&prettyFlag, "pretty", false, "Pretty-print JSON output")
	rootCmd.PersistentFlags().BoolVar(&verboseFlag, "verbose", false, "Enable verbose diagnostic output on stderr")
}

// output writes a success envelope for the given data to stdout.
func output(w io.Writer, data any) error {
	envelope, err := contract.Success(data)
	if err != nil {
		return fmt.Errorf("failed to build success envelope: %w", err)
	}
	return contract.Render(w, envelope, prettyFlag)
}

// handleError writes an error envelope to stdout and returns the appropriate exit code.
func handleError(w io.Writer, err error) int {
	ce, ok := err.(*contract.ContractError)
	if !ok {
		ce = &contract.ContractError{Code: contract.UnknownError, Message: err.Error()}
	}
	envelope, marshalErr := contract.ErrorFromContractError(ce)
	if marshalErr != nil {
		// Last resort: write raw error to stderr
		fmt.Fprintf(os.Stderr, "fatal: %v\n", marshalErr)
		return 1
	}
	contract.Render(w, envelope, prettyFlag)
	return 1
}
