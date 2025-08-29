package transliterate

import (
	"context"
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
		{"Cyrillic", "Привет", "cyrillic"},
		{"Chinese", "你好", "chinese"},
		{"Arabic", "مرحبا", "arabic"},
		{"Latin", "Hello", "latin"},
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

// TestConfidenceCalculation tests confidence score calculation
func TestConfidenceCalculation(t *testing.T) {
	tests := []struct {
		name         string
		inputScript  string
		outputScript string
		expectedMin  float64
		expectedMax  float64
	}{
		{"Latin to ASCII", "latin", "ascii", 0.90, 1.0},
		{"Chinese to Latin", "chinese", "latin", 0.50, 0.70},
		{"Cyrillic to Latin", "cyrillic", "latin", 0.70, 0.80},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			confidence := calculateConfidence("test", "test", tt.inputScript, tt.outputScript)
			if confidence < tt.expectedMin || confidence > tt.expectedMax {
				t.Errorf("calculateConfidence(%q, %q) = %f, want between %f and %f",
					tt.inputScript, tt.outputScript, confidence, tt.expectedMin, tt.expectedMax)
			}
		})
	}
}