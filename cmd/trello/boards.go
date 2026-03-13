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

var boardsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a board",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("desc")
		defaultLists, _ := cmd.Flags().GetBool("default-lists")
		defaultLabels, _ := cmd.Flags().GetBool("default-labels")
		organizationID, _ := cmd.Flags().GetString("organization")
		sourceBoardID, _ := cmd.Flags().GetString("source-board")

		if err := contract.RequireFlag("name", name); err != nil {
			return err
		}

		params := trello.CreateBoardParams{Name: name}
		if desc != "" {
			params.Desc = &desc
		}
		if cmd.Flags().Changed("default-lists") {
			params.DefaultLists = &defaultLists
		}
		if cmd.Flags().Changed("default-labels") {
			params.DefaultLabels = &defaultLabels
		}
		if organizationID != "" {
			params.IDOrganization = &organizationID
		}
		if sourceBoardID != "" {
			params.IDBoardSource = &sourceBoardID
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		board, err := getAPIClient(creds).CreateBoard(cmd.Context(), params)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), board)
	},
}

func init() {
	boardsGetCmd.Flags().String("board", "", "Board ID")
	boardsCreateCmd.Flags().String("name", "", "Board name")
	boardsCreateCmd.Flags().String("desc", "", "Board description")
	boardsCreateCmd.Flags().Bool("default-lists", false, "Create default lists on the new board")
	boardsCreateCmd.Flags().Bool("default-labels", false, "Create default labels on the new board")
	boardsCreateCmd.Flags().String("organization", "", "Workspace or organization ID")
	boardsCreateCmd.Flags().String("source-board", "", "Source board ID to copy from")

	boardsCmd.AddCommand(boardsListCmd)
	boardsCmd.AddCommand(boardsGetCmd)
	boardsCmd.AddCommand(boardsCreateCmd)
	rootCmd.AddCommand(boardsCmd)
}
