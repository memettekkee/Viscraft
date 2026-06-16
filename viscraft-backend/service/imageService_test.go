package service

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"viscraft-backend/constant"
)

func TestHashPrompt_Deterministic(t *testing.T) {
	svc := &ImageService{}

	// Requirement 7.3: identical inputs always produce the same 64-char hex string
	result1 := svc.hashPrompt("dark armor knight", "fantasy", "character", "dark")
	result2 := svc.hashPrompt("dark armor knight", "fantasy", "character", "dark")

	if result1 != result2 {
		t.Errorf("hashPrompt is not deterministic: got %q and %q", result1, result2)
	}
}

func TestHashPrompt_Returns64CharHex(t *testing.T) {
	svc := &ImageService{}

	result := svc.hashPrompt("dark armor knight", "fantasy", "character", "dark")

	if len(result) != 64 {
		t.Errorf("expected 64-char hex string, got length %d: %q", len(result), result)
	}

	// Verify it's valid hex
	_, err := hex.DecodeString(result)
	if err != nil {
		t.Errorf("result is not valid hex: %v", err)
	}
}

func TestHashPrompt_PipeDelimiter(t *testing.T) {
	svc := &ImageService{}

	// Requirement 7.1: hash from (prompt, genre, assetType, mood) joined by pipe delimiter
	// Manually compute expected hash
	input := "dark armor knight|fantasy|character|dark"
	expected := sha256.Sum256([]byte(input))
	expectedHex := hex.EncodeToString(expected[:])

	result := svc.hashPrompt("dark armor knight", "fantasy", "character", "dark")

	if result != expectedHex {
		t.Errorf("hashPrompt does not use pipe delimiter correctly\nexpected: %s\ngot:      %s", expectedHex, result)
	}
}

func TestHashPrompt_DifferentInputsDifferentHashes(t *testing.T) {
	svc := &ImageService{}

	hash1 := svc.hashPrompt("dark armor knight", "fantasy", "character", "dark")
	hash2 := svc.hashPrompt("bright elf mage", "fantasy", "character", "epic")

	if hash1 == hash2 {
		t.Error("different inputs produced the same hash")
	}
}

func TestHashPrompt_OrderMatters(t *testing.T) {
	svc := &ImageService{}

	hash1 := svc.hashPrompt("prompt", "genre", "assetType", "mood")
	hash2 := svc.hashPrompt("genre", "prompt", "assetType", "mood")

	if hash1 == hash2 {
		t.Error("different ordering produced the same hash, order should matter")
	}
}

func TestHashPrompt_SinglePart(t *testing.T) {
	svc := &ImageService{}

	result := svc.hashPrompt("single")
	expected := sha256.Sum256([]byte("single"))
	expectedHex := hex.EncodeToString(expected[:])

	if result != expectedHex {
		t.Errorf("single part hash mismatch\nexpected: %s\ngot:      %s", expectedHex, result)
	}
}

func TestHashPrompt_EmptyParts(t *testing.T) {
	svc := &ImageService{}

	// Should still produce a valid 64-char hex string even with empty parts
	result := svc.hashPrompt("", "", "", "")

	if len(result) != 64 {
		t.Errorf("expected 64-char hex string for empty parts, got length %d", len(result))
	}

	// Verify it matches expected pipe-delimited empty string hash
	expected := sha256.Sum256([]byte("|||"))
	expectedHex := hex.EncodeToString(expected[:])
	if result != expectedHex {
		t.Errorf("empty parts hash mismatch\nexpected: %s\ngot:      %s", expectedHex, result)
	}
}

func TestValidatePrompt_ValidPrompts(t *testing.T) {
	svc := &ImageService{}

	tests := []struct {
		name   string
		prompt string
	}{
		{"minimum length (3 chars)", "abc"},
		{"maximum length (300 chars)", strings.Repeat("a", 300)},
		{"normal prompt", "a dark knight in a forest"},
		{"with leading/trailing whitespace", "   valid prompt   "},
		{"exactly 3 chars after trim", "  abc  "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.validatePrompt(tt.prompt)
			if err != nil {
				t.Errorf("expected nil error for prompt %q, got %v", tt.prompt, err)
			}
		})
	}
}

func TestValidatePrompt_TooShort(t *testing.T) {
	svc := &ImageService{}

	tests := []struct {
		name   string
		prompt string
	}{
		{"empty string", ""},
		{"single char", "a"},
		{"two chars", "ab"},
		{"only whitespace", "   "},
		{"two chars with whitespace", "  ab  "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.validatePrompt(tt.prompt)
			if err == nil {
				t.Errorf("expected error for prompt %q, got nil", tt.prompt)
			}
			if err != nil && err.Code != constant.ErrInvalidPrompt.Code {
				t.Errorf("expected error code %s, got %s", constant.ErrInvalidPrompt.Code, err.Code)
			}
		})
	}
}

func TestValidatePrompt_TooLong(t *testing.T) {
	svc := &ImageService{}

	tests := []struct {
		name   string
		prompt string
	}{
		{"301 chars", strings.Repeat("a", 301)},
		{"500 chars", strings.Repeat("b", 500)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.validatePrompt(tt.prompt)
			if err == nil {
				t.Errorf("expected error for prompt of length %d, got nil", len(tt.prompt))
			}
			if err != nil && err.Code != constant.ErrInvalidPrompt.Code {
				t.Errorf("expected error code %s, got %s", constant.ErrInvalidPrompt.Code, err.Code)
			}
		})
	}
}

func TestValidatePrompt_BlockedWords(t *testing.T) {
	svc := &ImageService{}

	tests := []struct {
		name   string
		prompt string
	}{
		{"contains nude", "a nude character"},
		{"contains explicit", "explicit content here"},
		{"contains nsfw", "this is nsfw material"},
		{"contains gore", "gore scene in forest"},
		{"uppercase blocked word", "NUDE warrior"},
		{"mixed case blocked word", "ExPlIcIt violence"},
		{"blocked word as substring", "denude the tree"},
		{"blocked word in middle", "something nsfwork related"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.validatePrompt(tt.prompt)
			if err == nil {
				t.Errorf("expected error for prompt %q, got nil", tt.prompt)
			}
			if err != nil && err.Code != constant.ErrInvalidPrompt.Code {
				t.Errorf("expected error code %s, got %s", constant.ErrInvalidPrompt.Code, err.Code)
			}
		})
	}
}

func TestValidatePrompt_WhitespaceTrimming(t *testing.T) {
	svc := &ImageService{}

	// "ab" after trim is only 2 chars - should fail
	err := svc.validatePrompt("   ab   ")
	if err == nil {
		t.Error("expected error for prompt with only 2 chars after trim")
	}

	// "abc" after trim is 3 chars - should pass
	err = svc.validatePrompt("   abc   ")
	if err != nil {
		t.Errorf("expected nil for prompt with 3 chars after trim, got %v", err)
	}
}
