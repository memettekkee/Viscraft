package gemini

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"viscraft-backend/constant"
)

// newTestClient creates a Client pointing at the given test server URL.
func newTestClient(serverURL string) *Client {
	return &Client{
		apiKey:  "test-api-key",
		model:   "gemini-2.0-flash-exp-image-generation",
		baseURL: serverURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func TestNewClient(t *testing.T) {
	client := NewClient("my-key", "my-model")
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.apiKey != "my-key" {
		t.Errorf("expected apiKey 'my-key', got '%s'", client.apiKey)
	}
	if client.model != "my-model" {
		t.Errorf("expected model 'my-model', got '%s'", client.model)
	}
	if client.httpClient == nil {
		t.Fatal("expected non-nil httpClient")
	}
	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("expected 30s timeout, got %v", client.httpClient.Timeout)
	}
	if client.baseURL != defaultBaseURL {
		t.Errorf("expected baseURL %q, got %q", defaultBaseURL, client.baseURL)
	}
}

func TestGenerate_Success(t *testing.T) {
	imageData := []byte("fake-png-image-data")
	b64Data := base64.StdEncoding.EncodeToString(imageData)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and content type
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", ct)
		}

		// Verify URL path contains model and action
		if !strings.Contains(r.URL.Path, "generateContent") {
			t.Errorf("expected URL to contain generateContent, got %s", r.URL.Path)
		}

		// Verify API key is passed
		if key := r.URL.Query().Get("key"); key != "test-api-key" {
			t.Errorf("expected key=test-api-key, got key=%s", key)
		}

		// Verify request body structure
		var reqBody generateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if len(reqBody.Contents) != 1 || len(reqBody.Contents[0].Parts) != 1 {
			t.Fatal("unexpected request body structure")
		}
		if reqBody.Contents[0].Parts[0].Text != "a dark fantasy knight" {
			t.Errorf("unexpected prompt: %s", reqBody.Contents[0].Parts[0].Text)
		}
		if len(reqBody.GenerationConfig.ResponseModalities) != 1 || reqBody.GenerationConfig.ResponseModalities[0] != "IMAGE" {
			t.Error("expected responseModalities to contain IMAGE")
		}

		resp := generateResponse{
			Candidates: []candidate{
				{
					Content: candidateContent{
						Parts: []responsePart{
							{
								InlineData: &inlineData{
									MimeType: "image/png",
									Data:     b64Data,
								},
							},
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.Generate(context.Background(), "a dark fantasy knight")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if string(result) != string(imageData) {
		t.Errorf("expected image data %q, got %q", string(imageData), string(result))
	}
}

func TestGenerate_Non200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal error"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.Generate(context.Background(), "test prompt")
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
	if !strings.Contains(err.Error(), constant.ErrGeminiBadResponse.Message) {
		t.Errorf("expected error to contain %q, got %q", constant.ErrGeminiBadResponse.Message, err.Error())
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to mention status code 500, got %q", err.Error())
	}
}

func TestGenerate_EmptyImageData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := generateResponse{
			Candidates: []candidate{
				{
					Content: candidateContent{
						Parts: []responsePart{
							{Text: "I cannot generate that image"},
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.Generate(context.Background(), "test prompt")
	if err == nil {
		t.Fatal("expected error for empty image data")
	}
	if !strings.Contains(err.Error(), constant.ErrGeminiBadResponse.Message) {
		t.Errorf("expected error to contain %q, got %q", constant.ErrGeminiBadResponse.Message, err.Error())
	}
}

func TestGenerate_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`not valid json`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.Generate(context.Background(), "test prompt")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), constant.ErrGeminiBadResponse.Message) {
		t.Errorf("expected error to contain %q, got %q", constant.ErrGeminiBadResponse.Message, err.Error())
	}
}

func TestGenerate_ContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &Client{
		apiKey:  "test-key",
		model:   "test-model",
		baseURL: server.URL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	// Use a context with very short deadline
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := client.Generate(ctx, "test prompt")
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if !strings.Contains(err.Error(), constant.ErrGeminiTimeout.Message) {
		t.Errorf("expected error to contain %q, got %q", constant.ErrGeminiTimeout.Message, err.Error())
	}
}

func TestGenerate_InvalidBase64(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := generateResponse{
			Candidates: []candidate{
				{
					Content: candidateContent{
						Parts: []responsePart{
							{
								InlineData: &inlineData{
									MimeType: "image/png",
									Data:     "!!!not-valid-base64!!!",
								},
							},
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.Generate(context.Background(), "test prompt")
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
	if !strings.Contains(err.Error(), constant.ErrGeminiBadResponse.Message) {
		t.Errorf("expected error to contain %q, got %q", constant.ErrGeminiBadResponse.Message, err.Error())
	}
}

func TestExtractImageBytes_Success(t *testing.T) {
	imageData := []byte("test-image-bytes")
	b64Data := base64.StdEncoding.EncodeToString(imageData)

	resp := generateResponse{
		Candidates: []candidate{
			{
				Content: candidateContent{
					Parts: []responsePart{
						{
							InlineData: &inlineData{
								MimeType: "image/png",
								Data:     b64Data,
							},
						},
					},
				},
			},
		},
	}

	result, err := extractImageBytes(resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != string(imageData) {
		t.Errorf("expected %q, got %q", string(imageData), string(result))
	}
}

func TestExtractImageBytes_NoCandidates(t *testing.T) {
	resp := generateResponse{Candidates: []candidate{}}
	_, err := extractImageBytes(resp)
	if err == nil {
		t.Fatal("expected error for empty candidates")
	}
}

func TestExtractImageBytes_NoInlineData(t *testing.T) {
	resp := generateResponse{
		Candidates: []candidate{
			{
				Content: candidateContent{
					Parts: []responsePart{
						{Text: "some text response"},
					},
				},
			},
		},
	}
	_, err := extractImageBytes(resp)
	if err == nil {
		t.Fatal("expected error for missing inlineData")
	}
}

func TestExtractImageBytes_EmptyData(t *testing.T) {
	resp := generateResponse{
		Candidates: []candidate{
			{
				Content: candidateContent{
					Parts: []responsePart{
						{
							InlineData: &inlineData{
								MimeType: "image/png",
								Data:     "",
							},
						},
					},
				},
			},
		},
	}
	_, err := extractImageBytes(resp)
	if err == nil {
		t.Fatal("expected error for empty Data field")
	}
}
