package main

import (
	"context"
	"io"

	"github.com/brettmcdowell/trello-cli/internal/auth"
	"github.com/brettmcdowell/trello-cli/internal/contract"
	"github.com/brettmcdowell/trello-cli/internal/credentials"
	"github.com/spf13/cobra"
)

// trelloBaseURL is the Trello API base URL. Variable (not const) so tests can override it.
var trelloBaseURL = "https://api.trello.com"

// credStore is the credential store used by auth commands.
// Overridden in tests with a MemoryStore.
var credStore credentials.Store

type loginCommandContext = context.Context
type loginBrowserOpener = auth.BrowserOpener
type loginOutputWriter = io.Writer

var runAuthLogin = auth.Login

func getCredStore() credentials.Store {
	if credStore != nil {
		return credStore
	}
	credStore = credentials.NewFallbackStore(
		credentials.NewKeyringStore(),
		credentials.NewEnvStore(),
	)
	return credStore
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage Trello authentication",
}

var authSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Store API key and token for subsequent commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, _ := cmd.Flags().GetString("api-key")
		token, _ := cmd.Flags().GetString("token")

		if err := contract.RequireFlag("api-key", apiKey); err != nil {
			return err
		}
		if err := contract.RequireFlag("token", token); err != nil {
			return err
		}

		result, err := auth.Set(getCredStore(), "default", apiKey, token)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

var authSetKeyCmd = &cobra.Command{
	Use:   "set-key",
	Short: "Store the Trello API key for interactive login or key rotation",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, _ := cmd.Flags().GetString("api-key")

		if err := contract.RequireFlag("api-key", apiKey); err != nil {
			return err
		}

		result, err := auth.SetKey(getCredStore(), "default", apiKey)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

var authClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove stored credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := auth.Clear(getCredStore(), "default")
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication state",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := auth.Status(cmd.Context(), getCredStore(), "default", trelloBaseURL)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate via interactive browser login",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := runAuthLogin(
			cmd.Context(),
			getCredStore(),
			"default",
			trelloBaseURL,
			"",
			nil,
			cmd.ErrOrStderr(),
		)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

func init() {
	authSetCmd.Flags().String("api-key", "", "Trello API key")
	authSetCmd.Flags().String("token", "", "Trello user token")
	authSetKeyCmd.Flags().String("api-key", "", "Trello API key")

	authCmd.AddCommand(authSetCmd)
	authCmd.AddCommand(authSetKeyCmd)
	authCmd.AddCommand(authClearCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authLoginCmd)
	rootCmd.AddCommand(authCmd)
}
