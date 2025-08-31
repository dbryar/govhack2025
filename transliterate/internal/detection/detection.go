// Package detection provides script and language detection capabilities for text analysis.
package detection

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// ScriptInfo contains information about the detected script
type ScriptInfo struct {
	Script     string  // Primary script (e.g., "latin", "cyrillic", "chinese")
	Confidence float64 // Confidence score (0.0-1.0)
	Details    map[string]int // Character counts per script
}

// LanguageHint provides hints about the likely language
type LanguageHint struct {
	Language   string  // Language code (e.g., "vi", "zh", "ru")
	Confidence float64 // Confidence score (0.0-1.0)
	Indicators []string // What led to this detection
}

// DetectScript identifies the primary script used in the text
func DetectScript(text string) ScriptInfo {
	if text == "" {
		return ScriptInfo{Script: "unknown", Confidence: 0.0}
	}

	// Count characters by script
	scriptCounts := make(map[string]int)
	totalLetters := 0

	for _, r := range text {
		if unicode.IsLetter(r) {
			totalLetters++
			script := classifyRune(r)
			scriptCounts[script]++
		}
	}

	if totalLetters == 0 {
		return ScriptInfo{Script: "unknown", Confidence: 0.0, Details: scriptCounts}
	}

	// Find the dominant script, prioritizing specific language variants
	maxScript := "unknown"
	maxCount := 0
	
	// Check for Vietnamese first (it's more specific than general latin)
	if scriptCounts["vietnamese"] > 0 {
		maxScript = "vietnamese"
		maxCount = scriptCounts["vietnamese"]
	} else if scriptCounts["german"] > 0 {
		maxScript = "german"
		maxCount = scriptCounts["german"]
	} else {
		// Fall back to highest count
		for script, count := range scriptCounts {
			if count > maxCount {
				maxScript = script
				maxCount = count
			}
		}
	}

	// Calculate confidence
	confidence := float64(maxCount) / float64(totalLetters)
	
	// Boost confidence for clear script dominance
	if confidence > 0.8 {
		confidence = 0.95
	} else if confidence > 0.6 {
		confidence = 0.85
	} else if confidence > 0.4 {
		confidence = 0.70
	} else if maxScript == "latin" && scriptCounts["latin"] > 0 {
		// Default to Latin with lower confidence if any Latin chars found
		confidence = 0.60
	}

	return ScriptInfo{
		Script:     maxScript,
		Confidence: confidence,
		Details:    scriptCounts,
	}
}

// DetectLanguage attempts to identify the language based on text characteristics
func DetectLanguage(text string, scriptInfo ScriptInfo) LanguageHint {
	indicators := make([]string, 0)
	lowerText := strings.ToLower(text)
	
	switch scriptInfo.Script {
	case "vietnamese", "latin":
		if isVietnamese(lowerText) {
			indicators = append(indicators, "vietnamese_diacritics")
			return LanguageHint{Language: "vi", Confidence: 0.85, Indicators: indicators}
		}
		if isGerman(lowerText) {
			indicators = append(indicators, "german_umlauts")
			return LanguageHint{Language: "de", Confidence: 0.80, Indicators: indicators}
		}
		if isSpanish(lowerText) {
			indicators = append(indicators, "spanish_characters")
			return LanguageHint{Language: "es", Confidence: 0.75, Indicators: indicators}
		}
		
	case "chinese":
		if isTraditionalChinese(text) {
			indicators = append(indicators, "traditional_characters")
			return LanguageHint{Language: "zh-TW", Confidence: 0.80, Indicators: indicators}
		}
		indicators = append(indicators, "simplified_characters")
		return LanguageHint{Language: "zh-CN", Confidence: 0.75, Indicators: indicators}
		
	case "japanese":
		indicators = append(indicators, "hiragana_katakana")
		return LanguageHint{Language: "ja", Confidence: 0.90, Indicators: indicators}
		
	case "cyrillic":
		if isRussian(lowerText) {
			indicators = append(indicators, "russian_patterns")
			return LanguageHint{Language: "ru", Confidence: 0.80, Indicators: indicators}
		}
		indicators = append(indicators, "cyrillic_script")
		return LanguageHint{Language: "ru", Confidence: 0.60, Indicators: indicators}
		
	case "arabic":
		indicators = append(indicators, "arabic_script")
		return LanguageHint{Language: "ar", Confidence: 0.75, Indicators: indicators}
		
	case "greek":
		indicators = append(indicators, "greek_script")
		return LanguageHint{Language: "el", Confidence: 0.90, Indicators: indicators}
	}

	return LanguageHint{Language: "unknown", Confidence: 0.1, Indicators: indicators}
}

// classifyRune determines which script family a rune belongs to
func classifyRune(r rune) string {
	switch {
	// Cyrillic
	case r >= 0x0400 && r <= 0x04FF:
		return "cyrillic"
	case r >= 0x0500 && r <= 0x052F: // Cyrillic Supplement
		return "cyrillic"
	case r >= 0x2DE0 && r <= 0x2DFF: // Cyrillic Extended-A
		return "cyrillic"

	// CJK
	case r >= 0x4E00 && r <= 0x9FFF: // CJK Unified Ideographs
		return "chinese"
	case r >= 0x3400 && r <= 0x4DBF: // CJK Extension A
		return "chinese"
	case r >= 0x20000 && r <= 0x2A6DF: // CJK Extension B
		return "chinese"

	// Japanese
	case r >= 0x3040 && r <= 0x309F: // Hiragana
		return "japanese"
	case r >= 0x30A0 && r <= 0x30FF: // Katakana
		return "japanese"

	// Arabic
	case r >= 0x0600 && r <= 0x06FF: // Arabic
		return "arabic"
	case r >= 0x0750 && r <= 0x077F: // Arabic Supplement
		return "arabic"
	case r >= 0x08A0 && r <= 0x08FF: // Arabic Extended-A
		return "arabic"

	// Greek
	case r >= 0x0370 && r <= 0x03FF: // Greek and Coptic
		return "greek"
	case r >= 0x1F00 && r <= 0x1FFF: // Greek Extended
		return "greek"

	// Latin (including Vietnamese, German, etc.)
	case r >= 0x0041 && r <= 0x005A: // Basic Latin uppercase
		return "latin"
	case r >= 0x0061 && r <= 0x007A: // Basic Latin lowercase
		return "latin"
	case r >= 0x00C0 && r <= 0x024F: // Latin-1, Latin Extended-A/B
		return detectLatinVariant(r)
	case r >= 0x1E00 && r <= 0x1EFF: // Latin Extended Additional
		return detectLatinVariant(r)

	// Hebrew
	case r >= 0x0590 && r <= 0x05FF:
		return "hebrew"

	// Thai
	case r >= 0x0E00 && r <= 0x0E7F:
		return "thai"

	// Korean
	case r >= 0xAC00 && r <= 0xD7AF: // Hangul Syllables
		return "korean"
	case r >= 0x1100 && r <= 0x11FF: // Hangul Jamo
		return "korean"

	default:
		return "unknown"
	}
}

// detectLatinVariant determines if extended Latin characters indicate specific languages
func detectLatinVariant(r rune) string {
	switch {
	// Vietnamese diacritics
	case r == 'ă' || r == 'Ă' || r == 'đ' || r == 'Đ' ||
		 r == 'ư' || r == 'Ư' || r == 'ơ' || r == 'Ơ' ||
		 (r >= 0x1EA0 && r <= 0x1EF9): // Vietnamese combining marks
		return "vietnamese"
	// German umlauts and ß
	case r == 'ä' || r == 'Ä' || r == 'ö' || r == 'Ö' || 
		 r == 'ü' || r == 'Ü' || r == 'ß':
		return "german"
	default:
		return "latin"
	}
}

// isVietnamese checks for Vietnamese-specific patterns
func isVietnamese(text string) bool {
	vietnameseMarkers := []string{
		"ă", "â", "đ", "ê", "ô", "ơ", "ư", "ạ", "ả", "ã", "á", "à",
		"ẹ", "ẻ", "ẽ", "é", "è", "ị", "ỉ", "ĩ", "í", "ì",
		"ọ", "ỏ", "õ", "ó", "ò", "ụ", "ủ", "ũ", "ú", "ù",
		"ỳ", "ỷ", "ỹ", "ý", "ỳ",
	}

	count := 0
	for _, marker := range vietnameseMarkers {
		if strings.Contains(text, marker) {
			count++
		}
	}

	return count >= 2 // Need at least 2 Vietnamese diacritics
}

// isGerman checks for German-specific patterns
func isGerman(text string) bool {
	germanMarkers := []string{"ä", "ö", "ü", "ß", "ae", "oe", "ue"}
	count := 0
	for _, marker := range germanMarkers {
		if strings.Contains(text, marker) {
			count++
		}
	}
	return count >= 1
}

// isSpanish checks for Spanish-specific patterns
func isSpanish(text string) bool {
	spanishMarkers := []string{"ñ", "á", "é", "í", "ó", "ú", "ü"}
	count := 0
	for _, marker := range spanishMarkers {
		if strings.Contains(text, marker) {
			count++
		}
	}
	return count >= 1
}

// isTraditionalChinese checks for Traditional Chinese characters
func isTraditionalChinese(text string) bool {
	// Some characters that are different in Traditional vs Simplified
	traditionalMarkers := []rune{
		'龍', '鳳', '學', '國', '長', '開', '關', '門', '間', '問',
		'風', '飛', '馬', '鳥', '魚', '車', '門', '電', '話', '語',
	}

	for _, r := range text {
		for _, marker := range traditionalMarkers {
			if r == marker {
				return true
			}
		}
	}
	return false
}

// isRussian checks for Russian-specific patterns
func isRussian(text string) bool {
	// Russian has specific letter frequencies and patterns
	russianMarkers := []string{"ъ", "ы", "ь", "э", "ю", "я", "ё"}
	count := 0
	for _, marker := range russianMarkers {
		if strings.Contains(text, marker) {
			count++
		}
	}
	return count >= 1
}

// IsValidUTF8 checks if the text is valid UTF-8
func IsValidUTF8(text string) bool {
	return utf8.ValidString(text)
}

// ContainsScript checks if text contains characters from a specific script
func ContainsScript(text, script string) bool {
	for _, r := range text {
		if unicode.IsLetter(r) && classifyRune(r) == script {
			return true
		}
	}
	return false
}

// Note: Removed whatlanggo dependency due to compilation issues