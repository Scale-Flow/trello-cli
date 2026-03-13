package main

import (
	"github.com/brettmcdowell/trello-cli/internal/auth"
	"github.com/brettmcdowell/trello-cli/internal/contract"
	"github.com/brettmcdowell/trello-cli/internal/trello"
	"github.com/spf13/cobra"
)

var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Manage labels",
}

var labelsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List labels on a board",
	RunE: func(cmd *cobra.Command, args []string) error {
		boardID, _ := cmd.Flags().GetString("board")
		if err := contract.RequireFlag("board", boardID); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		labels, err := getAPIClient(creds).ListLabels(cmd.Context(), boardID)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), labels)
	},
}

var labelsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a label on a board",
	RunE: func(cmd *cobra.Command, args []string) error {
		boardID, _ := cmd.Flags().GetString("board")
		name, _ := cmd.Flags().GetString("name")
		color, _ := cmd.Flags().GetString("color")
		if err := contract.RequireFlag("board", boardID); err != nil {
			return err
		}
		if err := contract.RequireFlag("name", name); err != nil {
			return err
		}
		if err := contract.RequireFlag("color", color); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		label, err := getAPIClient(creds).CreateLabel(cmd.Context(), boardID, name, color)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), label)
	},
}

var labelsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a label to a card",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		labelID, _ := cmd.Flags().GetString("label")
		if err := contract.RequireFlag("card", cardID); err != nil {
			return err
		}
		if err := contract.RequireFlag("label", labelID); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		if err := getAPIClient(creds).AddLabelToCard(cmd.Context(), cardID, labelID); err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), trello.ActionResult{Success: true, ID: labelID})
	},
}

var labelsRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a label from a card",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardID, _ := cmd.Flags().GetString("card")
		labelID, _ := cmd.Flags().GetString("label")
		if err := contract.RequireFlag("card", cardID); err != nil {
			return err
		}
		if err := contract.RequireFlag("label", labelID); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		if err := getAPIClient(creds).RemoveLabelFromCard(cmd.Context(), cardID, labelID); err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), trello.ActionResult{Success: true, ID: labelID})
	},
}

func init() {
	labelsListCmd.Flags().String("board", "", "Board ID")
	labelsCreateCmd.Flags().String("board", "", "Board ID")
	labelsCreateCmd.Flags().String("name", "", "Label name")
	labelsCreateCmd.Flags().String("color", "", "Label color")
	labelsAddCmd.Flags().String("card", "", "Card ID")
	labelsAddCmd.Flags().String("label", "", "Label ID")
	labelsRemoveCmd.Flags().String("card", "", "Card ID")
	labelsRemoveCmd.Flags().String("label", "", "Label ID")

	labelsCmd.AddCommand(labelsListCmd)
	labelsCmd.AddCommand(labelsCreateCmd)
	labelsCmd.AddCommand(labelsAddCmd)
	labelsCmd.AddCommand(labelsRemoveCmd)
	rootCmd.AddCommand(labelsCmd)
}
