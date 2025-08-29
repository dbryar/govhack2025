// Service transliterate converts names and text between different scripts and locales.
// It tracks transliterations for confidence scoring and learning from user feedback.
package transliterate

import (
	"context"
	"fmt"
	"strings"
	"time"

	"encore.dev/storage/sqldb"
)

// TransliterationRequest represents a request to transliterate text
type TransliterationRequest struct {
	Text         string  `json:"text"`                    // Text to transliterate
	InputScript  string  `json:"input_script,omitempty"`  // e.g., 'cyrillic', 'chinese', 'arabic' (optional - can auto-detect)
	OutputScript string  `json:"output_script"`           // e.g., 'latin', 'ascii'
	InputLocale  *string `json:"input_locale,omitempty"`  // e.g., 'zh-CN', 'ru-RU' (optional)
}

// TransliterationResponse represents the result of transliteration
type TransliterationResponse struct {
	ID               string   `json:"id"`
	InputText        string   `json:"input_text"`
	OutputText       string   `json:"output_text"`
	InputScript      string   `json:"input_script"`
	OutputScript     string   `json:"output_script"`
	InputLocale      *string  `json:"input_locale,omitempty"`
	ConfidenceScore  *float64 `json:"confidence_score"`
	AlternativeForms []string `json:"alternative_forms,omitempty"`
}

// FeedbackRequest represents user feedback on transliteration results
type FeedbackRequest struct {
	TransliterationID string `json:"transliteration_id"`
	SuggestedOutput   string `json:"suggested_output"`
	FeedbackType      string `json:"feedback_type"` // 'correction', 'alternative', 'preferred'
	UserContext       string `json:"user_context,omitempty"`
}

// Transliterate converts text from one script to another
//
//encore:api public method=POST path=/transliterate
func Transliterate(ctx context.Context, req *TransliterationRequest) (*TransliterationResponse, error) {
	// Detect input script if not provided
	inputScript := req.InputScript
	if inputScript == "" {
		inputScript = detectScript(req.Text)
	}

	// Check if we have this transliteration cached
	cached, err := getCachedTransliteration(ctx, req.Text, inputScript, req.OutputScript, req.InputLocale)
	if err == nil && cached != nil {
		// Update usage count
		_, _ = db.Exec(ctx, `
			UPDATE transliterations 
			SET usage_count = usage_count + 1, updated_at = NOW() 
			WHERE id = $1
		`, cached.ID)
		return cached, nil
	}

	// Perform transliteration
	outputText := performTransliteration(req.Text, inputScript, req.OutputScript, req.InputLocale)
	
	// Calculate confidence score
	confidenceScore := calculateConfidence(req.Text, outputText, inputScript, req.OutputScript)

	// Store the result
	result, err := storeTransliteration(ctx, req.Text, outputText, inputScript, req.OutputScript, req.InputLocale, confidenceScore)
	if err != nil {
		return nil, fmt.Errorf("failed to store transliteration: %w", err)
	}

	return result, nil
}

// GetTransliteration retrieves a previously stored transliteration by ID
//
//encore:api public method=GET path=/transliterate/:id
func GetTransliteration(ctx context.Context, id string) (*TransliterationResponse, error) {
	var result TransliterationResponse
	var inputLocale *string

	err := db.QueryRow(ctx, `
		SELECT id, input_text, output_text, input_script, output_script, input_locale, confidence_score
		FROM transliterations 
		WHERE id = $1
	`, id).Scan(&result.ID, &result.InputText, &result.OutputText, &result.InputScript, 
		&result.OutputScript, &inputLocale, &result.ConfidenceScore)
	
	if err != nil {
		return nil, fmt.Errorf("transliteration not found: %w", err)
	}

	result.InputLocale = inputLocale
	return &result, nil
}

// SubmitFeedback allows users to provide feedback on transliteration results
//
//encore:api public method=POST path=/transliterate/:id/feedback
func SubmitFeedback(ctx context.Context, id string, req *FeedbackRequest) error {
	// Verify the transliteration exists
	_, err := GetTransliteration(ctx, id)
	if err != nil {
		return fmt.Errorf("invalid transliteration ID: %w", err)
	}

	// Store feedback
	_, err = db.Exec(ctx, `
		INSERT INTO transliteration_feedback (transliteration_id, suggested_output, feedback_type, user_context)
		VALUES ($1, $2, $3, $4)
	`, id, req.SuggestedOutput, req.FeedbackType, req.UserContext)

	return err
}

// Helper functions

func getCachedTransliteration(ctx context.Context, inputText, inputScript, outputScript string, inputLocale *string) (*TransliterationResponse, error) {
	var result TransliterationResponse
	var cachedInputLocale *string

	err := db.QueryRow(ctx, `
		SELECT id, input_text, output_text, input_script, output_script, input_locale, confidence_score
		FROM transliterations 
		WHERE input_text = $1 AND input_script = $2 AND output_script = $3 
		AND ($4::text IS NULL OR input_locale = $4)
		ORDER BY usage_count DESC, updated_at DESC
		LIMIT 1
	`, inputText, inputScript, outputScript, inputLocale).Scan(
		&result.ID, &result.InputText, &result.OutputText, 
		&result.InputScript, &result.OutputScript, &cachedInputLocale, &result.ConfidenceScore)

	if err != nil {
		return nil, err
	}

	result.InputLocale = cachedInputLocale
	return &result, nil
}

func storeTransliteration(ctx context.Context, inputText, outputText, inputScript, outputScript string, inputLocale *string, confidenceScore float64) (*TransliterationResponse, error) {
	var id string
	err := db.QueryRow(ctx, `
		INSERT INTO transliterations (input_text, output_text, input_script, output_script, input_locale, confidence_score)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, inputText, outputText, inputScript, outputScript, inputLocale, confidenceScore).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &TransliterationResponse{
		ID:              id,
		InputText:       inputText,
		OutputText:      outputText,
		InputScript:     inputScript,
		OutputScript:    outputScript,
		InputLocale:     inputLocale,
		ConfidenceScore: &confidenceScore,
	}, nil
}

// detectScript attempts to identify the script of the input text
func detectScript(text string) string {
	// Simple script detection - can be enhanced with more sophisticated algorithms
	for _, r := range text {
		switch {
		case r >= 0x0400 && r <= 0x04FF: // Cyrillic
			return "cyrillic"
		case r >= 0x4E00 && r <= 0x9FFF: // CJK Ideographs
			return "chinese"
		case r >= 0x0600 && r <= 0x06FF: // Arabic
			return "arabic"
		case r >= 0x0370 && r <= 0x03FF: // Greek
			return "greek"
		}
	}
	return "latin" // Default assumption
}

// performTransliteration is a placeholder for the actual transliteration logic
func performTransliteration(text, inputScript, outputScript string, inputLocale *string) string {
	// This is a basic placeholder - real implementation would use sophisticated
	// transliteration algorithms, possibly involving machine learning models
	
	// For now, perform simple character-by-character mapping
	result := strings.Builder{}
	
	for _, r := range text {
		// Look up character mapping in database or use default rules
		mapped := getCharacterMapping(string(r), inputScript, outputScript, inputLocale)
		if mapped != "" {
			result.WriteString(mapped)
		} else {
			// Fallback: keep original character if no mapping found
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

// getCharacterMapping retrieves character mapping from database
func getCharacterMapping(sourceChar, sourceScript, targetScript string, locale *string) string {
	// In a real implementation, this would query the character_mappings table
	// For now, return a placeholder
	
	// Simple hardcoded examples for demonstration
	switch sourceChar {
	case "Привет":
		if sourceScript == "cyrillic" && targetScript == "latin" {
			return "Privet"
		}
	case "你好":
		if sourceScript == "chinese" && targetScript == "latin" {
			return "ni hao"
		}
	}
	
	return ""
}

// calculateConfidence computes a confidence score for the transliteration
func calculateConfidence(inputText, outputText, inputScript, outputScript string) float64 {
	// Placeholder confidence calculation
	// Real implementation would consider factors like:
	// - Character mapping frequency
	// - Historical accuracy
	// - User feedback patterns
	// - Script compatibility
	
	baseConfidence := 0.75 // Default confidence
	
	// Adjust based on script pairing
	switch {
	case inputScript == "latin" && outputScript == "ascii":
		baseConfidence = 0.95 // High confidence for similar scripts
	case inputScript == "chinese" && outputScript == "latin":
		baseConfidence = 0.60 // Lower confidence for complex scripts
	}
	
	return baseConfidence
}

// Database connection
var db = sqldb.NewDatabase("transliterate", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})