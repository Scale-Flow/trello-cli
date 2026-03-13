package trello

import (
	"context"
	"fmt"
)

func (c *Client) ListCustomFieldsByBoard(ctx context.Context, boardID string) ([]CustomField, error) {
	path := fmt.Sprintf("/1/boards/%s/customFields", boardID)
	var fields []CustomField
	err := c.Get(ctx, path, nil, &fields)
	return fields, err
}

func (c *Client) GetCustomField(ctx context.Context, fieldID string) (CustomField, error) {
	path := fmt.Sprintf("/1/customFields/%s", fieldID)
	var field CustomField
	err := c.Get(ctx, path, nil, &field)
	return field, err
}

func (c *Client) CreateCustomField(ctx context.Context, params CreateCustomFieldParams) (CustomField, error) {
	var field CustomField
	err := c.PostJSON(ctx, "/1/customFields", params, &field)
	return field, err
}

func (c *Client) UpdateCustomField(ctx context.Context, fieldID string, params UpdateCustomFieldParams) (CustomField, error) {
	path := fmt.Sprintf("/1/customFields/%s", fieldID)
	var field CustomField
	err := c.PutJSON(ctx, path, params, &field)
	return field, err
}

func (c *Client) DeleteCustomField(ctx context.Context, fieldID string) error {
	path := fmt.Sprintf("/1/customFields/%s", fieldID)
	return c.Delete(ctx, path, nil)
}

func (c *Client) ListCustomFieldOptions(ctx context.Context, fieldID string) ([]CustomFieldOption, error) {
	path := fmt.Sprintf("/1/customFields/%s/options", fieldID)
	var options []CustomFieldOption
	err := c.Get(ctx, path, nil, &options)
	return options, err
}

func (c *Client) CreateCustomFieldOption(ctx context.Context, fieldID string, params CreateCustomFieldOptionParams) (CustomFieldOption, error) {
	path := fmt.Sprintf("/1/customFields/%s/options", fieldID)
	var option CustomFieldOption
	err := c.PostJSON(ctx, path, params, &option)
	return option, err
}

func (c *Client) UpdateCustomFieldOption(ctx context.Context, fieldID, optionID string, params UpdateCustomFieldOptionParams) (CustomFieldOption, error) {
	path := fmt.Sprintf("/1/customFields/%s/options/%s", fieldID, optionID)
	var option CustomFieldOption
	err := c.PutJSON(ctx, path, params, &option)
	return option, err
}

func (c *Client) DeleteCustomFieldOption(ctx context.Context, fieldID, optionID string) error {
	path := fmt.Sprintf("/1/customFields/%s/options/%s", fieldID, optionID)
	return c.Delete(ctx, path, nil)
}

func (c *Client) ListCardCustomFieldItems(ctx context.Context, cardID string) ([]CardCustomFieldItem, error) {
	path := fmt.Sprintf("/1/cards/%s/customFieldItems", cardID)
	var items []CardCustomFieldItem
	err := c.Get(ctx, path, nil, &items)
	return items, err
}

func (c *Client) SetCardCustomFieldItem(ctx context.Context, cardID, fieldID string, params SetCardCustomFieldItemParams) (CardCustomFieldItem, error) {
	path := fmt.Sprintf("/1/cards/%s/customField/%s/item", cardID, fieldID)
	var item CardCustomFieldItem
	err := c.PutJSON(ctx, path, params, &item)
	if err == nil && item.IDValue != "" && item.Value.IDValue == "" {
		item.Value.IDValue = item.IDValue
	}
	return item, err
}

func (c *Client) ClearCardCustomFieldItem(ctx context.Context, cardID, fieldID string) error {
	path := fmt.Sprintf("/1/cards/%s/customField/%s/item", cardID, fieldID)
	return c.Put(ctx, path, nil, nil)
}
