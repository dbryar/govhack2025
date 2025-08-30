package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

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

const baseURL = "http://localhost:4000"

func main() {
	// Wait a moment for the API to be ready
	time.Sleep(2 * time.Second)
	
	fmt.Println("Running API validation tests...")
	
	testCases := []struct {
		name  string
		input string
	}{
		{"Vietnamese male with title", "Doctor Nguyễn Văn Minh"},
		{"German with umlaut and eszett", "Prof. Jürgen Groß"},
		{"Chinese name", "李小龍"},
		{"Japanese honorific", "Tanaka-san Yoko"},
		{"Spanish with particle", "Maria del Carmen Núñez"},
	}
	
	totalTests := 0
	passedTests := 0
	
	for _, tc := range testCases {
		fmt.Printf("\n--- Testing: %s ---\n", tc.name)
		fmt.Printf("Input: %s\n", tc.input)
		
		response, err := callAPI(tc.input)
		if err != nil {
			fmt.Printf("❌ API call failed: %v\n", err)
			totalTests++
			continue
		}
		
		totalTests++
		issues := validateResponse(response)
		
		fmt.Printf("Output Text: %s\n", response.OutputText)
		if response.Name != nil {
			fmt.Printf("Family: %s\n", response.Name.Family)
			fmt.Printf("First: %s\n", response.Name.First)
			fmt.Printf("Middle: %v\n", response.Name.Middle)
			fmt.Printf("Titles: %v\n", response.Name.Titles)
			fmt.Printf("FullASCII: %s\n", response.Name.FullASCII)
		}
		if response.Gender != nil {
			fmt.Printf("Gender: %s (%.2f confidence)\n", response.Gender.Value, response.Gender.Confidence)
		}
		
		if len(issues) == 0 {
			fmt.Printf("✅ Test passed\n")
			passedTests++
		} else {
			fmt.Printf("❌ Test failed with issues:\n")
			for _, issue := range issues {
				fmt.Printf("  - %s\n", issue)
			}
		}
	}
	
	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Passed: %d/%d tests\n", passedTests, totalTests)
	
	if passedTests == totalTests {
		fmt.Println("🎉 All tests passed!")
	} else {
		fmt.Printf("💡 %d tests need fixing\n", totalTests-passedTests)
	}
}

func callAPI(input string) (*APIResponse, error) {
	payload := map[string]interface{}{
		"text":          input,
		"output_script": "ascii",
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(baseURL+"/transliterate", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, body)
	}

	var apiResponse APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiResponse, nil
}

func validateResponse(response *APIResponse) []string {
	var issues []string
	
	// Check for Unicode characters (should be ASCII only)
	if !isASCIIOnly(response.OutputText) {
		issues = append(issues, fmt.Sprintf("Output contains non-ASCII characters: %q", response.OutputText))
	}
	
	// Check for question marks (should not contain ?)
	if strings.Contains(response.OutputText, "?") {
		issues = append(issues, fmt.Sprintf("Output contains question marks: %q", response.OutputText))
	}
	
	if response.Name != nil {
		// Check FullASCII
		if !isASCIIOnly(response.Name.FullASCII) {
			issues = append(issues, fmt.Sprintf("FullASCII contains non-ASCII: %q", response.Name.FullASCII))
		}
		if strings.Contains(response.Name.FullASCII, "?") {
			issues = append(issues, fmt.Sprintf("FullASCII contains question marks: %q", response.Name.FullASCII))
		}
		
		// Check that titles are not in name fields
		titlePatterns := []string{"DR", "DOCTOR", "PROF", "PROFESSOR", "MR", "MRS", "MS", "MISS"}
		for _, pattern := range titlePatterns {
			if strings.Contains(strings.ToUpper(response.Name.First), pattern) {
				issues = append(issues, fmt.Sprintf("Title %q found in first name: %q", pattern, response.Name.First))
			}
			if strings.Contains(strings.ToUpper(response.Name.Family), pattern) {
				issues = append(issues, fmt.Sprintf("Title %q found in family name: %q", pattern, response.Name.Family))
			}
			for i, middle := range response.Name.Middle {
				if strings.Contains(strings.ToUpper(middle), pattern) {
					issues = append(issues, fmt.Sprintf("Title %q found in middle[%d]: %q", pattern, i, middle))
				}
			}
		}
		
		// Check for non-ASCII in name components
		if !isASCIIOnly(response.Name.Family) {
			issues = append(issues, fmt.Sprintf("Family name contains non-ASCII: %q", response.Name.Family))
		}
		if !isASCIIOnly(response.Name.First) {
			issues = append(issues, fmt.Sprintf("First name contains non-ASCII: %q", response.Name.First))
		}
		for _, middle := range response.Name.Middle {
			if !isASCIIOnly(middle) {
				issues = append(issues, fmt.Sprintf("Middle name contains non-ASCII: %q", middle))
			}
		}
	} else {
		issues = append(issues, "Name structure is missing")
	}
	
	return issues
}

func isASCIIOnly(s string) bool {
	for _, r := range s {
		if r > 127 {
			return false
		}
	}
	return true
}