# Name Formatting Feature Implementation

## Task Overview
Implement the core functionality for ASCII name transliteration with proper cultural name handling, as specified in docs/project.md.

## Key Requirements (from docs/project.md)

### Primary Goal
Transform Unicode names to ASCII-compatible structured JSON for legacy systems, addressing real-world problems:
- Vietnamese gender markers (Thị/Văn) being recorded as first names
- Single names (mononyms) from Karen, Kareni, and Chin ethnic groups
- Chinese romanisation causing identity confusion
- Malaysian patronymic structures (binti/bin)

### Expected Output Format
```json
{
  "name": {
    "family": "NGUYEN",
    "first": "Minh", 
    "middle": ["Van"],
    "full_ascii": "DR NGUYEN MINH VAN"
  },
  "gender": { "value": "M", "confidence": 0.65 }
}
```

### Cultural Handling Requirements
- Vietnamese: Proper handling of Thị (female) and Văn (male) markers
- Indonesian/Myanmar: Support for mononyms (single names)
- Chinese: Handle romanisation variants correctly
- Malaysian: Process patronymic binti/bin structures
- Titles: Extract and format titles (Dr, Prof, etc.)

## Implementation Plan
1. Extend API to return structured name output
2. Add cultural name parsing logic
3. Implement gender inference with confidence scoring
4. Add validation test cases for each cultural naming convention
5. Update frontend to display structured output

## Validation Cases Needed
- Vietnamese names with gender markers
- Mononym handling
- Chinese romanisation variations
- Malaysian patronymic structures
- Title extraction
- Mixed cultural names

## Current Status
- Basic transliteration working
- Need to implement structured name parsing
- Need to add gender inference
- Need comprehensive cultural name handling