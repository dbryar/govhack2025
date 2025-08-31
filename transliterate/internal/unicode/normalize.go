// Package unicode provides Unicode normalization and character handling utilities.
package unicode

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// NormalizeOptions configures Unicode normalization behavior
type NormalizeOptions struct {
	Form           norm.Form // NFC, NFD, NFKC, NFKD
	RemoveDiacritics bool     // Remove combining diacritical marks
	CaseFolding     bool     // Apply case folding for comparison
	ASCIIOnly       bool     // Convert to ASCII-compatible characters
}

// DefaultNormalizeOptions provides sensible defaults for most use cases
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		Form:           norm.NFD,
		RemoveDiacritics: false,
		CaseFolding:     false,
		ASCIIOnly:       false,
	}
}

// NormalizeText applies Unicode normalization according to the specified options
func NormalizeText(text string, opts NormalizeOptions) (string, error) {
	if !utf8.ValidString(text) {
		return "", ErrInvalidUTF8
	}

	// Start with the specified normalization form
	result := opts.Form.String(text)

	// Create transformation chain
	var transformations []transform.Transformer

	// Add normalization
	transformations = append(transformations, opts.Form)

	// Remove diacritics if requested
	if opts.RemoveDiacritics {
		transformations = append(transformations, runes.Remove(runes.In(unicode.Mn)))
	}

	// Apply case folding
	if opts.CaseFolding {
		transformations = append(transformations, runes.Map(unicode.ToLower))
	}

	// ASCII conversion if requested
	if opts.ASCIIOnly {
		transformations = append(transformations, NewASCIITransformer())
	}

	// Apply all transformations
	if len(transformations) > 1 {
		chain := transform.Chain(transformations...)
		var err error
		result, _, err = transform.String(chain, text)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}

// StripDiacritics removes diacritical marks while preserving base characters
func StripDiacritics(text string) (string, error) {
	opts := NormalizeOptions{
		Form:           norm.NFD,
		RemoveDiacritics: true,
	}
	return NormalizeText(text, opts)
}

// ToASCII converts text to ASCII-compatible form with intelligent character mapping
func ToASCII(text string) (string, error) {
	if !utf8.ValidString(text) {
		return "", ErrInvalidUTF8
	}

	var result strings.Builder
	result.Grow(len(text))

	for _, r := range text {
		ascii := mapToASCII(r)
		result.WriteString(ascii)
	}

	return result.String(), nil
}

// mapToASCII maps a single rune to its ASCII representation
func mapToASCII(r rune) string {
	// Already ASCII
	if r < 128 {
		return string(r)
	}

	// Language-specific mappings (preserve cultural accuracy)
	if mapped := getLanguageSpecificASCII(r); mapped != "" {
		return mapped
	}

	// Diacritical mark removal
	if mapped := getDiacriticalMapping(r); mapped != "" {
		return mapped
	}

	// Script-specific conversions
	if mapped := getScriptMapping(r); mapped != "" {
		return mapped
	}

	// Fallback for letters
	if unicode.IsLetter(r) {
		// Try to get the base character by normalizing and removing marks
		normalized := norm.NFD.String(string(r))
		for _, nr := range normalized {
			if nr < 128 && unicode.IsLetter(nr) {
				return string(nr)
			}
		}
		return "?" // Unknown letter
	}

	// Handle other character types
	if unicode.IsDigit(r) {
		// Try to extract digit value
		if digit := unicode.ToLower(r); digit >= '0' && digit <= '9' {
			return string(digit)
		}
		return "0"
	}

	if unicode.IsSpace(r) {
		return " "
	}

	if unicode.IsPunct(r) {
		return getPunctuationMapping(r)
	}

	// Skip or replace with empty string for other characters
	return ""
}

// getLanguageSpecificASCII provides culturally-aware ASCII mappings
func getLanguageSpecificASCII(r rune) string {
	mappings := map[rune]string{
		// German umlauts and ß
		'Ä': "AE", 'ä': "ae",
		'Ö': "OE", 'ö': "oe", 
		'Ü': "UE", 'ü': "ue",
		'ß': "ss",

		// Scandinavian
		'Å': "AA", 'å': "aa",
		'Æ': "AE", 'æ': "ae",
		'Ø': "OE", 'ø': "oe",

		// Dutch/Flemish
		'ĳ': "ij", 'Ĳ': "IJ",

		// French ligatures
		'œ': "oe", 'Œ': "OE",

		// Spanish
		'Ñ': "N", 'ñ': "n",

		// Portuguese
		'ã': "a", 'Ã': "A",
		'õ': "o", 'Õ': "O",

		// Czech/Slovak
		'č': "c", 'Č': "C",
		'š': "s", 'Š': "S", 
		'ž': "z", 'Ž': "Z",
		'ř': "r", 'Ř': "R",

		// Polish
		'ł': "l", 'Ł': "L",
		'ą': "a", 'Ą': "A",
		'ę': "e", 'Ę': "E",
		'ć': "c", 'Ć': "C",
		'ń': "n", 'Ń': "N",
		'ś': "s", 'Ś': "S",
		'ź': "z", 'Ź': "Z",
		'ż': "z", 'Ż': "Z",

		// Vietnamese (preserve base characters)
		'Đ': "D", 'đ': "d",

		// Turkish
		'ı': "i", 'İ': "I",
		'ğ': "g", 'Ğ': "G",
		'ş': "s", 'Ş': "S",
		'ç': "c", 'Ç': "C",

		// Romanian
		'ă': "a", 'Ă': "A",
		'â': "a", 'Â': "A", 
		'î': "i", 'Î': "I",
		'ș': "s", 'Ș': "S",
		'ț': "t", 'Ț': "T",
	}

	return mappings[r]
}

// getDiacriticalMapping handles basic diacritical marks
func getDiacriticalMapping(r rune) string {
	// Common diacritical mappings
	mappings := map[rune]string{
		// A variants
		'À': "A", 'Á': "A", 'Â': "A", 'Ã': "A", 'Ā': "A", 'Ă': "A",
		'à': "a", 'á': "a", 'â': "a", 'ã': "a", 'ā': "a", 'ă': "a",

		// E variants  
		'È': "E", 'É': "E", 'Ê': "E", 'Ë': "E", 'Ē': "E", 'Ĕ': "E",
		'è': "e", 'é': "e", 'ê': "e", 'ë': "e", 'ē': "e", 'ĕ': "e",

		// I variants
		'Ì': "I", 'Í': "I", 'Î': "I", 'Ï': "I", 'Ī': "I", 'Ĭ': "I",
		'ì': "i", 'í': "i", 'î': "i", 'ï': "i", 'ī': "i", 'ĭ': "i",

		// O variants
		'Ò': "O", 'Ó': "O", 'Ô': "O", 'Õ': "O", 'Ō': "O", 'Ŏ': "O",
		'ò': "o", 'ó': "o", 'ô': "o", 'õ': "o", 'ō': "o", 'ŏ': "o",

		// U variants
		'Ù': "U", 'Ú': "U", 'Û': "U", 'Ū': "U", 'Ŭ': "U",
		'ù': "u", 'ú': "u", 'û': "u", 'ū': "u", 'ŭ': "u",

		// Y variants
		'Ỳ': "Y", 'Ý': "Y", 'Ŷ': "Y", 'Ÿ': "Y",
		'ỳ': "y", 'ý': "y", 'ŷ': "y", 'ÿ': "y",

		// Other common accented characters
		'Ç': "C", 'ç': "c",
	}

	return mappings[r]
}

// getScriptMapping handles script-to-Latin conversions
func getScriptMapping(r rune) string {
	// This is a simplified version - in practice you'd want more comprehensive mappings
	// For now, return empty to let the fallback handle it
	return ""
}

// getPunctuationMapping maps punctuation to ASCII equivalents
func getPunctuationMapping(r rune) string {
	mappings := map[rune]string{
		0x201C: "\"", 0x201D: "\"", // Smart quotes
		0x2018: "'", 0x2019: "'",   // Smart apostrophes
		0x2026: "...",              // Ellipsis
		0x2013: "-", 0x2014: "-",   // En dash, em dash
		0x00AB: "\"", 0x00BB: "\"", // Guillemets
		0x2039: "'", 0x203A: "'",   // Single guillemets
		0x2022: "*",                // Bullet
		0x00B7: ".",                // Middle dot
		0x00A1: "!",                // Inverted exclamation
		0x00BF: "?",                // Inverted question
	}

	if mapped, exists := mappings[r]; exists {
		return mapped
	}

	// Default punctuation handling
	return "."
}

// NewASCIITransformer creates a transformer that converts text to ASCII
func NewASCIITransformer() transform.Transformer {
	return &asciiTransformer{}
}

// asciiTransformer implements transform.Transformer for ASCII conversion
type asciiTransformer struct{}

func (t *asciiTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	for nSrc < len(src) {
		r, size := utf8.DecodeRune(src[nSrc:])
		if r == utf8.RuneError && size == 1 {
			if !atEOF {
				break
			}
			// Invalid UTF-8
			return nDst, nSrc, transform.ErrShortSrc
		}

		ascii := mapToASCII(r)
		if len(ascii) > len(dst)-nDst {
			break
		}

		copy(dst[nDst:], ascii)
		nDst += len(ascii)
		nSrc += size
	}

	if nSrc < len(src) {
		return nDst, nSrc, transform.ErrShortDst
	}

	return nDst, nSrc, nil
}

func (t *asciiTransformer) Reset() {}

// Custom errors
var (
	ErrInvalidUTF8 = transform.ErrShortSrc
)