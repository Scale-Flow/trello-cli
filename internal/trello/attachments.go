package trello

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/brettmcdowell/trello-cli/internal/contract"
)

func (c *Client) ListAttachments(ctx context.Context, cardID string) ([]Attachment, error) {
	var attachments []Attachment
	err := c.Get(ctx, fmt.Sprintf("/1/cards/%s/attachments", cardID), nil, &attachments)
	return attachments, err
}

func (c *Client) AddURLAttachment(ctx context.Context, cardID, urlStr string, name *string) (Attachment, error) {
	queryParams := map[string]string{"url": urlStr}
	if name != nil {
		queryParams["name"] = *name
	}
	var attachment Attachment
	err := c.Post(ctx, fmt.Sprintf("/1/cards/%s/attachments", cardID), queryParams, &attachment)
	return attachment, err
}

func (c *Client) AddFileAttachment(ctx context.Context, cardID, filePath string, name *string) (Attachment, error) {
	queryParams := map[string]string{}
	if name != nil {
		queryParams["name"] = *name
	}
	var attachment Attachment
	err := c.postMultipartFile(ctx, fmt.Sprintf("/1/cards/%s/attachments", cardID), filePath, queryParams, &attachment)
	return attachment, err
}

func (c *Client) DeleteAttachment(ctx context.Context, cardID, attachmentID string) error {
	return c.Delete(ctx, fmt.Sprintf("/1/cards/%s/attachments/%s", cardID, attachmentID), nil)
}

// postMultipartFile handles multipart/form-data file uploads.
func (c *Client) postMultipartFile(ctx context.Context, path, filePath string, params map[string]string, result any) error {
	file, err := os.Open(filePath)
	if err != nil {
		return contract.NewError(contract.FileNotFound, fmt.Sprintf("cannot open file: %s", filePath))
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return contract.NewError(contract.UnknownError, fmt.Sprintf("failed to create form file: %v", err))
	}
	if _, err := io.Copy(part, file); err != nil {
		return contract.NewError(contract.UnknownError, fmt.Sprintf("failed to read file: %v", err))
	}
	for k, v := range params {
		if err := writer.WriteField(k, v); err != nil {
			return contract.NewError(contract.UnknownError, fmt.Sprintf("failed to write form field: %v", err))
		}
	}
	if err := writer.Close(); err != nil {
		return contract.NewError(contract.UnknownError, fmt.Sprintf("failed to finalize multipart body: %v", err))
	}

	fullURL := c.buildURL(path, nil)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return contract.NewError(contract.HTTPError, fmt.Sprintf("upload failed: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return mapHTTPError(resp)
	}
	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}
	return nil
}
