// Package nameparser provides cultural name parsing and formatting capabilities.
package nameparser

import (
	"strings"
)

// NameStructure represents parsed name components with cultural awareness
type NameStructure struct {
	Family       string   `json:"family"`                // Family/surname (UPPERCASE for display)
	First        string   `json:"first"`                 // Given/first name (Title Case)
	Middle       []string `json:"middle,omitempty"`      // Middle names/patronymics
	Titles       []string `json:"titles,omitempty"`      // Extracted titles (Dr, Prof, etc)
	Suffixes     []string `json:"suffixes,omitempty"`    // Jr., Sr., III, etc.
	Particles    []string `json:"particles,omitempty"`   // de, van, von, del, etc.
	FullASCII    string   `json:"full_ascii"`            // Complete formatted ASCII name
	OriginalForm string   `json:"original_form"`         // Original input for reference
	Order        string   `json:"order"`                 // "western" or "eastern"
}

// CulturalContext provides information about naming conventions
type CulturalContext struct {
	Culture           string   `json:"culture"`             // "western", "chinese", "vietnamese", etc.
	NameOrder         string   `json:"name_order"`          // "family-first" or "given-first"
	HasGenderMarkers  bool     `json:"has_gender_markers"`  // Whether culture uses gender markers
	HasPatronymics    bool     `json:"has_patronymics"`     // Whether culture uses patronymics
	ParticlePrefix    bool     `json:"particle_prefix"`     // Whether particles come before surnames
	CaseSensitive     bool     `json:"case_sensitive"`      // Whether proper case is culturally important
	PreservedElements []string `json:"preserved_elements"`  // Elements that should not be altered
}

// Parser handles name parsing with cultural awareness
type Parser struct {
	preserveOriginal bool
	strictCultural   bool
}

// NewParser creates a new name parser
func NewParser(preserveOriginal, strictCultural bool) *Parser {
	return &Parser{
		preserveOriginal: preserveOriginal,
		strictCultural:   strictCultural,
	}
}

// ParseName analyzes and structures a name according to cultural conventions
func (p *Parser) ParseName(originalText, transliteratedText, culture, language string) *NameStructure {
	if transliteratedText == "" {
		return &NameStructure{
			OriginalForm: originalText,
			FullASCII:    "",
		}
	}

	// Extract titles first
	titles := p.extractTitles(transliteratedText)
	cleanText := p.removeTitles(transliteratedText, titles)

	// Extract suffixes
	suffixes := p.extractSuffixes(cleanText)
	cleanText = p.removeSuffixes(cleanText, suffixes)

	// Determine cultural context
	context := p.getCulturalContext(culture, language, originalText)

	// Parse according to cultural conventions
	var result *NameStructure
	switch context.Culture {
	case "vietnamese":
		result = p.parseVietnamese(originalText, cleanText, context)
	case "chinese":
		result = p.parseChinese(cleanText, context)
	case "japanese":
		result = p.parseJapanese(cleanText, context)
	case "arabic":
		result = p.parseArabic(cleanText, context)
	case "korean":
		result = p.parseKorean(cleanText, context)
	case "indian":
		result = p.parseIndian(cleanText, context)
	case "indonesian", "malaysian":
		result = p.parseIndonesian(cleanText, context)
	case "thai":
		result = p.parseThai(cleanText, context)
	default:
		result = p.parseWestern(cleanText, context)
	}

	// Add metadata
	result.Titles = titles
	result.Suffixes = suffixes
	result.OriginalForm = originalText
	result.Order = context.NameOrder
	result.FullASCII = p.formatFullName(result, context)

	return result
}

// extractTitles identifies and extracts titles from text
func (p *Parser) extractTitles(text string) []string {
	titleMapping := map[string]string{
		// English titles
		"dr": "Dr", "doctor": "Dr", "prof": "Prof", "professor": "Prof",
		"mr": "Mr", "mrs": "Mrs", "ms": "Ms", "miss": "Miss", "mx": "Mx",
		"sir": "Sir", "dame": "Dame", "lord": "Lord", "lady": "Lady",
		"hon": "Hon", "honourable": "Hon", "rev": "Rev", "reverend": "Rev",
		
		// Academic/Professional
		"phd": "PhD", "md": "MD", "jd": "JD", "esq": "Esq",
		
		// International variants
		"herr": "Mr", "frau": "Mrs", "fraulein": "Ms",
		"señor": "Mr", "señora": "Mrs", "señorita": "Ms",
		"monsieur": "Mr", "madame": "Mrs", "mademoiselle": "Ms",
	}

	var titles []string
	words := strings.Fields(text)

	for _, word := range words {
		cleaned := strings.ToLower(strings.Trim(word, ".,"))
		if title, exists := titleMapping[cleaned]; exists {
			titles = append(titles, title)
		}
	}

	return titles
}

// extractSuffixes identifies generational and other suffixes
func (p *Parser) extractSuffixes(text string) []string {
	suffixMapping := map[string]string{
		"jr": "Jr", "junior": "Jr", "sr": "Sr", "senior": "Sr",
		"ii": "II", "iii": "III", "iv": "IV", "v": "V",
		"2nd": "II", "3rd": "III", "4th": "IV", "5th": "V",
		"phd": "PhD", "md": "MD", "esq": "Esq",
	}

	var suffixes []string
	words := strings.Fields(text)

	// Check last few words for suffixes
	for i := len(words) - 1; i >= 0 && len(suffixes) < 3; i-- {
		cleaned := strings.ToLower(strings.Trim(words[i], ".,"))
		if suffix, exists := suffixMapping[cleaned]; exists {
			suffixes = append([]string{suffix}, suffixes...) // Prepend to maintain order
		} else if len(suffixes) > 0 {
			break // Stop if we hit a non-suffix after finding suffixes
		}
	}

	return suffixes
}

// removeTitles removes identified titles from text
func (p *Parser) removeTitles(text string, titles []string) string {
	if len(titles) == 0 {
		return text
	}

	words := strings.Fields(text)
	var resultWords []string

	titleSet := make(map[string]bool)
	for _, title := range titles {
		titleSet[strings.ToLower(strings.Trim(title, "."))] = true
	}

	for _, word := range words {
		cleaned := strings.ToLower(strings.Trim(word, ".,"))
		if !titleSet[cleaned] {
			resultWords = append(resultWords, word)
		}
	}

	return strings.TrimSpace(strings.Join(resultWords, " "))
}

// removeSuffixes removes identified suffixes from text
func (p *Parser) removeSuffixes(text string, suffixes []string) string {
	if len(suffixes) == 0 {
		return text
	}

	words := strings.Fields(text)
	suffixSet := make(map[string]bool)
	for _, suffix := range suffixes {
		suffixSet[strings.ToLower(strings.Trim(suffix, "."))] = true
	}

	// Remove suffixes from the end
	for len(words) > 0 {
		lastWord := strings.ToLower(strings.Trim(words[len(words)-1], ".,"))
		if suffixSet[lastWord] {
			words = words[:len(words)-1]
		} else {
			break
		}
	}

	return strings.TrimSpace(strings.Join(words, " "))
}

// getCulturalContext determines naming conventions based on culture/language
func (p *Parser) getCulturalContext(culture, language, originalText string) CulturalContext {
	switch {
	case culture == "vietnamese" || language == "vi" || p.looksVietnamese(originalText):
		return CulturalContext{
			Culture:          "vietnamese",
			NameOrder:        "family-first",
			HasGenderMarkers: true,
			ParticlePrefix:   false,
			CaseSensitive:    false,
		}
		
	case culture == "chinese" || language == "zh" || language == "zh-CN" || language == "zh-TW" || p.looksChinese(originalText):
		return CulturalContext{
			Culture:       "chinese",
			NameOrder:     "family-first",
			CaseSensitive: true,
		}
		
	case culture == "japanese" || language == "ja" || p.looksJapanese(originalText):
		return CulturalContext{
			Culture:       "japanese",
			NameOrder:     "family-first",
			CaseSensitive: true,
		}
		
	case culture == "korean" || language == "ko" || p.looksKorean(originalText):
		return CulturalContext{
			Culture:       "korean",
			NameOrder:     "family-first",
			CaseSensitive: true,
		}
		
	case culture == "arabic" || language == "ar" || p.looksArabic(originalText):
		return CulturalContext{
			Culture:        "arabic",
			NameOrder:      "given-first",
			HasPatronymics: true,
			ParticlePrefix: false,
		}
		
	case strings.Contains(language, "in") || culture == "indonesian" || culture == "malaysian":
		return CulturalContext{
			Culture:        "indonesian",
			NameOrder:      "given-first",
			HasPatronymics: true,
		}
		
	default:
		return CulturalContext{
			Culture:        "western",
			NameOrder:      "given-first",
			ParticlePrefix: true,
			CaseSensitive:  true,
		}
	}
}

// parseVietnamese handles Vietnamese naming conventions
func (p *Parser) parseVietnamese(original, text string, context CulturalContext) *NameStructure {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return &NameStructure{}
	}

	var result NameStructure
	
	if len(parts) >= 2 {
		// Vietnamese: Family name first, then middle names, then given name
		result.Family = strings.ToUpper(parts[0])
		result.First = p.toTitleCase(parts[len(parts)-1])
		
		// Handle middle names (excluding gender markers)
		for i := 1; i < len(parts)-1; i++ {
			part := parts[i]
			partLower := strings.ToLower(part)
			
			// Vietnamese gender markers - preserve but don't capitalize
			if partLower == "van" || partLower == "văn" || partLower == "thi" || partLower == "thị" {
				result.Middle = append(result.Middle, p.toTitleCase(part))
			} else {
				result.Middle = append(result.Middle, p.toTitleCase(part))
			}
		}
	} else {
		result.First = p.toTitleCase(parts[0])
	}

	return &result
}

// parseChinese handles Chinese naming conventions
func (p *Parser) parseChinese(text string, context CulturalContext) *NameStructure {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return &NameStructure{}
	}

	var result NameStructure

	if len(parts) >= 2 {
		// Chinese: Family name first, then given names
		result.Family = strings.ToUpper(parts[0])
		if len(parts) == 2 {
			result.First = p.toTitleCase(parts[1])
		} else {
			// Multiple given names - last one is primary
			result.First = p.toTitleCase(parts[len(parts)-1])
			for i := 1; i < len(parts)-1; i++ {
				result.Middle = append(result.Middle, p.toTitleCase(parts[i]))
			}
		}
	} else {
		result.First = p.toTitleCase(parts[0])
	}

	return &result
}

// parseJapanese handles Japanese naming conventions
func (p *Parser) parseJapanese(text string, context CulturalContext) *NameStructure {
	// Remove honorifics like -san, -kun, -chan
	text = p.removeJapaneseHonorifics(text)
	
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return &NameStructure{}
	}

	var result NameStructure

	if len(parts) >= 2 {
		// Japanese: Family name first, then given name
		result.Family = strings.ToUpper(parts[0])
		result.First = p.toTitleCase(parts[len(parts)-1])
		
		// Middle names are rare in Japanese
		for i := 1; i < len(parts)-1; i++ {
			result.Middle = append(result.Middle, p.toTitleCase(parts[i]))
		}
	} else {
		result.First = p.toTitleCase(parts[0])
	}

	return &result
}

// parseArabic handles Arabic naming conventions
func (p *Parser) parseArabic(text string, context CulturalContext) *NameStructure {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return &NameStructure{}
	}

	var result NameStructure

	if len(parts) >= 2 {
		// Arabic: Given name first, patronymics/family name last
		result.First = p.toTitleCase(parts[0])
		result.Family = strings.ToUpper(parts[len(parts)-1])

		// Handle patronymics (ibn, bin, bint, etc.)
		for i := 1; i < len(parts)-1; i++ {
			part := parts[i]
			result.Middle = append(result.Middle, p.toTitleCase(part))
		}
	} else {
		result.First = p.toTitleCase(parts[0])
	}

	return &result
}

// parseKorean handles Korean naming conventions
func (p *Parser) parseKorean(text string, context CulturalContext) *NameStructure {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return &NameStructure{}
	}

	var result NameStructure

	if len(parts) >= 2 {
		// Korean: Family name first, then given name
		result.Family = strings.ToUpper(parts[0])
		if len(parts) == 2 {
			result.First = p.toTitleCase(parts[1])
		} else {
			// Multiple given names
			result.First = p.toTitleCase(parts[len(parts)-1])
			for i := 1; i < len(parts)-1; i++ {
				result.Middle = append(result.Middle, p.toTitleCase(parts[i]))
			}
		}
	} else {
		result.First = p.toTitleCase(parts[0])
	}

	return &result
}

// parseIndian handles Indian naming conventions
func (p *Parser) parseIndian(text string, context CulturalContext) *NameStructure {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return &NameStructure{}
	}

	// Indian names are very diverse - use Western pattern as default
	return p.parseWestern(text, context)
}

// parseIndonesian handles Indonesian/Malaysian naming conventions
func (p *Parser) parseIndonesian(text string, context CulturalContext) *NameStructure {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return &NameStructure{}
	}

	var result NameStructure

	// Check for patronymic indicators
	hasPatronymic := false
	for _, part := range parts {
		lower := strings.ToLower(part)
		if lower == "bin" || lower == "binti" || lower == "binte" {
			hasPatronymic = true
			break
		}
	}

	if len(parts) == 1 {
		// Mononym
		result.First = p.toTitleCase(parts[0])
	} else if hasPatronymic {
		// Patronymic structure
		result.First = p.toTitleCase(parts[0])
		for i := 1; i < len(parts); i++ {
			result.Middle = append(result.Middle, p.toTitleCase(parts[i]))
		}
	} else {
		// Regular multi-part name
		result.First = p.toTitleCase(parts[0])
		if len(parts) > 1 {
			result.Family = strings.ToUpper(parts[len(parts)-1])
		}
		for i := 1; i < len(parts)-1; i++ {
			result.Middle = append(result.Middle, p.toTitleCase(parts[i]))
		}
	}

	return &result
}

// parseThai handles Thai naming conventions
func (p *Parser) parseThai(text string, context CulturalContext) *NameStructure {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return &NameStructure{}
	}

	var result NameStructure

	if len(parts) >= 2 {
		// Thai: Given name first, family name last
		result.First = p.toTitleCase(parts[0])
		result.Family = strings.ToUpper(parts[len(parts)-1])
		
		for i := 1; i < len(parts)-1; i++ {
			result.Middle = append(result.Middle, p.toTitleCase(parts[i]))
		}
	} else {
		result.First = p.toTitleCase(parts[0])
	}

	return &result
}

// parseWestern handles Western naming conventions
func (p *Parser) parseWestern(text string, context CulturalContext) *NameStructure {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return &NameStructure{}
	}

	var result NameStructure

	// Extract particles (de, van, von, del, etc.)
	particles, cleanParts := p.extractParticles(parts)
	result.Particles = particles

	if len(cleanParts) == 1 {
		result.First = p.toTitleCase(cleanParts[0])
	} else if len(cleanParts) == 2 {
		result.First = p.toTitleCase(cleanParts[0])
		result.Family = strings.ToUpper(cleanParts[1])
	} else {
		// First + middle names + last
		result.First = p.toTitleCase(cleanParts[0])
		result.Family = strings.ToUpper(cleanParts[len(cleanParts)-1])
		for i := 1; i < len(cleanParts)-1; i++ {
			result.Middle = append(result.Middle, p.toTitleCase(cleanParts[i]))
		}
	}

	return &result
}

// extractParticles identifies and extracts nobiliary particles
func (p *Parser) extractParticles(parts []string) ([]string, []string) {
	particleSet := map[string]bool{
		"de": true, "del": true, "della": true, "di": true, "da": true,
		"van": true, "von": true, "der": true, "den": true, "ter": true,
		"le": true, "la": true, "du": true, "des": true,
		"bin": true, "binti": true, "ibn": true, "bint": true,
		"al": true, "el": true,
	}

	var particles []string
	var cleanParts []string

	for i, part := range parts {
		lower := strings.ToLower(part)
		if particleSet[lower] && i > 0 && i < len(parts)-1 {
			// Keep particles with family name
			particles = append(particles, strings.ToLower(part))
		} else {
			cleanParts = append(cleanParts, part)
		}
	}

	// If we found particles, add them back to the family name
	if len(particles) > 0 && len(cleanParts) > 1 {
		familyIndex := len(cleanParts) - 1
		familyParts := append(particles, cleanParts[familyIndex])
		cleanParts[familyIndex] = strings.Join(familyParts, " ")
	}

	return particles, cleanParts
}

// removeJapaneseHonorifics removes Japanese honorific suffixes
func (p *Parser) removeJapaneseHonorifics(text string) string {
	honorifics := []string{"-san", "-kun", "-chan", "-sama", "-sensei", "-senpai"}
	
	for _, honorific := range honorifics {
		text = strings.ReplaceAll(text, honorific, "")
	}
	
	return strings.TrimSpace(text)
}

// formatFullName creates the complete formatted ASCII name
func (p *Parser) formatFullName(name *NameStructure, context CulturalContext) string {
	var parts []string

	// Add titles
	for _, title := range name.Titles {
		parts = append(parts, title)
	}

	// Add name components based on cultural order
	if context.NameOrder == "family-first" {
		if name.Family != "" {
			parts = append(parts, name.Family)
		}
		if name.First != "" {
			parts = append(parts, name.First)
		}
		for _, middle := range name.Middle {
			if middle != "" {
				parts = append(parts, middle)
			}
		}
	} else {
		// Given-first order
		if name.First != "" {
			parts = append(parts, name.First)
		}
		for _, middle := range name.Middle {
			if middle != "" {
				parts = append(parts, middle)
			}
		}
		if name.Family != "" {
			parts = append(parts, name.Family)
		}
	}

	// Add suffixes
	for _, suffix := range name.Suffixes {
		parts = append(parts, suffix)
	}

	return strings.Join(parts, " ")
}

// toTitleCase converts text to title case
func (p *Parser) toTitleCase(text string) string {
	return strings.Title(strings.ToLower(text))
}

// Helper methods for cultural detection

func (p *Parser) looksVietnamese(text string) bool {
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

func (p *Parser) looksChinese(text string) bool {
	for _, r := range text {
		if r >= 0x4E00 && r <= 0x9FFF {
			return true
		}
	}
	return false
}

func (p *Parser) looksJapanese(text string) bool {
	for _, r := range text {
		if (r >= 0x3040 && r <= 0x309F) || (r >= 0x30A0 && r <= 0x30FF) {
			return true
		}
	}
	return false
}

func (p *Parser) looksKorean(text string) bool {
	for _, r := range text {
		if r >= 0xAC00 && r <= 0xD7AF {
			return true
		}
	}
	return false
}

func (p *Parser) looksArabic(text string) bool {
	for _, r := range text {
		if r >= 0x0600 && r <= 0x06FF {
			return true
		}
	}
	return false
}