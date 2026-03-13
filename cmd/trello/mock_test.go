package main

import (
	"bytes"
	"context"
	"testing"

	"github.com/Scale-Flow/trello-cli/internal/auth"
	"github.com/Scale-Flow/trello-cli/internal/contract"
	"github.com/Scale-Flow/trello-cli/internal/credentials"
	"github.com/Scale-Flow/trello-cli/internal/trello"
)

// mockAPI implements trello.API for command testing.
type mockAPI struct {
	trello.API
	listBoardsFn               func(ctx context.Context) ([]trello.Board, error)
	getBoardFn                 func(ctx context.Context, id string) (trello.Board, error)
	createBoardFn              func(ctx context.Context, params trello.CreateBoardParams) (trello.Board, error)
	listListsFn                func(ctx context.Context, boardID string) ([]trello.List, error)
	createListFn               func(ctx context.Context, boardID, name string) (trello.List, error)
	updateListFn               func(ctx context.Context, listID string, params trello.UpdateListParams) (trello.List, error)
	archiveListFn              func(ctx context.Context, listID string) (trello.List, error)
	moveListFn                 func(ctx context.Context, listID, boardID string, pos *float64) (trello.List, error)
	listCardsByBoardFn         func(ctx context.Context, boardID string) ([]trello.Card, error)
	listCardsByListFn          func(ctx context.Context, listID string) ([]trello.Card, error)
	getCardFn                  func(ctx context.Context, cardID string) (trello.Card, error)
	createCardFn               func(ctx context.Context, params trello.CreateCardParams) (trello.Card, error)
	updateCardFn               func(ctx context.Context, cardID string, params trello.UpdateCardParams) (trello.Card, error)
	moveCardFn                 func(ctx context.Context, cardID, listID string, pos *float64) (trello.Card, error)
	archiveCardFn              func(ctx context.Context, cardID string) (trello.Card, error)
	deleteCardFn               func(ctx context.Context, cardID string) error
	listCommentsFn             func(ctx context.Context, cardID string) ([]trello.Comment, error)
	addCommentFn               func(ctx context.Context, cardID, text string) (trello.Comment, error)
	updateCommentFn            func(ctx context.Context, actionID, text string) (trello.Comment, error)
	deleteCommentFn            func(ctx context.Context, actionID string) error
	listChecklistsFn           func(ctx context.Context, cardID string) ([]trello.Checklist, error)
	createChecklistFn          func(ctx context.Context, cardID, name string) (trello.Checklist, error)
	deleteChecklistFn          func(ctx context.Context, checklistID string) error
	addCheckItemFn             func(ctx context.Context, checklistID, name string) (trello.CheckItem, error)
	updateCheckItemFn          func(ctx context.Context, cardID, itemID, state string) (trello.CheckItem, error)
	deleteCheckItemFn          func(ctx context.Context, checklistID, itemID string) error
	listAttachmentsFn          func(ctx context.Context, cardID string) ([]trello.Attachment, error)
	addFileAttachmentFn        func(ctx context.Context, cardID, filePath string, name *string) (trello.Attachment, error)
	addURLAttachmentFn         func(ctx context.Context, cardID, urlStr string, name *string) (trello.Attachment, error)
	deleteAttachmentFn         func(ctx context.Context, cardID, attachmentID string) error
	listLabelsFn               func(ctx context.Context, boardID string) ([]trello.Label, error)
	createLabelFn              func(ctx context.Context, boardID, name, color string) (trello.Label, error)
	addLabelToCardFn           func(ctx context.Context, cardID, labelID string) error
	removeLabelFromCardFn      func(ctx context.Context, cardID, labelID string) error
	listCustomFieldsByBoardFn  func(ctx context.Context, boardID string) ([]trello.CustomField, error)
	getCustomFieldFn           func(ctx context.Context, fieldID string) (trello.CustomField, error)
	createCustomFieldFn        func(ctx context.Context, params trello.CreateCustomFieldParams) (trello.CustomField, error)
	updateCustomFieldFn        func(ctx context.Context, fieldID string, params trello.UpdateCustomFieldParams) (trello.CustomField, error)
	deleteCustomFieldFn        func(ctx context.Context, fieldID string) error
	listCustomFieldOptionsFn   func(ctx context.Context, fieldID string) ([]trello.CustomFieldOption, error)
	createCustomFieldOptionFn  func(ctx context.Context, fieldID string, params trello.CreateCustomFieldOptionParams) (trello.CustomFieldOption, error)
	updateCustomFieldOptionFn  func(ctx context.Context, fieldID, optionID string, params trello.UpdateCustomFieldOptionParams) (trello.CustomFieldOption, error)
	deleteCustomFieldOptionFn  func(ctx context.Context, fieldID, optionID string) error
	listCardCustomFieldItemsFn func(ctx context.Context, cardID string) ([]trello.CardCustomFieldItem, error)
	setCardCustomFieldItemFn   func(ctx context.Context, cardID, fieldID string, params trello.SetCardCustomFieldItemParams) (trello.CardCustomFieldItem, error)
	clearCardCustomFieldItemFn func(ctx context.Context, cardID, fieldID string) error
	listMembersFn              func(ctx context.Context, boardID string) ([]trello.Member, error)
	addMemberToCardFn          func(ctx context.Context, cardID, memberID string) error
	removeMemberFromCardFn     func(ctx context.Context, cardID, memberID string) error
	searchCardsFn              func(ctx context.Context, query string) (trello.CardSearchResult, error)
	searchBoardsFn             func(ctx context.Context, query string) (trello.BoardSearchResult, error)
	getMeFn                    func(ctx context.Context) (trello.Member, error)
}

// Boards
func (m *mockAPI) ListBoards(ctx context.Context) ([]trello.Board, error) {
	if m.listBoardsFn != nil {
		return m.listBoardsFn(ctx)
	}
	return nil, nil
}

func (m *mockAPI) GetBoard(ctx context.Context, id string) (trello.Board, error) {
	if m.getBoardFn != nil {
		return m.getBoardFn(ctx, id)
	}
	return trello.Board{}, nil
}

func (m *mockAPI) CreateBoard(ctx context.Context, params trello.CreateBoardParams) (trello.Board, error) {
	if m.createBoardFn != nil {
		return m.createBoardFn(ctx, params)
	}
	return trello.Board{}, nil
}

// Lists
func (m *mockAPI) ListLists(ctx context.Context, boardID string) ([]trello.List, error) {
	if m.listListsFn != nil {
		return m.listListsFn(ctx, boardID)
	}
	return nil, nil
}

func (m *mockAPI) CreateList(ctx context.Context, boardID, name string) (trello.List, error) {
	if m.createListFn != nil {
		return m.createListFn(ctx, boardID, name)
	}
	return trello.List{}, nil
}

func (m *mockAPI) UpdateList(ctx context.Context, listID string, params trello.UpdateListParams) (trello.List, error) {
	if m.updateListFn != nil {
		return m.updateListFn(ctx, listID, params)
	}
	return trello.List{}, nil
}

func (m *mockAPI) ArchiveList(ctx context.Context, listID string) (trello.List, error) {
	if m.archiveListFn != nil {
		return m.archiveListFn(ctx, listID)
	}
	return trello.List{}, nil
}

func (m *mockAPI) MoveList(ctx context.Context, listID, boardID string, pos *float64) (trello.List, error) {
	if m.moveListFn != nil {
		return m.moveListFn(ctx, listID, boardID, pos)
	}
	return trello.List{}, nil
}

// Cards
func (m *mockAPI) ListCardsByBoard(ctx context.Context, boardID string) ([]trello.Card, error) {
	if m.listCardsByBoardFn != nil {
		return m.listCardsByBoardFn(ctx, boardID)
	}
	return nil, nil
}

func (m *mockAPI) ListCardsByList(ctx context.Context, listID string) ([]trello.Card, error) {
	if m.listCardsByListFn != nil {
		return m.listCardsByListFn(ctx, listID)
	}
	return nil, nil
}

func (m *mockAPI) GetCard(ctx context.Context, cardID string) (trello.Card, error) {
	if m.getCardFn != nil {
		return m.getCardFn(ctx, cardID)
	}
	return trello.Card{}, nil
}

func (m *mockAPI) CreateCard(ctx context.Context, params trello.CreateCardParams) (trello.Card, error) {
	if m.createCardFn != nil {
		return m.createCardFn(ctx, params)
	}
	return trello.Card{}, nil
}

func (m *mockAPI) UpdateCard(ctx context.Context, cardID string, params trello.UpdateCardParams) (trello.Card, error) {
	if m.updateCardFn != nil {
		return m.updateCardFn(ctx, cardID, params)
	}
	return trello.Card{}, nil
}

func (m *mockAPI) MoveCard(ctx context.Context, cardID, listID string, pos *float64) (trello.Card, error) {
	if m.moveCardFn != nil {
		return m.moveCardFn(ctx, cardID, listID, pos)
	}
	return trello.Card{}, nil
}

func (m *mockAPI) ArchiveCard(ctx context.Context, cardID string) (trello.Card, error) {
	if m.archiveCardFn != nil {
		return m.archiveCardFn(ctx, cardID)
	}
	return trello.Card{}, nil
}

func (m *mockAPI) DeleteCard(ctx context.Context, cardID string) error {
	if m.deleteCardFn != nil {
		return m.deleteCardFn(ctx, cardID)
	}
	return nil
}

// Comments
func (m *mockAPI) ListComments(ctx context.Context, cardID string) ([]trello.Comment, error) {
	if m.listCommentsFn != nil {
		return m.listCommentsFn(ctx, cardID)
	}
	return nil, nil
}

func (m *mockAPI) AddComment(ctx context.Context, cardID, text string) (trello.Comment, error) {
	if m.addCommentFn != nil {
		return m.addCommentFn(ctx, cardID, text)
	}
	return trello.Comment{}, nil
}

func (m *mockAPI) UpdateComment(ctx context.Context, actionID, text string) (trello.Comment, error) {
	if m.updateCommentFn != nil {
		return m.updateCommentFn(ctx, actionID, text)
	}
	return trello.Comment{}, nil
}

func (m *mockAPI) DeleteComment(ctx context.Context, actionID string) error {
	if m.deleteCommentFn != nil {
		return m.deleteCommentFn(ctx, actionID)
	}
	return nil
}

// Checklists
func (m *mockAPI) ListChecklists(ctx context.Context, cardID string) ([]trello.Checklist, error) {
	if m.listChecklistsFn != nil {
		return m.listChecklistsFn(ctx, cardID)
	}
	return nil, nil
}

func (m *mockAPI) CreateChecklist(ctx context.Context, cardID, name string) (trello.Checklist, error) {
	if m.createChecklistFn != nil {
		return m.createChecklistFn(ctx, cardID, name)
	}
	return trello.Checklist{}, nil
}

func (m *mockAPI) DeleteChecklist(ctx context.Context, checklistID string) error {
	if m.deleteChecklistFn != nil {
		return m.deleteChecklistFn(ctx, checklistID)
	}
	return nil
}

func (m *mockAPI) AddCheckItem(ctx context.Context, checklistID, name string) (trello.CheckItem, error) {
	if m.addCheckItemFn != nil {
		return m.addCheckItemFn(ctx, checklistID, name)
	}
	return trello.CheckItem{}, nil
}

func (m *mockAPI) UpdateCheckItem(ctx context.Context, cardID, itemID, state string) (trello.CheckItem, error) {
	if m.updateCheckItemFn != nil {
		return m.updateCheckItemFn(ctx, cardID, itemID, state)
	}
	return trello.CheckItem{}, nil
}

func (m *mockAPI) DeleteCheckItem(ctx context.Context, checklistID, itemID string) error {
	if m.deleteCheckItemFn != nil {
		return m.deleteCheckItemFn(ctx, checklistID, itemID)
	}
	return nil
}

// Attachments
func (m *mockAPI) ListAttachments(ctx context.Context, cardID string) ([]trello.Attachment, error) {
	if m.listAttachmentsFn != nil {
		return m.listAttachmentsFn(ctx, cardID)
	}
	return nil, nil
}

func (m *mockAPI) AddFileAttachment(ctx context.Context, cardID, filePath string, name *string) (trello.Attachment, error) {
	if m.addFileAttachmentFn != nil {
		return m.addFileAttachmentFn(ctx, cardID, filePath, name)
	}
	return trello.Attachment{}, nil
}

func (m *mockAPI) AddURLAttachment(ctx context.Context, cardID, urlStr string, name *string) (trello.Attachment, error) {
	if m.addURLAttachmentFn != nil {
		return m.addURLAttachmentFn(ctx, cardID, urlStr, name)
	}
	return trello.Attachment{}, nil
}

func (m *mockAPI) DeleteAttachment(ctx context.Context, cardID, attachmentID string) error {
	if m.deleteAttachmentFn != nil {
		return m.deleteAttachmentFn(ctx, cardID, attachmentID)
	}
	return nil
}

// Labels
func (m *mockAPI) ListLabels(ctx context.Context, boardID string) ([]trello.Label, error) {
	if m.listLabelsFn != nil {
		return m.listLabelsFn(ctx, boardID)
	}
	return nil, nil
}

func (m *mockAPI) CreateLabel(ctx context.Context, boardID, name, color string) (trello.Label, error) {
	if m.createLabelFn != nil {
		return m.createLabelFn(ctx, boardID, name, color)
	}
	return trello.Label{}, nil
}

func (m *mockAPI) AddLabelToCard(ctx context.Context, cardID, labelID string) error {
	if m.addLabelToCardFn != nil {
		return m.addLabelToCardFn(ctx, cardID, labelID)
	}
	return nil
}

func (m *mockAPI) RemoveLabelFromCard(ctx context.Context, cardID, labelID string) error {
	if m.removeLabelFromCardFn != nil {
		return m.removeLabelFromCardFn(ctx, cardID, labelID)
	}
	return nil
}

// Custom Fields
func (m *mockAPI) ListCustomFieldsByBoard(ctx context.Context, boardID string) ([]trello.CustomField, error) {
	if m.listCustomFieldsByBoardFn != nil {
		return m.listCustomFieldsByBoardFn(ctx, boardID)
	}
	return nil, nil
}

func (m *mockAPI) GetCustomField(ctx context.Context, fieldID string) (trello.CustomField, error) {
	if m.getCustomFieldFn != nil {
		return m.getCustomFieldFn(ctx, fieldID)
	}
	return trello.CustomField{}, nil
}

func (m *mockAPI) CreateCustomField(ctx context.Context, params trello.CreateCustomFieldParams) (trello.CustomField, error) {
	if m.createCustomFieldFn != nil {
		return m.createCustomFieldFn(ctx, params)
	}
	return trello.CustomField{}, nil
}

func (m *mockAPI) UpdateCustomField(ctx context.Context, fieldID string, params trello.UpdateCustomFieldParams) (trello.CustomField, error) {
	if m.updateCustomFieldFn != nil {
		return m.updateCustomFieldFn(ctx, fieldID, params)
	}
	return trello.CustomField{}, nil
}

func (m *mockAPI) DeleteCustomField(ctx context.Context, fieldID string) error {
	if m.deleteCustomFieldFn != nil {
		return m.deleteCustomFieldFn(ctx, fieldID)
	}
	return nil
}

func (m *mockAPI) ListCustomFieldOptions(ctx context.Context, fieldID string) ([]trello.CustomFieldOption, error) {
	if m.listCustomFieldOptionsFn != nil {
		return m.listCustomFieldOptionsFn(ctx, fieldID)
	}
	return nil, nil
}

func (m *mockAPI) CreateCustomFieldOption(ctx context.Context, fieldID string, params trello.CreateCustomFieldOptionParams) (trello.CustomFieldOption, error) {
	if m.createCustomFieldOptionFn != nil {
		return m.createCustomFieldOptionFn(ctx, fieldID, params)
	}
	return trello.CustomFieldOption{}, nil
}

func (m *mockAPI) UpdateCustomFieldOption(ctx context.Context, fieldID, optionID string, params trello.UpdateCustomFieldOptionParams) (trello.CustomFieldOption, error) {
	if m.updateCustomFieldOptionFn != nil {
		return m.updateCustomFieldOptionFn(ctx, fieldID, optionID, params)
	}
	return trello.CustomFieldOption{}, nil
}

func (m *mockAPI) DeleteCustomFieldOption(ctx context.Context, fieldID, optionID string) error {
	if m.deleteCustomFieldOptionFn != nil {
		return m.deleteCustomFieldOptionFn(ctx, fieldID, optionID)
	}
	return nil
}

func (m *mockAPI) ListCardCustomFieldItems(ctx context.Context, cardID string) ([]trello.CardCustomFieldItem, error) {
	if m.listCardCustomFieldItemsFn != nil {
		return m.listCardCustomFieldItemsFn(ctx, cardID)
	}
	return nil, nil
}

func (m *mockAPI) SetCardCustomFieldItem(ctx context.Context, cardID, fieldID string, params trello.SetCardCustomFieldItemParams) (trello.CardCustomFieldItem, error) {
	if m.setCardCustomFieldItemFn != nil {
		return m.setCardCustomFieldItemFn(ctx, cardID, fieldID, params)
	}
	return trello.CardCustomFieldItem{}, nil
}

func (m *mockAPI) ClearCardCustomFieldItem(ctx context.Context, cardID, fieldID string) error {
	if m.clearCardCustomFieldItemFn != nil {
		return m.clearCardCustomFieldItemFn(ctx, cardID, fieldID)
	}
	return nil
}

// Members
func (m *mockAPI) ListMembers(ctx context.Context, boardID string) ([]trello.Member, error) {
	if m.listMembersFn != nil {
		return m.listMembersFn(ctx, boardID)
	}
	return nil, nil
}

func (m *mockAPI) AddMemberToCard(ctx context.Context, cardID, memberID string) error {
	if m.addMemberToCardFn != nil {
		return m.addMemberToCardFn(ctx, cardID, memberID)
	}
	return nil
}

func (m *mockAPI) RemoveMemberFromCard(ctx context.Context, cardID, memberID string) error {
	if m.removeMemberFromCardFn != nil {
		return m.removeMemberFromCardFn(ctx, cardID, memberID)
	}
	return nil
}

// Search
func (m *mockAPI) SearchCards(ctx context.Context, query string) (trello.CardSearchResult, error) {
	if m.searchCardsFn != nil {
		return m.searchCardsFn(ctx, query)
	}
	return trello.CardSearchResult{}, nil
}

func (m *mockAPI) SearchBoards(ctx context.Context, query string) (trello.BoardSearchResult, error) {
	if m.searchBoardsFn != nil {
		return m.searchBoardsFn(ctx, query)
	}
	return trello.BoardSearchResult{}, nil
}

// Auth
func (m *mockAPI) GetMe(ctx context.Context) (trello.Member, error) {
	if m.getMeFn != nil {
		return m.getMeFn(ctx)
	}
	return trello.Member{}, nil
}

// setupTestAuth resets test state to avoid cross-test contamination.
func setupTestAuth(t *testing.T) {
	t.Helper()
	// Use a memory store for tests
	credStore = credentials.NewMemoryStore()
	apiClient = nil
	runAuthLogin = auth.Login
	runAuthLoginWithDeviceFlow = auth.LoginWithDeviceFlow
	// Reset root command output state to avoid cross-test contamination
	rootCmd.SetOut(nil)
	rootCmd.SetArgs(nil)
	if err := authSetCmd.Flags().Set("api-key", ""); err != nil {
		t.Fatalf("failed to reset api-key flag: %v", err)
	}
	if err := authSetCmd.Flags().Set("token", ""); err != nil {
		t.Fatalf("failed to reset token flag: %v", err)
	}
	if err := authSetKeyCmd.Flags().Set("api-key", ""); err != nil {
		t.Fatalf("failed to reset set-key api-key flag: %v", err)
	}
	resetFlag := func(cmdName, flagName, value string) {
		t.Helper()
		cmd, _, err := rootCmd.Find([]string{cmdName})
		if err != nil {
			return
		}
		flag := cmd.Flags().Lookup(flagName)
		if flag == nil {
			return
		}
		if err := cmd.Flags().Set(flagName, value); err != nil {
			t.Fatalf("failed to reset %s --%s flag: %v", cmdName, flagName, err)
		}
		flag.Changed = false
	}
	resetSubFlag := func(parentName, childName, flagName, value string) {
		t.Helper()
		parent, _, err := rootCmd.Find([]string{parentName})
		if err != nil {
			return
		}
		cmd, _, err := parent.Find([]string{childName})
		if err != nil {
			return
		}
		flag := cmd.Flags().Lookup(flagName)
		if flag == nil {
			return
		}
		if err := cmd.Flags().Set(flagName, value); err != nil {
			t.Fatalf("failed to reset %s %s --%s flag: %v", parentName, childName, flagName, err)
		}
		flag.Changed = false
	}
	resetPathFlag := func(path []string, flagName, value string) {
		t.Helper()
		cmd, _, err := rootCmd.Find(path)
		if err != nil {
			return
		}
		flag := cmd.Flags().Lookup(flagName)
		if flag == nil {
			return
		}
		if err := cmd.Flags().Set(flagName, value); err != nil {
			t.Fatalf("failed to reset %v --%s flag: %v", path, flagName, err)
		}
		flag.Changed = false
	}
	resetSubFlag("boards", "get", "board", "")
	resetSubFlag("boards", "create", "name", "")
	resetSubFlag("boards", "create", "desc", "")
	resetSubFlag("boards", "create", "default-lists", "false")
	resetSubFlag("boards", "create", "default-labels", "false")
	resetSubFlag("boards", "create", "organization", "")
	resetSubFlag("boards", "create", "source-board", "")
	resetSubFlag("lists", "list", "board", "")
	resetSubFlag("lists", "create", "board", "")
	resetSubFlag("lists", "create", "name", "")
	resetSubFlag("lists", "update", "list", "")
	resetSubFlag("lists", "update", "name", "")
	resetSubFlag("lists", "update", "pos", "0")
	resetSubFlag("lists", "archive", "list", "")
	resetSubFlag("lists", "move", "list", "")
	resetSubFlag("lists", "move", "board", "")
	resetSubFlag("lists", "move", "pos", "0")
	resetSubFlag("cards", "list", "board", "")
	resetSubFlag("cards", "list", "list", "")
	resetSubFlag("cards", "get", "card", "")
	resetSubFlag("cards", "create", "list", "")
	resetSubFlag("cards", "create", "name", "")
	resetSubFlag("cards", "create", "desc", "")
	resetSubFlag("cards", "create", "due", "")
	resetSubFlag("cards", "create", "labels", "")
	resetSubFlag("cards", "create", "members", "")
	resetSubFlag("cards", "update", "card", "")
	resetSubFlag("cards", "update", "name", "")
	resetSubFlag("cards", "update", "desc", "")
	resetSubFlag("cards", "update", "due", "")
	resetSubFlag("cards", "update", "labels", "")
	resetSubFlag("cards", "update", "members", "")
	resetSubFlag("cards", "move", "card", "")
	resetSubFlag("cards", "move", "list", "")
	resetSubFlag("cards", "move", "pos", "0")
	resetSubFlag("cards", "archive", "card", "")
	resetSubFlag("cards", "delete", "card", "")
	resetSubFlag("comments", "list", "card", "")
	resetSubFlag("comments", "add", "card", "")
	resetSubFlag("comments", "add", "text", "")
	resetSubFlag("comments", "update", "action", "")
	resetSubFlag("comments", "update", "text", "")
	resetSubFlag("comments", "delete", "action", "")
	resetSubFlag("checklists", "list", "card", "")
	resetSubFlag("checklists", "create", "card", "")
	resetSubFlag("checklists", "create", "name", "")
	resetSubFlag("checklists", "delete", "checklist", "")
	resetPathFlag([]string{"checklists", "items", "add"}, "checklist", "")
	resetPathFlag([]string{"checklists", "items", "add"}, "name", "")
	resetPathFlag([]string{"checklists", "items", "update"}, "card", "")
	resetPathFlag([]string{"checklists", "items", "update"}, "item", "")
	resetPathFlag([]string{"checklists", "items", "update"}, "state", "")
	resetPathFlag([]string{"checklists", "items", "delete"}, "checklist", "")
	resetPathFlag([]string{"checklists", "items", "delete"}, "item", "")
	resetSubFlag("attachments", "list", "card", "")
	resetSubFlag("attachments", "add-file", "card", "")
	resetSubFlag("attachments", "add-file", "path", "")
	resetSubFlag("attachments", "add-file", "name", "")
	resetSubFlag("attachments", "add-url", "card", "")
	resetSubFlag("attachments", "add-url", "url", "")
	resetSubFlag("attachments", "add-url", "name", "")
	resetSubFlag("attachments", "delete", "card", "")
	resetSubFlag("attachments", "delete", "attachment", "")
	resetSubFlag("labels", "list", "board", "")
	resetSubFlag("labels", "create", "board", "")
	resetSubFlag("labels", "create", "name", "")
	resetSubFlag("labels", "create", "color", "")
	resetSubFlag("labels", "add", "card", "")
	resetSubFlag("labels", "add", "label", "")
	resetSubFlag("labels", "remove", "card", "")
	resetSubFlag("labels", "remove", "label", "")
	resetSubFlag("members", "list", "board", "")
	resetSubFlag("members", "add", "card", "")
	resetSubFlag("members", "add", "member", "")
	resetSubFlag("members", "remove", "card", "")
	resetSubFlag("members", "remove", "member", "")
	resetSubFlag("search", "cards", "query", "")
	resetSubFlag("search", "boards", "query", "")
	// custom-fields definition subcommands
	resetSubFlag("custom-fields", "list", "board", "")
	resetSubFlag("custom-fields", "get", "field", "")
	resetSubFlag("custom-fields", "create", "board", "")
	resetSubFlag("custom-fields", "create", "name", "")
	resetSubFlag("custom-fields", "create", "type", "")
	resetSubFlag("custom-fields", "create", "card-front", "false")
	// Reset StringArray flag for custom-fields create --option
	if cfCreate, _, err := rootCmd.Find([]string{"custom-fields", "create"}); err == nil {
		if f := cfCreate.Flags().Lookup("option"); f != nil {
			if r, ok := f.Value.(interface{ Replace([]string) error }); ok {
				r.Replace(nil)
			}
			f.Changed = false
		}
	}
	resetSubFlag("custom-fields", "update", "field", "")
	resetSubFlag("custom-fields", "update", "name", "")
	resetSubFlag("custom-fields", "update", "card-front", "false")
	resetSubFlag("custom-fields", "delete", "field", "")
	// custom-fields options subcommands
	resetPathFlag([]string{"custom-fields", "options", "list"}, "field", "")
	resetPathFlag([]string{"custom-fields", "options", "add"}, "field", "")
	resetPathFlag([]string{"custom-fields", "options", "add"}, "text", "")
	resetPathFlag([]string{"custom-fields", "options", "add"}, "color", "")
	resetPathFlag([]string{"custom-fields", "options", "update"}, "field", "")
	resetPathFlag([]string{"custom-fields", "options", "update"}, "option", "")
	resetPathFlag([]string{"custom-fields", "options", "update"}, "text", "")
	resetPathFlag([]string{"custom-fields", "options", "update"}, "color", "")
	resetPathFlag([]string{"custom-fields", "options", "delete"}, "field", "")
	resetPathFlag([]string{"custom-fields", "options", "delete"}, "option", "")
	// custom-fields items subcommands
	resetPathFlag([]string{"custom-fields", "items", "list"}, "card", "")
	resetPathFlag([]string{"custom-fields", "items", "set"}, "card", "")
	resetPathFlag([]string{"custom-fields", "items", "set"}, "field", "")
	resetPathFlag([]string{"custom-fields", "items", "set"}, "text", "")
	resetPathFlag([]string{"custom-fields", "items", "set"}, "number", "")
	resetPathFlag([]string{"custom-fields", "items", "set"}, "date", "")
	resetPathFlag([]string{"custom-fields", "items", "set"}, "checked", "")
	resetPathFlag([]string{"custom-fields", "items", "set"}, "option", "")
	resetPathFlag([]string{"custom-fields", "items", "clear"}, "card", "")
	resetPathFlag([]string{"custom-fields", "items", "clear"}, "field", "")
	_ = resetFlag
}

// assertContractCode validates that err is a *ContractError with the expected code.
func assertContractCode(t *testing.T, err error, want string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error with code %s, got nil", want)
	}
	ce, ok := err.(*contract.ContractError)
	if !ok {
		t.Fatalf("expected *ContractError, got %T: %v", err, err)
	}
	if ce.Code != want {
		t.Fatalf("error code = %s, want %s", ce.Code, want)
	}
}

// executeRootArgs executes the root command with the given arguments and returns the error.
func executeRootArgs(args ...string) error {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}
