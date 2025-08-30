// Service transliterate converts names and text between different scripts and locales.
// It tracks transliterations for confidence scoring and learning from user feedback.
package transliterate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"encore.dev/storage/sqldb"
)

// TransliterationRequest represents a request to transliterate text
type TransliterationRequest struct {
	Text         string  `json:"text"`                    // Text to transliterate
	InputScript  string  `json:"input_script,omitempty"`  // e.g., 'cyrillic', 'chinese', 'arabic' (optional - can auto-detect)
	OutputScript string  `json:"output_script"`           // e.g., 'latin', 'ascii'
	InputLocale  *string `json:"input_locale,omitempty"`  // e.g., 'zh-CN', 'ru-RU' (optional)
}

// NameStructure represents parsed name components
type NameStructure struct {
	Family    string   `json:"family"`              // Family/surname in UPPERCASE
	First     string   `json:"first"`               // Given/first name in Title Case
	Middle    []string `json:"middle,omitempty"`    // Middle names/patronymics
	Titles    []string `json:"titles,omitempty"`    // Extracted titles (Dr, Prof, etc)
	FullASCII string   `json:"full_ascii"`          // Complete formatted ASCII name
}

// GenderInference represents inferred gender with confidence
type GenderInference struct {
	Value      string  `json:"value"`      // M, F, or X (unknown/non-binary)
	Confidence float64 `json:"confidence"` // 0.0 to 1.0
	Source     string  `json:"source"`     // "cultural_marker", "statistical", "unknown"
}

// TransliterationResponse represents the result of transliteration
type TransliterationResponse struct {
	ID               string           `json:"id"`
	InputText        string           `json:"input_text"`
	OutputText       string           `json:"output_text"`
	InputScript      string           `json:"input_script"`
	OutputScript     string           `json:"output_script"`
	InputLocale      *string          `json:"input_locale,omitempty"`
	ConfidenceScore  *float64         `json:"confidence_score"`
	AlternativeForms []string         `json:"alternative_forms,omitempty"`
	Name             *NameStructure   `json:"name,omitempty"`           // Structured name parsing
	Gender           *GenderInference `json:"gender,omitempty"`         // Gender inference
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
//encore:api public method=POST path=/api/transliterate
func Transliterate(ctx context.Context, req *TransliterationRequest) (*TransliterationResponse, error) {
	// Validate input
	if err := validateTransliterationRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Detect input script if not provided
	inputScript := req.InputScript
	if inputScript == "" {
		inputScript = detectScript(req.Text)
		if inputScript == "unknown" {
			return nil, errors.New("unable to detect input script")
		}
	}

	// Validate script combination
	if !isSupportedScriptPair(inputScript, req.OutputScript) {
		return nil, fmt.Errorf("unsupported script conversion: %s to %s", inputScript, req.OutputScript)
	}

	// Check if we have this transliteration cached
	cached, err := getCachedTransliteration(ctx, req.Text, inputScript, req.OutputScript, req.InputLocale)
	if err == nil && cached != nil {
		// Parse name structure and gender for cached results (they may not be stored)
		if cached.Name == nil {
			cached.Name = parseName(req.Text, cached.OutputText, inputScript)
		}
		if cached.Gender == nil {
			cached.Gender = inferGender(req.Text, cached.OutputText, inputScript)
		}

		// Update usage count
		_, updateErr := db.Exec(ctx, `
			UPDATE transliterations
			SET usage_count = usage_count + 1, updated_at = NOW()
			WHERE id = $1
		`, cached.ID)
		if updateErr != nil {
			// Log but don't fail - return cached result anyway
		}
		return cached, nil
	}

	// Perform transliteration
	outputText, err := performTransliterationWithValidation(req.Text, inputScript, req.OutputScript, req.InputLocale)
	if err != nil {
		return nil, fmt.Errorf("transliteration failed: %w", err)
	}

	// Calculate confidence score
	confidenceScore := calculateConfidence(req.Text, outputText, inputScript, req.OutputScript)

	// Parse name structure from transliterated text
	nameStructure := parseName(req.Text, outputText, inputScript)

	// Infer gender from name and cultural markers
	genderInference := inferGender(req.Text, outputText, inputScript)

	// Store the result
	result, err := storeTransliteration(ctx, req.Text, outputText, inputScript, req.OutputScript, req.InputLocale, confidenceScore)
	if err != nil {
		return nil, fmt.Errorf("failed to store transliteration: %w", err)
	}

	// Add structured name parsing and gender inference to response
	result.Name = nameStructure
	result.Gender = genderInference

	return result, nil
}

// GetTransliteration retrieves a previously stored transliteration by ID
//
//encore:api public method=GET path=/api/transliterate/:id
func GetTransliteration(ctx context.Context, id string) (*TransliterationResponse, error) {
	// Validate UUID format
	if !isValidUUID(id) {
		return nil, errors.New("invalid transliteration ID format")
	}

	var result TransliterationResponse
	var inputLocale *string

	err := db.QueryRow(ctx, `
		SELECT id, input_text, output_text, input_script, output_script, input_locale, confidence_score
		FROM transliterations
		WHERE id = $1
	`, id).Scan(&result.ID, &result.InputText, &result.OutputText, &result.InputScript,
		&result.OutputScript, &inputLocale, &result.ConfidenceScore)

	if err == sql.ErrNoRows {
		return nil, errors.New("transliteration not found")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	result.InputLocale = inputLocale

	// Add name parsing and gender inference for retrieved records
	result.Name = parseName(result.InputText, result.OutputText, result.InputScript)
	result.Gender = inferGender(result.InputText, result.OutputText, result.InputScript)

	return &result, nil
}

// SubmitFeedback allows users to provide feedback on transliteration results
//
//encore:api public method=POST path=/api/transliterate/:id/feedback
func SubmitFeedback(ctx context.Context, id string, req *FeedbackRequest) error {
	// Validate feedback request
	if err := validateFeedbackRequest(req); err != nil {
		return fmt.Errorf("invalid feedback: %w", err)
	}

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

	if err != nil {
		return fmt.Errorf("failed to store feedback: %w", err)
	}

	return nil
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
	if text == "" {
		return "unknown"
	}

	// Count characters by script
	scriptCounts := make(map[string]int)
	totalChars := 0

	for _, r := range text {
		if unicode.IsLetter(r) {
			totalChars++
			switch {
			case r >= 0x0400 && r <= 0x04FF: // Cyrillic
				scriptCounts["cyrillic"]++
			case r >= 0x4E00 && r <= 0x9FFF: // CJK Ideographs
				scriptCounts["chinese"]++
			case r >= 0x0600 && r <= 0x06FF: // Arabic
				scriptCounts["arabic"]++
			case r >= 0x0370 && r <= 0x03FF: // Greek
				scriptCounts["greek"]++
			case (r >= 0x0041 && r <= 0x005A) || (r >= 0x0061 && r <= 0x007A): // Basic Latin
				scriptCounts["latin"]++
			case r >= 0x0080 && r <= 0x024F: // Extended Latin
				scriptCounts["latin"]++
			default:
				scriptCounts["unknown"]++
			}
		}
	}

	if totalChars == 0 {
		return "unknown"
	}

	// Find the most common script (needs >50% confidence)
	maxScript := "unknown"
	maxCount := 0
	for script, count := range scriptCounts {
		if count > maxCount && float64(count)/float64(totalChars) > 0.5 {
			maxScript = script
			maxCount = count
		}
	}

	if maxScript == "unknown" && scriptCounts["latin"] > 0 {
		return "latin" // Default to latin if some latin chars found
	}

	return maxScript
}

// performTransliteration converts text using character mappings and fallback rules
func performTransliteration(text, inputScript, outputScript string, inputLocale *string) string {
	result := strings.Builder{}

	// Process text character by character
	for _, r := range text {
		sourceChar := string(r)

		// First try database lookup
		mapped := getCharacterMapping(sourceChar, inputScript, outputScript, inputLocale)
		if mapped != "" {
			result.WriteString(mapped)
			continue
		}

		// Apply built-in transliteration rules
		builtinMapped := applyBuiltinRules(r, inputScript, outputScript)
		if builtinMapped != "" {
			result.WriteString(builtinMapped)
			continue
		}

		// Fallback: ASCII approximation or keep original
		if outputScript == "ascii" {
			asciiApprox := approximateToASCII(r)
			result.WriteString(asciiApprox)
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// getCharacterMapping retrieves character mapping from database
func getCharacterMapping(sourceChar, sourceScript, targetScript string, locale *string) string {
	var targetChar string

	// Query character mappings table for exact match
	err := db.QueryRow(context.Background(), `
		SELECT target_char
		FROM character_mappings
		WHERE source_char = $1
			AND source_script = $2
			AND target_script = $3
			AND ($4::text IS NULL OR locale = $4 OR locale IS NULL)
		ORDER BY
			CASE WHEN locale = $4 THEN 1 ELSE 2 END,
			frequency_weight DESC
		LIMIT 1
	`, sourceChar, sourceScript, targetScript, locale).Scan(&targetChar)

	if err == sql.ErrNoRows {
		return "" // No mapping found
	}
	if err != nil {
		// Log error but don't fail transliteration
		return ""
	}

	return targetChar
}

// calculateConfidence computes a confidence score based on multiple factors
func calculateConfidence(inputText, outputText, inputScript, outputScript string) float64 {
	baseConfidence := 0.50 // Conservative baseline

	// Script compatibility scoring
	scriptBonus := calculateScriptCompatibility(inputScript, outputScript)
	baseConfidence += scriptBonus

	// Character coverage scoring
	coverageBonus := calculateCharacterCoverage(inputText, outputText)
	baseConfidence += coverageBonus

	// Length preservation bonus
	lengthRatio := float64(len(outputText)) / float64(len(inputText))
	if lengthRatio >= 0.5 && lengthRatio <= 2.0 {
		baseConfidence += 0.1 // Reasonable length preservation
	}

	// Ensure confidence stays within bounds
	if baseConfidence > 1.0 {
		baseConfidence = 1.0
	}
	if baseConfidence < 0.1 {
		baseConfidence = 0.1
	}

	return baseConfidence
}

// applyBuiltinRules applies hardcoded transliteration rules for common cases
func applyBuiltinRules(r rune, inputScript, outputScript string) string {
	// Cyrillic to Latin mappings
	if inputScript == "cyrillic" && outputScript == "latin" {
		return transliterateCyrillicToLatin(r)
	}

	// Chinese character approximations
	if inputScript == "chinese" && (outputScript == "latin" || outputScript == "ascii") {
		return transliterateChineseToLatin(r)
	}

	// Arabic to Latin mappings
	if inputScript == "arabic" && outputScript == "latin" {
		return transliterateArabicToLatin(r)
	}

	// Greek to Latin mappings
	if inputScript == "greek" && outputScript == "latin" {
		return transliterateGreekToLatin(r)
	}

	return ""
}

// transliterateCyrillicToLatin provides standard Cyrillic transliteration
func transliterateCyrillicToLatin(r rune) string {
	cyrillicMap := map[rune]string{
		// Uppercase
		'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D", 'Е': "E", 'Ё': "Yo",
		'Ж': "Zh", 'З': "Z", 'И': "I", 'Й': "Y", 'К': "K", 'Л': "L", 'М': "M",
		'Н': "N", 'О': "O", 'П': "P", 'Р': "R", 'С': "S", 'Т': "T", 'У': "U",
		'Ф': "F", 'Х': "Kh", 'Ц': "Ts", 'Ч': "Ch", 'Ш': "Sh", 'Щ': "Shch",
		'Ъ': "", 'Ы': "Y", 'Ь': "", 'Э': "E", 'Ю': "Yu", 'Я': "Ya",
		// Lowercase
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo",
		'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
		'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
		'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "shch",
		'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
	}

	if mapped, exists := cyrillicMap[r]; exists {
		return mapped
	}
	return ""
}

// transliterateChineseToLatin provides basic Chinese character mappings
func transliterateChineseToLatin(r rune) string {
	// Basic common Chinese characters - in practice would use proper Pinyin lookup
	chineseMap := map[rune]string{
		'你': "ni", '好': "hao", '是': "shi", '的': "de", '我': "wo",
		'他': "ta", '她': "ta", '们': "men", '有': "you", '在': "zai",
		'了': "le", '不': "bu", '就': "jiu", '人': "ren", '都': "dou",
		'一': "yi", '二': "er", '三': "san", '四': "si", '五': "wu",
		'六': "liu", '七': "qi", '八': "ba", '九': "jiu", '十': "shi",
	}

	if mapped, exists := chineseMap[r]; exists {
		return mapped
	}
	return ""
}

// transliterateArabicToLatin provides Arabic transliteration
func transliterateArabicToLatin(r rune) string {
	arabicMap := map[rune]string{
		'ا': "a", 'ب': "b", 'ت': "t", 'ث': "th", 'ج': "j", 'ح': "h",
		'خ': "kh", 'د': "d", 'ذ': "dh", 'ر': "r", 'ز': "z", 'س': "s",
		'ش': "sh", 'ص': "s", 'ض': "d", 'ط': "t", 'ظ': "z", 'ع': "'",
		'غ': "gh", 'ف': "f", 'ق': "q", 'ك': "k", 'ل': "l", 'م': "m",
		'ن': "n", 'ه': "h", 'و': "w", 'ي': "y",
	}

	if mapped, exists := arabicMap[r]; exists {
		return mapped
	}
	return ""
}

// transliterateGreekToLatin provides Greek transliteration
func transliterateGreekToLatin(r rune) string {
	greekMap := map[rune]string{
		'Α': "A", 'Β': "B", 'Γ': "G", 'Δ': "D", 'Ε': "E", 'Ζ': "Z",
		'Η': "H", 'Θ': "Th", 'Ι': "I", 'Κ': "K", 'Λ': "L", 'Μ': "M",
		'Ν': "N", 'Ξ': "X", 'Ο': "O", 'Π': "P", 'Ρ': "R", 'Σ': "S",
		'Τ': "T", 'Υ': "Y", 'Φ': "Ph", 'Χ': "Ch", 'Ψ': "Ps", 'Ω': "O",
		'α': "a", 'β': "b", 'γ': "g", 'δ': "d", 'ε': "e", 'ζ': "z",
		'η': "h", 'θ': "th", 'ι': "i", 'κ': "k", 'λ': "l", 'μ': "m",
		'ν': "n", 'ξ': "x", 'ο': "o", 'π': "p", 'ρ': "r", 'σ': "s", 'ς': "s",
		'τ': "t", 'υ': "y", 'φ': "ph", 'χ': "ch", 'ψ': "ps", 'ω': "o",
	}

	if mapped, exists := greekMap[r]; exists {
		return mapped
	}
	return ""
}

// approximateToASCII converts Unicode characters to closest ASCII equivalents
func approximateToASCII(r rune) string {
	// Handle accented characters
	asciiMap := map[rune]string{
		// Accented vowels
		'á': "a", 'à': "a", 'â': "a", 'ã': "a", 'ä': "a", 'å': "a", 'ā': "a",
		'é': "e", 'è': "e", 'ê': "e", 'ë': "e", 'ē': "e",
		'í': "i", 'ì': "i", 'î': "i", 'ï': "i", 'ī': "i",
		'ó': "o", 'ò': "o", 'ô': "o", 'õ': "o", 'ö': "o", 'ø': "o", 'ō': "o",
		'ú': "u", 'ù': "u", 'û': "u", 'ü': "u", 'ū': "u",
		// Uppercase versions
		'Á': "A", 'À': "A", 'Â': "A", 'Ã': "A", 'Ä': "A", 'Å': "A", 'Ā': "A",
		'É': "E", 'È': "E", 'Ê': "E", 'Ë': "E", 'Ē': "E",
		'Í': "I", 'Ì': "I", 'Î': "I", 'Ï': "I", 'Ī': "I",
		'Ó': "O", 'Ò': "O", 'Ô': "O", 'Õ': "O", 'Ö': "O", 'Ø': "O", 'Ō': "O",
		'Ú': "U", 'Ù': "U", 'Û': "U", 'Ü': "U", 'Ū': "U",
		// Other common characters
		'ç': "c", 'Ç': "C", 'ñ': "n", 'Ñ': "N",
		'ß': "ss", 'æ': "ae", 'Æ': "AE", 'œ': "oe", 'Œ': "OE",
	}

	if mapped, exists := asciiMap[r]; exists {
		return mapped
	}

	// If it's already ASCII, return as-is
	if r < 128 {
		return string(r)
	}

	// For other characters, try to approximate based on Unicode category
	if unicode.IsLetter(r) {
		return "?" // Unknown letter
	}
	if unicode.IsDigit(r) {
		return "0" // Unknown digit
	}
	if unicode.IsSpace(r) {
		return " "
	}
	if unicode.IsPunct(r) {
		return "."
	}

	return "" // Skip other characters
}

// calculateScriptCompatibility returns a bonus based on script pairing difficulty
func calculateScriptCompatibility(inputScript, outputScript string) float64 {
	// High compatibility pairs
	highCompatibility := map[string]map[string]bool{
		"latin": {"ascii": true},
		"ascii": {"latin": true},
	}

	// Medium compatibility pairs
	mediumCompatibility := map[string]map[string]bool{
		"cyrillic": {"latin": true, "ascii": true},
		"greek":    {"latin": true, "ascii": true},
	}

	// Low compatibility pairs (complex scripts)
	lowCompatibility := map[string]map[string]bool{
		"chinese": {"latin": true, "ascii": true},
		"arabic":  {"latin": true, "ascii": true},
	}

	if highCompatibility[inputScript] != nil && highCompatibility[inputScript][outputScript] {
		return 0.3
	}
	if mediumCompatibility[inputScript] != nil && mediumCompatibility[inputScript][outputScript] {
		return 0.2
	}
	if lowCompatibility[inputScript] != nil && lowCompatibility[inputScript][outputScript] {
		return 0.1
	}

	return 0.0 // Unknown or unsupported pairing
}

// calculateCharacterCoverage estimates how well the output covers the input
func calculateCharacterCoverage(inputText, outputText string) float64 {
	// Count non-whitespace characters
	inputChars := countNonWhitespaceChars(inputText)
	outputChars := countNonWhitespaceChars(outputText)

	if inputChars == 0 {
		return 0.0
	}

	// Penalise outputs that are too short (lost information)
	if outputChars == 0 {
		return -0.2
	}

	// Bonus for reasonable coverage
	coverageRatio := float64(outputChars) / float64(inputChars)
	if coverageRatio >= 0.5 && coverageRatio <= 1.5 {
		return 0.1
	}

	return 0.0
}

// countNonWhitespaceChars counts non-whitespace characters in a string
func countNonWhitespaceChars(text string) int {
	count := 0
	for _, r := range text {
		if !unicode.IsSpace(r) {
			count++
		}
	}
	return count
}

// Validation functions

// validateTransliterationRequest validates the input request
func validateTransliterationRequest(req *TransliterationRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if strings.TrimSpace(req.Text) == "" {
		return errors.New("text cannot be empty")
	}

	if len(req.Text) > 10000 { // Reasonable limit
		return errors.New("text too long (maximum 10,000 characters)")
	}

	if !utf8.ValidString(req.Text) {
		return errors.New("text contains invalid UTF-8 sequences")
	}

	if req.OutputScript == "" {
		return errors.New("output_script is required")
	}

	// Validate script names
	validScripts := map[string]bool{
		"latin": true, "ascii": true, "cyrillic": true,
		"chinese": true, "arabic": true, "greek": true,
		"vietnamese": true, "indonesian": true, "malayalam": true,
	}

	if req.InputScript != "" && !validScripts[req.InputScript] {
		return fmt.Errorf("unsupported input script: %s", req.InputScript)
	}

	if !validScripts[req.OutputScript] {
		return fmt.Errorf("unsupported output script: %s", req.OutputScript)
	}

	// Validate locale format if provided
	if req.InputLocale != nil && !isValidLocale(*req.InputLocale) {
		return fmt.Errorf("invalid locale format: %s", *req.InputLocale)
	}

	return nil
}

// validateFeedbackRequest validates feedback input
func validateFeedbackRequest(req *FeedbackRequest) error {
	if req == nil {
		return errors.New("feedback request cannot be nil")
	}

	if strings.TrimSpace(req.SuggestedOutput) == "" {
		return errors.New("suggested_output cannot be empty")
	}

	if len(req.SuggestedOutput) > 10000 {
		return errors.New("suggested_output too long")
	}

	validFeedbackTypes := map[string]bool{
		"correction": true, "alternative": true, "preferred": true,
	}

	if !validFeedbackTypes[req.FeedbackType] {
		return fmt.Errorf("invalid feedback_type: %s (must be 'correction', 'alternative', or 'preferred')", req.FeedbackType)
	}

	return nil
}

// isSupportedScriptPair checks if the script conversion is supported
func isSupportedScriptPair(inputScript, outputScript string) bool {
	supportedPairs := map[string]map[string]bool{
		"latin":      {"ascii": true, "latin": true},
		"ascii":      {"latin": true, "ascii": true},
		"cyrillic":   {"latin": true, "ascii": true},
		"chinese":    {"latin": true, "ascii": true},
		"arabic":     {"latin": true, "ascii": true},
		"greek":      {"latin": true, "ascii": true},
		"vietnamese": {"latin": true, "ascii": true},
		"indonesian": {"latin": true, "ascii": true},
		"malayalam":  {"latin": true, "ascii": true},
	}

	if targets, exists := supportedPairs[inputScript]; exists {
		return targets[outputScript]
	}

	return false
}

// isValidUUID checks if a string is a valid UUID format
func isValidUUID(uuid string) bool {
	// Basic UUID format validation (36 characters with hyphens in right places)
	if len(uuid) != 36 {
		return false
	}

	// Check hyphen positions
	if uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		return false
	}

	// Check all other characters are hex digits
	for i, r := range uuid {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			continue // Skip hyphens
		}
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return false
		}
	}

	return true
}

// isValidLocale checks if a locale string follows ISO format
func isValidLocale(locale string) bool {
	if locale == "" {
		return false
	}

	// Basic check for xx-XX format (language-country)
	parts := strings.Split(locale, "-")
	if len(parts) < 1 || len(parts) > 2 {
		return false
	}

	// Language code should be 2-3 lowercase letters
	lang := parts[0]
	if len(lang) < 2 || len(lang) > 3 {
		return false
	}

	for _, r := range lang {
		if r < 'a' || r > 'z' {
			return false
		}
	}

	// Country code should be 2 uppercase letters (if present)
	if len(parts) == 2 {
		country := parts[1]
		if len(country) != 2 {
			return false
		}
		for _, r := range country {
			if r < 'A' || r > 'Z' {
				return false
			}
		}
	}

	return true
}

// performTransliterationWithValidation wraps performTransliteration with error handling
func performTransliterationWithValidation(text, inputScript, outputScript string, inputLocale *string) (string, error) {
	if text == "" {
		return "", errors.New("empty input text")
	}

	result := performTransliteration(text, inputScript, outputScript, inputLocale)

	// Validate output
	if result == "" {
		return "", errors.New("transliteration produced empty result")
	}

	if !utf8.ValidString(result) {
		return "", errors.New("transliteration produced invalid UTF-8")
	}

	return result, nil
}

// parseName extracts structured name components from transliterated text
func parseName(originalText, transliteratedText, inputScript string) *NameStructure {
	if transliteratedText == "" {
		return nil
	}

	// Extract titles first
	titles := extractTitles(transliteratedText)
	textWithoutTitles := removeTitles(transliteratedText, titles)

	// Handle cultural naming conventions
	switch inputScript {
	case "chinese":
		return parseChineseName(textWithoutTitles, titles)
	case "vietnamese":
		return parseVietnameseName(originalText, textWithoutTitles, titles)
	case "arabic":
		return parseArabicName(textWithoutTitles, titles)
	case "indonesian", "malayalam":
		return parseMononymOrPatronymic(textWithoutTitles, titles, inputScript)
	default:
		return parseWesternName(textWithoutTitles, titles)
	}
}

// extractTitles identifies and extracts titles from text
func extractTitles(text string) []string {
	titlePatterns := []string{
		"DR", "DOCTOR", "PROF", "PROFESSOR", "MR", "MRS", "MS", "MISS",
		"SIR", "DAME", "LORD", "LADY", "HON", "HONOURABLE", "REV", "REVEREND",
	}

	var titles []string
	words := strings.Fields(strings.ToUpper(text))

	for _, word := range words {
		cleanWord := strings.Trim(word, ".,")
		for _, pattern := range titlePatterns {
			if cleanWord == pattern {
				titles = append(titles, formatTitle(cleanWord))
				break
			}
		}
	}

	return titles
}

// removeTitles removes identified titles from text
func removeTitles(text string, titles []string) string {
	if len(titles) == 0 {
		return text
	}

	result := text
	for _, title := range titles {
		patterns := []string{
			strings.ToUpper(title),
			strings.ToLower(title),
			strings.Title(strings.ToLower(title)),
		}

		for _, pattern := range patterns {
			result = strings.ReplaceAll(result, pattern+".", "")
			result = strings.ReplaceAll(result, pattern+" ", "")
		}
	}

	return strings.TrimSpace(result)
}

// formatTitle formats title for display
func formatTitle(title string) string {
	switch strings.ToUpper(title) {
	case "DR", "DOCTOR":
		return "DR"
	case "PROF", "PROFESSOR":
		return "PROF"
	case "REV", "REVEREND":
		return "REV"
	case "HON", "HONOURABLE":
		return "HON"
	default:
		return strings.ToUpper(title)
	}
}

// parseVietnameseName handles Vietnamese naming conventions
func parseVietnameseName(original, transliterated string, titles []string) *NameStructure {
	parts := strings.Fields(transliterated)
	if len(parts) == 0 {
		return &NameStructure{FullASCII: transliterated, Titles: titles}
	}

	// Gender markers are handled in gender inference, not needed here for name structure
	_ = strings.Contains(strings.ToLower(original), "văn") || strings.Contains(strings.ToLower(transliterated), "van")
	_ = strings.Contains(strings.ToLower(original), "thị") || strings.Contains(strings.ToLower(transliterated), "thi")

	var family, first string
	var middle []string

	if len(parts) >= 2 {
		// Vietnamese: Family name first, then middle names, then given name
		family = strings.ToUpper(parts[0])
		first = strings.Title(strings.ToLower(parts[len(parts)-1]))

		// Middle names (excluding gender markers)
		for i := 1; i < len(parts)-1; i++ {
			part := parts[i]
			partLower := strings.ToLower(part)

			// Skip Vietnamese gender markers
			if partLower != "van" && partLower != "thi" && partLower != "văn" && partLower != "thị" {
				middle = append(middle, strings.Title(strings.ToLower(part)))
			}
		}
	} else {
		// Single name - could be family or given
		first = strings.Title(strings.ToLower(parts[0]))
	}

	// Format full ASCII name
	fullName := formatFullName(family, first, middle, titles)

	return &NameStructure{
		Family:    family,
		First:     first,
		Middle:    middle,
		Titles:    titles,
		FullASCII: fullName,
	}
}

// parseChineseName handles Chinese naming conventions
func parseChineseName(text string, titles []string) *NameStructure {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return &NameStructure{FullASCII: text, Titles: titles}
	}

	var family, first string
	var middle []string

	if len(parts) >= 2 {
		// Chinese: Family name first, then given names
		family = strings.ToUpper(parts[0])
		if len(parts) == 2 {
			first = strings.Title(strings.ToLower(parts[1]))
		} else {
			// Multiple given names
			first = strings.Title(strings.ToLower(parts[len(parts)-1]))
			for i := 1; i < len(parts)-1; i++ {
				middle = append(middle, strings.Title(strings.ToLower(parts[i])))
			}
		}
	} else {
		// Single name
		first = strings.Title(strings.ToLower(parts[0]))
	}

	fullName := formatFullName(family, first, middle, titles)

	return &NameStructure{
		Family:    family,
		First:     first,
		Middle:    middle,
		Titles:    titles,
		FullASCII: fullName,
	}
}

// parseArabicName handles Arabic naming conventions
func parseArabicName(text string, titles []string) *NameStructure {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return &NameStructure{FullASCII: text, Titles: titles}
	}

	var family, first string
	var middle []string

	if len(parts) >= 2 {
		// Arabic: Given name first, patronymics/family name last
		first = strings.Title(strings.ToLower(parts[0]))
		family = strings.ToUpper(parts[len(parts)-1])

		// Middle names/patronymics (ibn, bin, bint, etc.)
		for i := 1; i < len(parts)-1; i++ {
			middle = append(middle, strings.Title(strings.ToLower(parts[i])))
		}
	} else {
		first = strings.Title(strings.ToLower(parts[0]))
	}

	fullName := formatFullName(family, first, middle, titles)

	return &NameStructure{
		Family:    family,
		First:     first,
		Middle:    middle,
		Titles:    titles,
		FullASCII: fullName,
	}
}

// parseMononymOrPatronymic handles single names or patronymic structures
func parseMononymOrPatronymic(text string, titles []string, script string) *NameStructure {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return &NameStructure{FullASCII: text, Titles: titles}
	}

	// Check for patronymic indicators
	hasPatronymic := false
	for _, part := range parts {
		lower := strings.ToLower(part)
		if lower == "bin" || lower == "binti" || lower == "ibn" || lower == "bint" {
			hasPatronymic = true
			break
		}
	}

	var family, first string
	var middle []string

	if len(parts) == 1 {
		// Mononym - single name
		first = strings.Title(strings.ToLower(parts[0]))
	} else if hasPatronymic {
		// Handle patronymic structure
		first = strings.Title(strings.ToLower(parts[0]))
		for i := 1; i < len(parts); i++ {
			middle = append(middle, strings.Title(strings.ToLower(parts[i])))
		}
	} else {
		// Regular multi-part name
		first = strings.Title(strings.ToLower(parts[0]))
		if len(parts) > 1 {
			family = strings.ToUpper(parts[len(parts)-1])
		}
		for i := 1; i < len(parts)-1; i++ {
			middle = append(middle, strings.Title(strings.ToLower(parts[i])))
		}
	}

	fullName := formatFullName(family, first, middle, titles)

	return &NameStructure{
		Family:    family,
		First:     first,
		Middle:    middle,
		Titles:    titles,
		FullASCII: fullName,
	}
}

// parseWesternName handles Western naming conventions
func parseWesternName(text string, titles []string) *NameStructure {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return &NameStructure{FullASCII: text, Titles: titles}
	}

	var family, first string
	var middle []string

	if len(parts) == 1 {
		first = strings.Title(strings.ToLower(parts[0]))
	} else if len(parts) == 2 {
		first = strings.Title(strings.ToLower(parts[0]))
		family = strings.ToUpper(parts[1])
	} else {
		// First + middle names + last
		first = strings.Title(strings.ToLower(parts[0]))
		family = strings.ToUpper(parts[len(parts)-1])
		for i := 1; i < len(parts)-1; i++ {
			middle = append(middle, strings.Title(strings.ToLower(parts[i])))
		}
	}

	fullName := formatFullName(family, first, middle, titles)

	return &NameStructure{
		Family:    family,
		First:     first,
		Middle:    middle,
		Titles:    titles,
		FullASCII: fullName,
	}
}

// formatFullName creates the complete formatted ASCII name
func formatFullName(family, first string, middle []string, titles []string) string {
	var parts []string

	// Add titles
	for _, title := range titles {
		parts = append(parts, title)
	}

	// Add first name
	if first != "" {
		parts = append(parts, first)
	}

	// Add middle names
	for _, m := range middle {
		if m != "" {
			parts = append(parts, m)
		}
	}

	// Add family name
	if family != "" {
		parts = append(parts, family)
	}

	return strings.Join(parts, " ")
}

// inferGender attempts to determine gender from name and cultural markers
func inferGender(originalText, transliteratedText, inputScript string) *GenderInference {
	// Default to unknown
	inference := &GenderInference{
		Value:      "X",
		Confidence: 0.1,
		Source:     "unknown",
	}

	// Vietnamese gender markers
	if inputScript == "vietnamese" || strings.Contains(inputScript, "vietnam") {
		original := strings.ToLower(originalText)
		transliterated := strings.ToLower(transliteratedText)

		if strings.Contains(original, "văn") || strings.Contains(transliterated, "van") {
			return &GenderInference{
				Value:      "M",
				Confidence: 0.85,
				Source:     "cultural_marker",
			}
		}

		if strings.Contains(original, "thị") || strings.Contains(transliterated, "thi") {
			return &GenderInference{
				Value:      "F",
				Confidence: 0.85,
				Source:     "cultural_marker",
			}
		}
	}

	// Arabic patronymic indicators
	if inputScript == "arabic" {
		text := strings.ToLower(transliteratedText)
		if strings.Contains(text, "bin ") || strings.Contains(text, "ibn ") {
			inference.Value = "M"
			inference.Confidence = 0.75
			inference.Source = "cultural_marker"
		} else if strings.Contains(text, "bint ") {
			inference.Value = "F"
			inference.Confidence = 0.75
			inference.Source = "cultural_marker"
		}
	}

	// Malaysian/Indonesian patronymic
	if inputScript == "indonesian" || inputScript == "malayalam" {
		text := strings.ToLower(transliteratedText)
		if strings.Contains(text, "bin ") {
			inference.Value = "M"
			inference.Confidence = 0.80
			inference.Source = "cultural_marker"
		} else if strings.Contains(text, "binti ") {
			inference.Value = "F"
			inference.Confidence = 0.80
			inference.Source = "cultural_marker"
		}
	}

	return inference
}

// ServeApp serves the Hugo-generated frontend by reading files from disk
//
//encore:api public raw method=GET path=/app/*path
func ServeApp(w http.ResponseWriter, req *http.Request) {
	// Extract the path after /app/
	path := req.URL.Path[5:] // Remove "/app/" prefix

	// Handle root app path
	if path == "" || path == "/" {
		path = "index.html"
	}

	// Build the file path to the Hugo dist directory
	filePath := filepath.Join("../frontend/dist", path)

	// Try to serve the requested file
	http.ServeFile(w, req, filePath)
}

// Database connection
var db = sqldb.NewDatabase("transliterate", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})