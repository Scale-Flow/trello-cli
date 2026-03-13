package main

import (
	"github.com/brettmcdowell/trello-cli/internal/auth"
	"github.com/brettmcdowell/trello-cli/internal/contract"
	"github.com/brettmcdowell/trello-cli/internal/trello"
	"github.com/spf13/cobra"
)

var commentsCmd = &cobra.Command{
	Use:   "comments",
	Short: "Manage card comments",
}

var commentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List comments on a card",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		if err := contract.RequireFlag("card", cardID); err != nil {
			return err
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		comments, err := getAPIClient(creds).ListComments(cmd.Context(), cardID)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), comments)
	},
}

var commentsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a comment to a card",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		text, _ := cmd.Flags().GetString("text")
		if err := contract.RequireFlag("card", cardID); err != nil {
			return err
		}
		if err := contract.RequireFlag("text", text); err != nil {
			return err
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		comment, err := getAPIClient(creds).AddComment(cmd.Context(), cardID, text)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), comment)
	},
}

var commentsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a comment by action ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		actionID, _ := cmd.Flags().GetString("action")
		text, _ := cmd.Flags().GetString("text")
		if err := contract.RequireFlag("action", actionID); err != nil {
			return err
		}
		if err := contract.RequireFlag("text", text); err != nil {
			return err
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		comment, err := getAPIClient(creds).UpdateComment(cmd.Context(), actionID, text)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), comment)
	},
}

var commentsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a comment by action ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		actionID, _ := cmd.Flags().GetString("action")
		if err := contract.RequireFlag("action", actionID); err != nil {
			return err
		}

		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		if err := getAPIClient(creds).DeleteComment(cmd.Context(), actionID); err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), trello.DeleteResult{Deleted: true, ID: actionID})
	},
}

func init() {
	commentsListCmd.Flags().String("card", "", "Card ID")
	commentsAddCmd.Flags().String("card", "", "Card ID")
	commentsAddCmd.Flags().String("text", "", "Comment text")
	commentsUpdateCmd.Flags().String("action", "", "Comment action ID")
	commentsUpdateCmd.Flags().String("text", "", "Updated comment text")
	commentsDeleteCmd.Flags().String("action", "", "Comment action ID")

	commentsCmd.AddCommand(commentsListCmd)
	commentsCmd.AddCommand(commentsAddCmd)
	commentsCmd.AddCommand(commentsUpdateCmd)
	commentsCmd.AddCommand(commentsDeleteCmd)
	rootCmd.AddCommand(commentsCmd)
}
