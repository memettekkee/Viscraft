package imagegen

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"viscraft-backend/constant"
)

const (
	pollinationsImageURL = "https://gen.pollinations.ai/image"
	pollinationsEditsURL = "https://gen.pollinations.ai/v1/images/edits"
)

// Client wraps the Pollinations image generation API.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// Generate calls Pollinations API.
// If referenceImage is nil → text-to-image (GET).
// If referenceImage is provided → image editing (POST multipart).
func (c *Client) Generate(ctx context.Context, prompt string, referenceImage []byte) ([]byte, error) {
	if referenceImage == nil {
		return c.generateTextOnly(ctx, prompt)
	}
	return c.generateWithReference(ctx, prompt, referenceImage)
}

func (c *Client) generateTextOnly(ctx context.Context, prompt string) ([]byte, error) {
	endpoint := pollinationsImageURL + "/" + url.PathEscape(prompt) + "?model=kontext&nologo=true"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", constant.ErrGeminiBadResponse.Message, err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("%s: %w", constant.ErrGeminiTimeout.Message, err)
		}
		return nil, fmt.Errorf("%s: %w", constant.ErrGeminiBadResponse.Message, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyPreview, _ := io.ReadAll(io.LimitReader(resp.Body, 500))
		bodyStr := string(bodyPreview)
		if resp.StatusCode == http.StatusBadRequest && (strings.Contains(bodyStr, "safety_error") || strings.Contains(bodyStr, "safety")) {
			return nil, fmt.Errorf("content_policy_violation: prompt was rejected by safety filter")
		}
		return nil, fmt.Errorf("%s: API returned status %d, body: %s", constant.ErrGeminiBadResponse.Message, resp.StatusCode, bodyStr)
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) generateWithReference(ctx context.Context, prompt string, referenceImage []byte) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("prompt", prompt)
	_ = writer.WriteField("model", "kontext")

	part, err := writer.CreateFormFile("image", "reference.png")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", constant.ErrGeminiBadResponse.Message, err)
	}
	if _, err := part.Write(referenceImage); err != nil {
		return nil, fmt.Errorf("%s: %w", constant.ErrGeminiBadResponse.Message, err)
	}
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, pollinationsEditsURL, body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", constant.ErrGeminiBadResponse.Message, err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("%s: %w", constant.ErrGeminiTimeout.Message, err)
		}
		return nil, fmt.Errorf("%s: %w", constant.ErrGeminiBadResponse.Message, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyPreview, _ := io.ReadAll(io.LimitReader(resp.Body, 500))
		bodyStr := string(bodyPreview)
		if resp.StatusCode == http.StatusBadRequest && (strings.Contains(bodyStr, "safety_error") || strings.Contains(bodyStr, "safety")) {
			return nil, fmt.Errorf("content_policy_violation: prompt was rejected by safety filter")
		}
		return nil, fmt.Errorf("%s: API returned status %d, body: %s", constant.ErrGeminiBadResponse.Message, resp.StatusCode, bodyStr)
	}

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", constant.ErrGeminiBadResponse.Message, err)
	}

	// Handle JSON response (base64 or URL)
	if len(imageBytes) > 0 && imageBytes[0] == '{' {
		return c.extractImageFromJSON(ctx, imageBytes)
	}

	return imageBytes, nil
}

func (c *Client) extractImageFromJSON(ctx context.Context, jsonData []byte) ([]byte, error) {
	type editsResponse struct {
		Output []string `json:"output"`
		Data   []struct {
			URL     string `json:"url"`
			B64JSON string `json:"b64_json"`
		} `json:"data"`
	}

	var parsed editsResponse
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		return nil, fmt.Errorf("%s: unexpected JSON response", constant.ErrGeminiBadResponse.Message)
	}

	if len(parsed.Data) > 0 && parsed.Data[0].B64JSON != "" {
		decoded, err := base64.StdEncoding.DecodeString(parsed.Data[0].B64JSON)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to decode b64_json: %w", constant.ErrGeminiBadResponse.Message, err)
		}
		return decoded, nil
	}

	var imageURL string
	if len(parsed.Output) > 0 {
		imageURL = parsed.Output[0]
	} else if len(parsed.Data) > 0 && parsed.Data[0].URL != "" {
		imageURL = parsed.Data[0].URL
	}

	if imageURL == "" {
		return nil, fmt.Errorf("%s: no image data in JSON response", constant.ErrGeminiBadResponse.Message)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", constant.ErrGeminiBadResponse.Message, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", constant.ErrGeminiBadResponse.Message, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: image download returned status %d", constant.ErrGeminiBadResponse.Message, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
