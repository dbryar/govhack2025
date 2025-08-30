package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

// ValidationTestCase represents a test case for API validation
type ValidationTestCase struct {
	Name        string
	Input       string
	ExpectedOut ExpectedOutput
}

// ExpectedOutput represents the expected API response structure
type ExpectedOutput struct {
	Title  string
	Family string
	First  string
	Middle []string
	Gender struct {
		Value      string
		Confidence float64
	}
	FullASCII string
	NoUnicode bool // Should not contain any non-ASCII characters
	NoQuestion bool // Should not contain ? characters
}

// APIResponse matches the expected API response structure
type APIResponse struct {
	ID              string           `json:"id"`
	InputText       string           `json:"input_text"`
	OutputText      string           `json:"output_text"`
	InputScript     string           `json:"input_script"`
	OutputScript    string           `json:"output_script"`
	ConfidenceScore *float64         `json:"confidence_score"`
	Name            *NameStructure   `json:"name,omitempty"`
	Gender          *GenderInference `json:"gender,omitempty"`
}

type NameStructure struct {
	Family    string   `json:"family"`
	First     string   `json:"first"`
	Middle    []string `json:"middle,omitempty"`
	Titles    []string `json:"titles,omitempty"`
	FullASCII string   `json:"full_ascii"`
}

type GenderInference struct {
	Value      string  `json:"value"`
	Confidence float64 `json:"confidence"`
	Source     string  `json:"source"`
}

// Test cases from idea.md lines 733-737
var validationTestCases = []ValidationTestCase{
	{
		Name:  "Vietnamese male with title",
		Input: "Doctor Nguyễn Văn Minh",
		ExpectedOut: ExpectedOutput{
			Title:  "Dr.",
			Family: "NGUYEN", 
			First:  "Minh",
			Middle: []string{"Van"},
			Gender: struct {
				Value      string
				Confidence float64
			}{
				Value:      "M",
				Confidence: 0.6,
			},
			FullASCII: "DR NGUYEN MINH VAN",
			NoUnicode: true,
			NoQuestion: true,
		},
	},
	{
		Name:  "German with umlaut and eszett",
		Input: "Prof. Jürgen Groß",
		ExpectedOut: ExpectedOutput{
			Title:  "Prof.",
			Family: "GROSS",
			First:  "Jurgen",
			Middle: []string{},
			Gender: struct {
				Value      string
				Confidence float64
			}{
				Value:      "U", // Unknown
				Confidence: 0.0,
			},
			FullASCII: "PROF JURGEN GROSS",
			NoUnicode: true,
			NoQuestion: true,
		},
	},
	{
		Name:  "Chinese name",
		Input: "李小龍",
		ExpectedOut: ExpectedOutput{
			Title:  "",
			Family: "LI",
			First:  "Xiaolong",
			Middle: []string{},
			Gender: struct {
				Value      string
				Confidence float64
			}{
				Value:      "U", // Unknown unless we have specific markers
				Confidence: 0.0,
			},
			FullASCII: "LI XIAOLONG",
			NoUnicode: true,
			NoQuestion: true,
		},
	},
	{
		Name:  "Japanese honorific (should not become title)",
		Input: "Tanaka-san Yoko",
		ExpectedOut: ExpectedOutput{
			Title:  "", // -san should NOT be mapped to title
			Family: "TANAKA",
			First:  "Yoko",
			Middle: []string{},
			Gender: struct {
				Value      string
				Confidence float64
			}{
				Value:      "U",
				Confidence: 0.0,
			},
			FullASCII: "TANAKA YOKO", // No -san in output
			NoUnicode: true,
			NoQuestion: true,
		},
	},
	{
		Name:  "Spanish with particle and accents",
		Input: "Maria del Carmen Núñez",
		ExpectedOut: ExpectedOutput{
			Title:  "",
			Family: "DEL CARMEN NUNEZ", // Keep del particle with family
			First:  "Maria",
			Middle: []string{},
			Gender: struct {
				Value      string
				Confidence float64
			}{
				Value:      "F", // Could infer from Maria
				Confidence: 0.5, // Lower confidence for statistical inference
			},
			FullASCII: "MARIA DEL CARMEN NUNEZ",
			NoUnicode: true,
			NoQuestion: true,
		},
	},
}

const baseURL = "http://localhost:4000"

// TestValidationCases runs all validation test cases against the running API
func TestValidationCases(t *testing.T) {
	for _, testCase := range validationTestCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// Call the transliterate API
			response, err := callTransliterateAPI(testCase.Input)
			if err != nil {
				t.Fatalf("API call failed: %v", err)
			}

			// Validate response structure
			validateResponse(t, testCase, response)
		})
	}
}

// callTransliterateAPI makes a request to the transliterate API
func callTransliterateAPI(input string) (*APIResponse, error) {
	// Prepare request payload
	payload := map[string]interface{}{
		"text":          input,
		"output_script": "ascii",
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Make HTTP request
	resp, err := http.Post(baseURL+"/transliterate", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, body)
	}

	// Parse response
	var apiResponse APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiResponse, nil
}

// validateResponse validates the API response against expected output
func validateResponse(t *testing.T, testCase ValidationTestCase, response *APIResponse) {
	expected := testCase.ExpectedOut

	// Check for Unicode characters (should be ASCII only)
	if expected.NoUnicode {
		if !isASCIIOnly(response.OutputText) {
			t.Errorf("Output contains non-ASCII characters: %q", response.OutputText)
		}
		
		if response.Name != nil {
			if !isASCIIOnly(response.Name.FullASCII) {
				t.Errorf("FullASCII contains non-ASCII characters: %q", response.Name.FullASCII)
			}
			if !isASCIIOnly(response.Name.Family) {
				t.Errorf("Family name contains non-ASCII characters: %q", response.Name.Family)
			}
			if !isASCIIOnly(response.Name.First) {
				t.Errorf("First name contains non-ASCII characters: %q", response.Name.First)
			}
			for _, middle := range response.Name.Middle {
				if !isASCIIOnly(middle) {
					t.Errorf("Middle name contains non-ASCII characters: %q", middle)
				}
			}
		}
	}

	// Check for question marks (should not contain ?)
	if expected.NoQuestion {
		if strings.Contains(response.OutputText, "?") {
			t.Errorf("Output contains question marks: %q", response.OutputText)
		}
		
		if response.Name != nil {
			if strings.Contains(response.Name.FullASCII, "?") {
				t.Errorf("FullASCII contains question marks: %q", response.Name.FullASCII)
			}
		}
	}

	// Validate structured name parsing
	if response.Name == nil {
		t.Fatal("Expected Name structure to be present")
	}

	// Check that titles are not in name fields
	validateTitlesNotInNameFields(t, response.Name)

	// Validate family name
	if expected.Family != "" {
		if response.Name.Family != expected.Family {
			t.Errorf("Family name mismatch: got %q, want %q", response.Name.Family, expected.Family)
		}
	}

	// Validate first name  
	if expected.First != "" {
		if response.Name.First != expected.First {
			t.Errorf("First name mismatch: got %q, want %q", response.Name.First, expected.First)
		}
	}

	// Validate middle names
	if len(expected.Middle) > 0 {
		if len(response.Name.Middle) != len(expected.Middle) {
			t.Errorf("Middle names count mismatch: got %d, want %d", len(response.Name.Middle), len(expected.Middle))
		} else {
			for i, expectedMiddle := range expected.Middle {
				if i < len(response.Name.Middle) && response.Name.Middle[i] != expectedMiddle {
					t.Errorf("Middle name[%d] mismatch: got %q, want %q", i, response.Name.Middle[i], expectedMiddle)
				}
			}
		}
	}

	// Validate titles (should be in titles field, not name fields)
	if expected.Title != "" {
		if len(response.Name.Titles) == 0 {
			t.Errorf("Expected title %q but got no titles", expected.Title)
		} else {
			found := false
			for _, title := range response.Name.Titles {
				if strings.EqualFold(strings.TrimSuffix(title, "."), strings.TrimSuffix(expected.Title, ".")) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected title %q not found in titles %v", expected.Title, response.Name.Titles)
			}
		}
	}

	// Validate gender inference
	if response.Gender != nil {
		if expected.Gender.Value != "" {
			if response.Gender.Value != expected.Gender.Value {
				t.Errorf("Gender value mismatch: got %q, want %q", response.Gender.Value, expected.Gender.Value)
			}
		}
		
		if expected.Gender.Confidence > 0 {
			if response.Gender.Confidence < expected.Gender.Confidence {
				t.Errorf("Gender confidence too low: got %f, want >= %f", response.Gender.Confidence, expected.Gender.Confidence)
			}
		}
	}

	// Validate FullASCII format
	if expected.FullASCII != "" {
		// Allow some flexibility in formatting but check key components are present
		actualNormalized := normalizeFullASCII(response.Name.FullASCII)
		expectedNormalized := normalizeFullASCII(expected.FullASCII)
		
		if actualNormalized != expectedNormalized {
			t.Errorf("FullASCII mismatch: got %q, want %q", response.Name.FullASCII, expected.FullASCII)
		}
	}
}

// validateTitlesNotInNameFields ensures titles like "DR" are not present in first/middle/family name fields
func validateTitlesNotInNameFields(t *testing.T, name *NameStructure) {
	titlePatterns := []string{
		"DR", "DOCTOR", "PROF", "PROFESSOR", "MR", "MRS", "MS", "MISS",
		"SIR", "DAME", "LORD", "LADY", "HON", "REV",
	}

	for _, pattern := range titlePatterns {
		// Check first name
		if strings.Contains(strings.ToUpper(name.First), pattern) {
			t.Errorf("Title %q found in first name: %q", pattern, name.First)
		}
		
		// Check family name
		if strings.Contains(strings.ToUpper(name.Family), pattern) {
			t.Errorf("Title %q found in family name: %q", pattern, name.Family)
		}
		
		// Check middle names
		for i, middle := range name.Middle {
			if strings.Contains(strings.ToUpper(middle), pattern) {
				t.Errorf("Title %q found in middle name[%d]: %q", pattern, i, middle)
			}
		}
	}
}

// isASCIIOnly checks if a string contains only ASCII characters
func isASCIIOnly(s string) bool {
	for _, r := range s {
		if r > 127 {
			return false
		}
	}
	return true
}

// normalizeFullASCII normalizes the full ASCII name for comparison
func normalizeFullASCII(s string) string {
	// Remove extra spaces and normalize to uppercase for comparison
	normalized := strings.Fields(strings.ToUpper(s))
	return strings.Join(normalized, " ")
}

// TestSpecificTransliterationCases tests specific character mappings
func TestSpecificTransliterationCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string // Output should contain these substrings
		notContains []string // Output should NOT contain these substrings
	}{
		{
			name:  "German eszett to ss",
			input: "Weiß",
			contains: []string{"ss"},
			notContains: []string{"ß", "?"},
		},
		{
			name:  "German umlauts",
			input: "Müller Jürgen Größe",
			contains: []string{"Muller", "Jurgen", "Grosse"},
			notContains: []string{"ü", "ö", "ä", "?"},
		},
		{
			name:  "Vietnamese diacritics",
			input: "Nguyễn",
			contains: []string{"Nguyen"},
			notContains: []string{"ễ", "?"},
		},
		{
			name:  "Chinese characters",
			input: "李小龍",
			contains: []string{"Li"}, // At minimum should get family name
			notContains: []string{"龍", "?"},
		},
		{
			name:  "Spanish accents and ñ",
			input: "Núñez José María",
			contains: []string{"Nunez", "Jose", "Maria"},
			notContains: []string{"ñ", "é", "í", "?"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := callTransliterateAPI(tt.input)
			if err != nil {
				t.Fatalf("API call failed: %v", err)
			}

			// Check contains
			for _, substr := range tt.contains {
				if !strings.Contains(response.OutputText, substr) {
					t.Errorf("Output %q should contain %q", response.OutputText, substr)
				}
			}

			// Check not contains
			for _, substr := range tt.notContains {
				if strings.Contains(response.OutputText, substr) {
					t.Errorf("Output %q should NOT contain %q", response.OutputText, substr)
				}
			}

			// Ensure it's ASCII only
			if !isASCIIOnly(response.OutputText) {
				t.Errorf("Output should be ASCII only, got: %q", response.OutputText)
			}
		})
	}
}

// TestJapaneseTransliteration tests Japanese to ASCII conversion
func TestJapaneseTransliteration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		notContains []string
	}{
		{
			name:  "Japanese hiragana",
			input: "たなか",
			notContains: []string{"た", "な", "か", "?"},
		},
		{
			name:  "Japanese katakana",
			input: "タナカ",
			notContains: []string{"タ", "ナ", "カ", "?"},
		},
		{
			name:  "Japanese kanji",
			input: "田中",
			notContains: []string{"田", "中", "?"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := callTransliterateAPI(tt.input)
			if err != nil {
				t.Fatalf("API call failed: %v", err)
			}

			// Check not contains
			for _, substr := range tt.notContains {
				if strings.Contains(response.OutputText, substr) {
					t.Errorf("Output %q should NOT contain %q", response.OutputText, substr)
				}
			}

			// Ensure it's ASCII only
			if !isASCIIOnly(response.OutputText) {
				t.Errorf("Output should be ASCII only, got: %q", response.OutputText)
			}

			// Ensure no question marks
			if strings.Contains(response.OutputText, "?") {
				t.Errorf("Output should not contain '?', got: %q", response.OutputText)
			}
		})
	}
}

func main() {
	t := &testing.T{}
	fmt.Println("Running validation tests manually...")
	
	// Run the main validation test
	TestValidationCases(t)
	
	// Run specific transliteration tests  
	TestSpecificTransliterationCases(t)
	
	// Run Japanese tests
	TestJapaneseTransliteration(t)
	
	fmt.Println("Validation tests completed.")
}

