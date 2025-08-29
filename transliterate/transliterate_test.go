package transliterate

import (
	"context"
	"strings"
	"testing"
)

// Run tests using `encore test`, which compiles the Encore app and then runs `go test`.
// It supports all the same flags that the `go test` command does.
// You automatically get tracing for tests in the local dev dash: http://localhost:9400
// Learn more: https://encore.dev/docs/go/develop/testing

// TestTransliterateAndRetrieve tests that transliteration is performed and stored correctly
func TestTransliterateAndRetrieve(t *testing.T) {
	testText := "Привет"
	req := TransliterationRequest{
		Text:         testText,
		InputScript:  "cyrillic",
		OutputScript: "latin",
	}
	
	resp, err := Transliterate(context.Background(), &req)
	if err != nil {
		t.Fatal(err)
	}
	
	if resp.InputText != testText {
		t.Errorf("got input %q, want %q", resp.InputText, testText)
	}
	
	if resp.InputScript != "cyrillic" {
		t.Errorf("got input script %q, want %q", resp.InputScript, "cyrillic")
	}
	
	if resp.OutputScript != "latin" {
		t.Errorf("got output script %q, want %q", resp.OutputScript, "latin")
	}
	
	if resp.ConfidenceScore == nil || *resp.ConfidenceScore <= 0 {
		t.Errorf("expected positive confidence score, got %v", resp.ConfidenceScore)
	}

	// Verify actual transliteration occurred
	expectedOutput := "Privet" // Cyrillic "Привет" -> Latin "Privet"
	if resp.OutputText != expectedOutput {
		t.Errorf("got output %q, want %q", resp.OutputText, expectedOutput)
	}

	// Test retrieval by ID
	retrieved, err := GetTransliteration(context.Background(), resp.ID)
	if err != nil {
		t.Fatal(err)
	}
	
	if retrieved.ID != resp.ID {
		t.Errorf("got ID %q, want %q", retrieved.ID, resp.ID)
	}
	
	if retrieved.InputText != resp.InputText {
		t.Errorf("got input text %q, want %q", retrieved.InputText, resp.InputText)
	}
}

// TestScriptDetection tests automatic script detection
func TestScriptDetection(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{"Cyrillic", "Привет мир", "cyrillic"},
		{"Chinese", "你好世界", "chinese"},
		{"Arabic", "مرحبا بالعالم", "arabic"},
		{"Greek", "Γεια σας κόσμος", "greek"},
		{"Latin", "Hello world", "latin"},
		{"Mixed favour Latin", "Hello мир", "unknown"}, // Mixed should return unknown
		{"Empty string", "", "unknown"},
		{"Numbers only", "12345", "unknown"},
		{"Spaces only", "   ", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectScript(tt.text)
			if result != tt.expected {
				t.Errorf("detectScript(%q) = %q, want %q", tt.text, result, tt.expected)
			}
		})
	}
}

// TestTransliterationBuiltinRules tests built-in transliteration rules
func TestTransliterationBuiltinRules(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		inputScript  string
		outputScript string
		expectedPart string // Part of expected output
	}{
		{"Cyrillic to Latin", "Привет", "cyrillic", "latin", "Privet"},
		{"Chinese to Latin", "你好", "chinese", "latin", "ni"},
		{"Greek to Latin", "Γεια", "greek", "latin", "G"},
		{"Arabic to Latin", "مرحبا", "arabic", "latin", "m"},
		{"Latin to ASCII", "café", "latin", "ascii", "cafe"}, // Accented chars removed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := performTransliteration(tt.input, tt.inputScript, tt.outputScript, nil)
			if !strings.Contains(result, tt.expectedPart) {
				t.Errorf("performTransliteration(%q, %q, %q) = %q, expected to contain %q",
					tt.input, tt.inputScript, tt.outputScript, result, tt.expectedPart)
			}
		})
	}
}

// TestValidation tests input validation
func TestValidation(t *testing.T) {
	tests := []struct {
		name        string
		req         *TransliterationRequest
		expectError bool
	}{
		{
			name: "Valid request",
			req: &TransliterationRequest{
				Text:         "Hello",
				OutputScript: "ascii",
			},
			expectError: false,
		},
		{
			name:        "Nil request",
			req:         nil,
			expectError: true,
		},
		{
			name: "Empty text",
			req: &TransliterationRequest{
				Text:         "",
				OutputScript: "ascii",
			},
			expectError: true,
		},
		{
			name: "No output script",
			req: &TransliterationRequest{
				Text: "Hello",
			},
			expectError: true,
		},
		{
			name: "Invalid output script",
			req: &TransliterationRequest{
				Text:         "Hello",
				OutputScript: "klingon",
			},
			expectError: true,
		},
		{
			name: "Too long text",
			req: &TransliterationRequest{
				Text:         strings.Repeat("x", 10001),
				OutputScript: "ascii",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTransliterationRequest(tt.req)
			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestFeedback tests the feedback submission functionality
func TestFeedback(t *testing.T) {
	// First create a transliteration
	req := TransliterationRequest{
		Text:         "Test",
		InputScript:  "latin",
		OutputScript: "ascii",
	}
	
	resp, err := Transliterate(context.Background(), &req)
	if err != nil {
		t.Fatal(err)
	}

	// Submit feedback
	feedbackReq := FeedbackRequest{
		TransliterationID: resp.ID,
		SuggestedOutput:   "Better Output",
		FeedbackType:      "correction",
		UserContext:       "Test feedback",
	}
	
	err = SubmitFeedback(context.Background(), resp.ID, &feedbackReq)
	if err != nil {
		t.Errorf("SubmitFeedback failed: %v", err)
	}
}

// TestFeedbackValidation tests feedback validation
func TestFeedbackValidation(t *testing.T) {
	tests := []struct {
		name        string
		req         *FeedbackRequest
		expectError bool
	}{
		{
			name: "Valid feedback",
			req: &FeedbackRequest{
				SuggestedOutput: "Better output",
				FeedbackType:    "correction",
			},
			expectError: false,
		},
		{
			name:        "Nil request",
			req:         nil,
			expectError: true,
		},
		{
			name: "Empty suggested output",
			req: &FeedbackRequest{
				SuggestedOutput: "",
				FeedbackType:    "correction",
			},
			expectError: true,
		},
		{
			name: "Invalid feedback type",
			req: &FeedbackRequest{
				SuggestedOutput: "Better output",
				FeedbackType:    "invalid",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFeedbackRequest(tt.req)
			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestConfidenceCalculation tests confidence score calculation
func TestConfidenceCalculation(t *testing.T) {
	tests := []struct {
		name         string
		inputText    string
		outputText   string
		inputScript  string
		outputScript string
		expectedMin  float64
		expectedMax  float64
	}{
		{"Latin to ASCII", "hello", "hello", "latin", "ascii", 0.8, 1.0},
		{"Chinese to Latin", "你好", "ni hao", "chinese", "latin", 0.5, 0.8},
		{"Cyrillic to Latin", "привет", "privet", "cyrillic", "latin", 0.6, 0.9},
		{"Empty output", "test", "", "latin", "ascii", 0.1, 0.1},
		{"Reasonable length preservation", "hello", "world", "latin", "ascii", 0.7, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			confidence := calculateConfidence(tt.inputText, tt.outputText, tt.inputScript, tt.outputScript)
			if confidence < tt.expectedMin || confidence > tt.expectedMax {
				t.Errorf("calculateConfidence(%q, %q, %q, %q) = %f, want between %f and %f",
					tt.inputText, tt.outputText, tt.inputScript, tt.outputScript, 
					confidence, tt.expectedMin, tt.expectedMax)
			}
		})
	}
}

// TestUUIDValidation tests UUID format validation
func TestUUIDValidation(t *testing.T) {
	tests := []struct {
		name     string
		uuid     string
		expected bool
	}{
		{"Valid UUID", "550e8400-e29b-41d4-a716-446655440000", true},
		{"Valid UUID uppercase", "550E8400-E29B-41D4-A716-446655440000", true},
		{"Invalid - too short", "550e8400-e29b-41d4-a716", false},
		{"Invalid - no hyphens", "550e8400e29b41d4a716446655440000", false},
		{"Invalid - wrong hyphen positions", "550e8400e-29b-41d4-a716-446655440000", false},
		{"Invalid - non-hex characters", "550e8400-e29b-41d4-a716-44665544000g", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidUUID(tt.uuid)
			if result != tt.expected {
				t.Errorf("isValidUUID(%q) = %v, want %v", tt.uuid, result, tt.expected)
			}
		})
	}
}

// TestLocaleValidation tests locale format validation
func TestLocaleValidation(t *testing.T) {
	tests := []struct {
		name     string
		locale   string
		expected bool
	}{
		{"Valid language only", "en", true},
		{"Valid language-country", "en-US", true},
		{"Valid 3-letter language", "chi", true},
		{"Invalid - empty", "", false},
		{"Invalid - too short", "e", false},
		{"Invalid - uppercase language", "EN", false},
		{"Invalid - lowercase country", "en-us", false},
		{"Invalid - too many parts", "en-US-variant", false},
		{"Invalid - wrong country format", "en-USA", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidLocale(tt.locale)
			if result != tt.expected {
				t.Errorf("isValidLocale(%q) = %v, want %v", tt.locale, result, tt.expected)
			}
		})
	}
}

// TestScriptPairSupport tests supported script combinations
func TestScriptPairSupport(t *testing.T) {
	tests := []struct {
		name         string
		inputScript  string
		outputScript string
		expected     bool
	}{
		{"Latin to ASCII", "latin", "ascii", true},
		{"Cyrillic to Latin", "cyrillic", "latin", true},
		{"Chinese to Latin", "chinese", "latin", true},
		{"Unsupported - Latin to Chinese", "latin", "chinese", false},
		{"Unsupported - Unknown script", "klingon", "latin", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSupportedScriptPair(tt.inputScript, tt.outputScript)
			if result != tt.expected {
				t.Errorf("isSupportedScriptPair(%q, %q) = %v, want %v", 
					tt.inputScript, tt.outputScript, result, tt.expected)
			}
		})
	}
}

// TestCaching tests that identical requests are cached
func TestCaching(t *testing.T) {
	req := TransliterationRequest{
		Text:         "Hello",
		InputScript:  "latin",
		OutputScript: "ascii",
	}
	
	// First request
	resp1, err := Transliterate(context.Background(), &req)
	if err != nil {
		t.Fatal(err)
	}
	
	// Second identical request should return cached result
	resp2, err := Transliterate(context.Background(), &req)
	if err != nil {
		t.Fatal(err)
	}
	
	// Should be the same result (cached)
	if resp1.ID != resp2.ID {
		t.Error("expected cached result to have same ID")
	}
	
	if resp1.OutputText != resp2.OutputText {
		t.Error("expected cached result to have same output")
	}
}

// TestAutoScriptDetection tests transliteration with auto-detection
func TestAutoScriptDetection(t *testing.T) {
	req := TransliterationRequest{
		Text:         "Привет", // Don't specify InputScript - should auto-detect
		OutputScript: "latin",
	}
	
	resp, err := Transliterate(context.Background(), &req)
	if err != nil {
		t.Fatal(err)
	}
	
	if resp.InputScript != "cyrillic" {
		t.Errorf("expected auto-detected script to be 'cyrillic', got %q", resp.InputScript)
	}
	
	if !strings.Contains(resp.OutputText, "Privet") {
		t.Errorf("expected transliteration to contain 'Privet', got %q", resp.OutputText)
	}
}