package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type callRoutingAudioFileResponse struct {
	ObjectKey     string  `json:"objectKey"`
	FileName      *string `json:"fileName"`
	ContentType   *string `json:"contentType"`
	FileSizeBytes *int64  `json:"fileSizeBytes"`
}

func (c *AllQuietAPIClient) UploadCallRoutingMedia(ctx context.Context, integrationId, purpose, filePath string) (*callRoutingAudioFileResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open audio file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("create multipart form: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("copy audio file: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	path := fmt.Sprintf("/inbound-integration/%s/call-routing-media/%s", url.PathEscape(integrationId), url.PathEscape(purpose))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.EndpointURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Authorization", c.APIKey)

	httpResp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result callRoutingAudioFileResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
