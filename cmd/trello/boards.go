package main

import (
	"github.com/brettmcdowell/trello-cli/internal/auth"
	"github.com/brettmcdowell/trello-cli/internal/contract"
	"github.com/brettmcdowell/trello-cli/internal/credentials"
	"github.com/brettmcdowell/trello-cli/internal/trello"
	"github.com/spf13/cobra"
)

// apiClient is the Trello API client. Overridden in tests with a mock.
var apiClient trello.API

func getAPIClient(creds credentials.Credentials) trello.API {
	if apiClient != nil {
		return apiClient
	}
	apiClient = trello.NewClient(trelloBaseURL, creds.APIKey, creds.Token, trello.DefaultClientOptions())
	return apiClient
}

var boardsCmd = &cobra.Command{
	Use:   "boards",
	Short: "Manage Trello boards",
}

var boardsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List visible boards",
	RunE: func(cmd *cobra.Command, args []string) error {
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		boards, err := getAPIClient(creds).ListBoards(cmd.Context())
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), boards)
	},
}

var boardsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a board by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		boardID, _ := cmd.Flags().GetString("board")
		if err := contract.RequireFlag("board", boardID); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		board, err := getAPIClient(creds).GetBoard(cmd.Context(), boardID)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), board)
	},
}

func init() {
	boardsGetCmd.Flags().String("board", "", "Board ID")

	boardsCmd.AddCommand(boardsListCmd)
	boardsCmd.AddCommand(boardsGetCmd)
	rootCmd.AddCommand(boardsCmd)
}
