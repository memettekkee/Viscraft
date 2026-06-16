package gemini

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"viscraft-backend/constant"
	"viscraft-backend/pkg/logger"
)

const defaultBaseURL = "https://generativelanguage.googleapis.com/v1beta/models"

// Client wraps the Google Gemini Flash Image API.
type Client struct {
	apiKey     string
	model      string
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Gemini API client with the given API key and model name.
func NewClient(apiKey, model string) *Client {
	return &Client{
		apiKey:  apiKey,
		model:   model,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// generateRequest represents the Gemini generateContent API request payload.
type generateRequest struct {
	Contents         []content        `json:"contents"`
	GenerationConfig generationConfig `json:"generationConfig"`
}

type content struct {
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text,omitempty"`
}

type generationConfig struct {
	ResponseModalities []string `json:"responseModalities"`
}

// generateResponse represents the Gemini generateContent API response.
type generateResponse struct {
	Candidates []candidate `json:"candidates"`
}

type candidate struct {
	Content candidateContent `json:"content"`
}

type candidateContent struct {
	Parts []responsePart `json:"parts"`
}

type responsePart struct {
	InlineData *inlineData `json:"inlineData,omitempty"`
	Text       string      `json:"text,omitempty"`
}

type inlineData struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

// Generate calls the Gemini Flash Image API with the given prompt and returns
// the generated image bytes. The context is used for timeout propagation.
// Returns structured errors: ErrGeminiTimeout for deadline exceeded,
// ErrGeminiBadResponse for non-200 status or invalid response data.
func (c *Client) Generate(ctx context.Context, prompt string) ([]byte, error) {
	requestID := ""
	if rid, ok := ctx.Value("requestId").(string); ok {
		requestID = rid
	}

	// Build request payload
	reqBody := generateRequest{
		Contents: []content{
			{
				Parts: []part{
					{Text: prompt},
				},
			},
		},
		GenerationConfig: generationConfig{
			ResponseModalities: []string{"IMAGE"},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		logger.Error(requestID, "failed to marshal Gemini request", err)
		return nil, fmt.Errorf("%s: %w", constant.ErrGeminiBadResponse.Message, err)
	}

	// Build HTTP request with context for timeout propagation
	url := fmt.Sprintf("%s/%s:generateContent?key=%s", c.baseURL, c.model, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		logger.Error(requestID, "failed to create Gemini HTTP request", err)
		return nil, fmt.Errorf("%s: %w", constant.ErrGeminiBadResponse.Message, err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			logger.Error(requestID, "Gemini API call timed out", err)
			return nil, fmt.Errorf("%s: %w", constant.ErrGeminiTimeout.Message, err)
		}
		logger.Error(requestID, "Gemini API call failed", err)
		return nil, fmt.Errorf("%s: %w", constant.ErrGeminiBadResponse.Message, err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(requestID, "failed to read Gemini response body", err)
		return nil, fmt.Errorf("%s: failed to read response body", constant.ErrGeminiBadResponse.Message)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		logger.Error(requestID, "Gemini API returned non-200 status",
			"status", resp.StatusCode, "body", string(respBody))
		return nil, fmt.Errorf("%s: API returned status %d", constant.ErrGeminiBadResponse.Message, resp.StatusCode)
	}

	// Parse response
	var genResp generateResponse
	if err := json.Unmarshal(respBody, &genResp); err != nil {
		logger.Error(requestID, "failed to unmarshal Gemini response", err)
		return nil, fmt.Errorf("%s: %w", constant.ErrGeminiBadResponse.Message, err)
	}

	// Extract image data from response
	imageBytes, err := extractImageBytes(genResp)
	if err != nil {
		logger.Error(requestID, "failed to extract image from Gemini response", err)
		return nil, fmt.Errorf("%s: %w", constant.ErrGeminiBadResponse.Message, err)
	}

	return imageBytes, nil
}

// extractImageBytes finds the first inline_data part in the response and decodes
// the base64 image data.
func extractImageBytes(resp generateResponse) ([]byte, error) {
	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			if part.InlineData != nil && part.InlineData.Data != "" {
				decoded, err := base64.StdEncoding.DecodeString(part.InlineData.Data)
				if err != nil {
					return nil, fmt.Errorf("failed to decode base64 image data: %w", err)
				}
				return decoded, nil
			}
		}
	}
	return nil, fmt.Errorf("no image data found in response")
}
