// Package gender provides culturally-aware gender inference capabilities.
package gender

import (
	"strings"
)

// Inference represents a gender inference with confidence and reasoning
type Inference struct {
	Value      string  `json:"value"`      // M, F, or X (unknown/non-binary)
	Confidence float64 `json:"confidence"` // 0.0 to 1.0
	Source     string  `json:"source"`     // "cultural_marker", "statistical", "unknown"
	Reason     string  `json:"reason"`     // Human-readable explanation
}

// Engine provides gender inference capabilities
type Engine struct {
	useStatistical bool
	culturalOnly   bool
}

// NewEngine creates a new gender inference engine
func NewEngine(useStatistical, culturalOnly bool) *Engine {
	return &Engine{
		useStatistical: useStatistical,
		culturalOnly:   culturalOnly,
	}
}

// InferGender attempts to determine gender from name and cultural context
func (e *Engine) InferGender(originalText, transliteratedText, culture, language string) *Inference {
	// Default to unknown
	result := &Inference{
		Value:      "X",
		Confidence: 0.1,
		Source:     "unknown",
		Reason:     "No gender indicators found",
	}

	// Try cultural markers first (highest confidence)
	if cultural := e.inferFromCulturalMarkers(originalText, transliteratedText, culture, language); cultural.Confidence > result.Confidence {
		result = cultural
	}

	// Try statistical inference if enabled and cultural inference is weak
	if e.useStatistical && result.Confidence < 0.5 {
		if statistical := e.inferFromStatisticalPatterns(transliteratedText, culture, language); statistical.Confidence > result.Confidence {
			result = statistical
		}
	}

	return result
}

// inferFromCulturalMarkers uses culture-specific gender markers
func (e *Engine) inferFromCulturalMarkers(original, transliterated, culture, language string) *Inference {
	switch {
	case culture == "vietnamese" || language == "vi" || e.looksVietnamese(original):
		return e.inferVietnamese(original, transliterated)
		
	case culture == "arabic" || language == "ar" || e.looksArabic(original):
		return e.inferArabic(transliterated)
		
	case culture == "indonesian" || culture == "malaysian" || strings.Contains(language, "id") || strings.Contains(language, "ms"):
		return e.inferIndonesian(transliterated)
		
	case culture == "chinese" || language == "zh" || language == "zh-CN" || language == "zh-TW":
		return e.inferChinese(original, transliterated)
		
	case culture == "japanese" || language == "ja":
		return e.inferJapanese(original, transliterated)
		
	case culture == "korean" || language == "ko":
		return e.inferKorean(original, transliterated)
		
	case culture == "indian" || language == "hi" || language == "ta" || language == "te":
		return e.inferIndian(transliterated)
		
	case culture == "thai" || language == "th":
		return e.inferThai(transliterated)
		
	default:
		return e.inferWestern(transliterated, language)
	}
}

// inferVietnamese uses Vietnamese gender markers
func (e *Engine) inferVietnamese(original, transliterated string) *Inference {
	originalLower := strings.ToLower(original)
	transliteratedLower := strings.ToLower(transliterated)
	
	// Vietnamese gender markers in middle names
	if strings.Contains(originalLower, "văn") || strings.Contains(transliteratedLower, "van") {
		return &Inference{
			Value:      "M",
			Confidence: 0.85,
			Source:     "cultural_marker",
			Reason:     "Vietnamese marker 'Văn' typically indicates male",
		}
	}
	
	if strings.Contains(originalLower, "thị") || strings.Contains(transliteratedLower, "thi") {
		return &Inference{
			Value:      "F",
			Confidence: 0.85,
			Source:     "cultural_marker",
			Reason:     "Vietnamese marker 'Thị' typically indicates female",
		}
	}
	
	// Check for other Vietnamese gendered names
	maleMarkers := []string{"minh", "duc", "hoang", "quang", "thanh", "tuan", "hung", "dung", "phong"}
	femaleMarkers := []string{"linh", "mai", "lan", "yen", "huong", "ngoc", "thuy", "anh", "ha"}
	
	for _, marker := range maleMarkers {
		if strings.Contains(transliteratedLower, marker) {
			return &Inference{
				Value:      "M",
				Confidence: 0.65,
				Source:     "cultural_marker",
				Reason:     "Vietnamese name pattern suggests male",
			}
		}
	}
	
	for _, marker := range femaleMarkers {
		if strings.Contains(transliteratedLower, marker) {
			return &Inference{
				Value:      "F",
				Confidence: 0.65,
				Source:     "cultural_marker",
				Reason:     "Vietnamese name pattern suggests female",
			}
		}
	}
	
	return &Inference{Value: "X", Confidence: 0.1, Source: "unknown", Reason: "No Vietnamese gender markers found"}
}

// inferArabic uses Arabic patronymic indicators
func (e *Engine) inferArabic(text string) *Inference {
	textLower := strings.ToLower(text)
	
	if strings.Contains(textLower, "bin ") || strings.Contains(textLower, "ibn ") {
		return &Inference{
			Value:      "M",
			Confidence: 0.90,
			Source:     "cultural_marker",
			Reason:     "Arabic patronymic 'bin/ibn' (son of) indicates male",
		}
	}
	
	if strings.Contains(textLower, "bint ") || strings.Contains(textLower, "binte ") {
		return &Inference{
			Value:      "F",
			Confidence: 0.90,
			Source:     "cultural_marker",
			Reason:     "Arabic patronymic 'bint' (daughter of) indicates female",
		}
	}
	
	// Common Arabic gendered names
	maleNames := []string{"ahmad", "muhammad", "ali", "omar", "khalid", "hassan", "ibrahim", "yousef", "abdullah"}
	femaleNames := []string{"fatima", "aisha", "sarah", "mariam", "zahra", "layla", "amina", "khadija", "nour"}
	
	for _, name := range maleNames {
		if strings.Contains(textLower, name) {
			return &Inference{
				Value:      "M",
				Confidence: 0.75,
				Source:     "cultural_marker",
				Reason:     "Common Arabic male name pattern",
			}
		}
	}
	
	for _, name := range femaleNames {
		if strings.Contains(textLower, name) {
			return &Inference{
				Value:      "F",
				Confidence: 0.75,
				Source:     "cultural_marker",
				Reason:     "Common Arabic female name pattern",
			}
		}
	}
	
	return &Inference{Value: "X", Confidence: 0.1, Source: "unknown", Reason: "No Arabic gender markers found"}
}

// inferIndonesian uses Indonesian/Malaysian patronymic patterns
func (e *Engine) inferIndonesian(text string) *Inference {
	textLower := strings.ToLower(text)
	
	if strings.Contains(textLower, "bin ") {
		return &Inference{
			Value:      "M",
			Confidence: 0.88,
			Source:     "cultural_marker",
			Reason:     "Malay/Indonesian patronymic 'bin' (son of) indicates male",
		}
	}
	
	if strings.Contains(textLower, "binti ") || strings.Contains(textLower, "binte ") {
		return &Inference{
			Value:      "F",
			Confidence: 0.88,
			Source:     "cultural_marker",
			Reason:     "Malay/Indonesian patronymic 'binti' (daughter of) indicates female",
		}
	}
	
	// Indonesian gendered name patterns
	maleNames := []string{"ahmad", "muhammad", "adi", "budi", "eko", "hadi", "indra", "joko", "rudi"}
	femaleNames := []string{"sari", "dewi", "rina", "maya", "indah", "fitri", "wati", "ning", "sri"}
	
	for _, name := range maleNames {
		if strings.Contains(textLower, name) {
			return &Inference{
				Value:      "M",
				Confidence: 0.70,
				Source:     "cultural_marker",
				Reason:     "Indonesian male name pattern",
			}
		}
	}
	
	for _, name := range femaleNames {
		if strings.Contains(textLower, name) {
			return &Inference{
				Value:      "F",
				Confidence: 0.70,
				Source:     "cultural_marker",
				Reason:     "Indonesian female name pattern",
			}
		}
	}
	
	return &Inference{Value: "X", Confidence: 0.1, Source: "unknown", Reason: "No Indonesian gender markers found"}
}

// inferChinese uses Chinese name patterns (limited accuracy)
func (e *Engine) inferChinese(original, transliterated string) *Inference {
	// Chinese gender inference is very difficult and unreliable
	// We can only make very general observations
	
	textLower := strings.ToLower(transliterated)
	
	// Some traditionally male-associated characters (very low confidence)
	maleIndicators := []string{"jian", "ming", "wei", "gang", "jun", "qiang", "lei", "bin"}
	femaleIndicators := []string{"li", "mei", "hua", "yan", "hong", "ping", "na", "jing", "xue"}
	
	for _, indicator := range maleIndicators {
		if strings.Contains(textLower, indicator) {
			return &Inference{
				Value:      "M",
				Confidence: 0.55,
				Source:     "statistical",
				Reason:     "Chinese name element suggests male (low confidence)",
			}
		}
	}
	
	for _, indicator := range femaleIndicators {
		if strings.Contains(textLower, indicator) {
			return &Inference{
				Value:      "F",
				Confidence: 0.55,
				Source:     "statistical",
				Reason:     "Chinese name element suggests female (low confidence)",
			}
		}
	}
	
	return &Inference{Value: "X", Confidence: 0.1, Source: "unknown", Reason: "Chinese names require cultural knowledge for gender inference"}
}

// inferJapanese uses Japanese name patterns (limited)
func (e *Engine) inferJapanese(original, transliterated string) *Inference {
	textLower := strings.ToLower(transliterated)
	
	// Common Japanese name endings
	if strings.HasSuffix(textLower, "ko") || strings.HasSuffix(textLower, "mi") || strings.HasSuffix(textLower, "ka") {
		return &Inference{
			Value:      "F",
			Confidence: 0.70,
			Source:     "cultural_marker",
			Reason:     "Japanese name ending suggests female",
		}
	}
	
	if strings.HasSuffix(textLower, "ro") || strings.HasSuffix(textLower, "ta") || strings.HasSuffix(textLower, "ki") {
		return &Inference{
			Value:      "M",
			Confidence: 0.60,
			Source:     "cultural_marker",
			Reason:     "Japanese name ending suggests male",
		}
	}
	
	return &Inference{Value: "X", Confidence: 0.1, Source: "unknown", Reason: "Japanese gender inference requires cultural context"}
}

// inferKorean uses Korean name patterns (very limited)
func (e *Engine) inferKorean(original, transliterated string) *Inference {
	// Korean gender inference is extremely difficult without cultural knowledge
	return &Inference{
		Value:      "X",
		Confidence: 0.1,
		Source:     "unknown",
		Reason:     "Korean names require cultural knowledge for gender inference",
	}
}

// inferIndian uses Indian name patterns
func (e *Engine) inferIndian(text string) *Inference {
	textLower := strings.ToLower(text)
	
	// Common Indian male names
	maleNames := []string{"raj", "kumar", "singh", "dev", "krishna", "ram", "sharma", "gupta", "anil", "sunil"}
	femaleNames := []string{"devi", "kumari", "priya", "sita", "gita", "lata", "rani", "shanti", "maya", "radha"}
	
	for _, name := range maleNames {
		if strings.Contains(textLower, name) {
			return &Inference{
				Value:      "M",
				Confidence: 0.75,
				Source:     "cultural_marker",
				Reason:     "Indian male name pattern",
			}
		}
	}
	
	for _, name := range femaleNames {
		if strings.Contains(textLower, name) {
			return &Inference{
				Value:      "F",
				Confidence: 0.75,
				Source:     "cultural_marker",
				Reason:     "Indian female name pattern",
			}
		}
	}
	
	return &Inference{Value: "X", Confidence: 0.1, Source: "unknown", Reason: "No Indian gender markers found"}
}

// inferThai uses Thai name patterns
func (e *Engine) inferThai(text string) *Inference {
	// Thai names are difficult to gender without cultural knowledge
	return &Inference{
		Value:      "X",
		Confidence: 0.1,
		Source:     "unknown",
		Reason:     "Thai names require cultural knowledge for gender inference",
	}
}

// inferWestern uses Western name patterns and statistical data
func (e *Engine) inferWestern(text string, language string) *Inference {
	textLower := strings.ToLower(text)
	
	// Common Western gendered names
	maleNames := []string{"john", "david", "michael", "james", "robert", "william", "richard", "thomas", "mark", "daniel"}
	femaleNames := []string{"mary", "patricia", "jennifer", "linda", "elizabeth", "barbara", "susan", "jessica", "sarah", "karen"}
	
	// Check for exact matches first
	words := strings.Fields(textLower)
	for _, word := range words {
		for _, name := range maleNames {
			if word == name {
				return &Inference{
					Value:      "M",
					Confidence: 0.85,
					Source:     "statistical",
					Reason:     "Common Western male name",
				}
			}
		}
		
		for _, name := range femaleNames {
			if word == name {
				return &Inference{
					Value:      "F",
					Confidence: 0.85,
					Source:     "statistical",
					Reason:     "Common Western female name",
				}
			}
		}
	}
	
	// Check name endings (lower confidence)
	if e.useStatistical {
		for _, word := range words {
			if len(word) > 2 {
				// Female name endings
				if strings.HasSuffix(word, "a") || strings.HasSuffix(word, "ia") || strings.HasSuffix(word, "ina") {
					return &Inference{
						Value:      "F",
						Confidence: 0.60,
						Source:     "statistical",
						Reason:     "Name ending pattern suggests female",
					}
				}
				
				// Male name endings
				if strings.HasSuffix(word, "er") || strings.HasSuffix(word, "on") || strings.HasSuffix(word, "us") {
					return &Inference{
						Value:      "M",
						Confidence: 0.55,
						Source:     "statistical",
						Reason:     "Name ending pattern suggests male",
					}
				}
			}
		}
	}
	
	return &Inference{Value: "X", Confidence: 0.1, Source: "unknown", Reason: "No Western gender indicators found"}
}

// inferFromStatisticalPatterns uses statistical analysis (placeholder for more sophisticated methods)
func (e *Engine) inferFromStatisticalPatterns(text, culture, language string) *Inference {
	if !e.useStatistical {
		return &Inference{Value: "X", Confidence: 0.0, Source: "disabled"}
	}
	
	// This would integrate with statistical models trained on name data
	// For now, return low-confidence unknown
	return &Inference{
		Value:      "X",
		Confidence: 0.2,
		Source:     "statistical",
		Reason:     "Statistical analysis inconclusive",
	}
}

// Helper methods for cultural detection

func (e *Engine) looksVietnamese(text string) bool {
	markers := []string{"ă", "â", "đ", "ê", "ô", "ơ", "ư", "thị", "văn"}
	count := 0
	lower := strings.ToLower(text)
	for _, marker := range markers {
		if strings.Contains(lower, marker) {
			count++
		}
	}
	return count >= 2
}

func (e *Engine) looksArabic(text string) bool {
	for _, r := range text {
		if r >= 0x0600 && r <= 0x06FF {
			return true
		}
	}
	return false
}

// GetGenderFromTitle extracts gender information from titles
func GetGenderFromTitle(title string) *Inference {
	titleLower := strings.ToLower(strings.Trim(title, "."))
	
	switch titleLower {
	case "mr", "sir", "lord", "herr", "señor", "monsieur":
		return &Inference{
			Value:      "M",
			Confidence: 0.95,
			Source:     "cultural_marker",
			Reason:     "Male-specific title",
		}
		
	case "mrs", "ms", "miss", "lady", "dame", "frau", "señora", "señorita", "madame", "mademoiselle":
		return &Inference{
			Value:      "F",
			Confidence: 0.95,
			Source:     "cultural_marker",
			Reason:     "Female-specific title",
		}
		
	case "mx":
		return &Inference{
			Value:      "X",
			Confidence: 0.95,
			Source:     "cultural_marker",
			Reason:     "Gender-neutral title",
		}
		
	default:
		// Dr, Prof, Rev, etc. are gender-neutral
		return &Inference{
			Value:      "X",
			Confidence: 0.1,
			Source:     "unknown",
			Reason:     "Gender-neutral or unknown title",
		}
	}
}