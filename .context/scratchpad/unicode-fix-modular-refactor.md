# Unicode Fix and Modular Refactor

## Task Summary
Fix Unicode handling errors in transliteration service and improve modular structure while preserving cultural naming conventions and gender markers.

## Key Issues to Address
1. Unicode characters becoming "?" in output
2. Monolithic `transliterate.go` needs better separation of concerns
3. Cultural names and gender markers need better preservation
4. Language/script detection needs improvement
5. Need proper test cases for the five examples from idea.md

## Target Test Cases (from idea.md lines 733-737)
1. `Doctor Nguyễn Văn Minh` → title `Dr.`, family `NGUYEN`, first `Minh`, middle `["Van"]`, gender `M?` (~0.6)
2. `Prof. Jürgen Groß` → title `Prof.`, family `GROSS`, first `Jurgen`, middle `[]`
3. `李小龍` (zh) → likely `Li Xiaolong` (unidecode), family `LI`, given `Xiaolong`
4. `Tanaka-san Yoko` → do **not** map `-san` to title; output note
5. `Maria del Carmen Núñez` → family `NUNEZ` (note: lost tilde), keep `del` particle with family in ASCII

## Required Modules
- Detection: Script and language detection
- Unicode parsing: Proper Unicode normalization and handling
- Transliteration: Character mapping and conversion
- Name parsing: Cultural naming conventions
- Gender inference: Cultural markers and statistical inference

## Status
- [x] Create scratchpad  
- [x] Install Unicode handling packages (golang.org/x/text v0.28.0)
- [x] Refactor into modular structure (5 new modules created)
- [x] Fix Unicode handling issues (proper normalization and ASCII mapping)
- [x] Implement proper test cases (comprehensive tests for 5 key examples)
- [x] Major compilation fixes (duplicate keys, imports, function references)
- [ ] Final compilation and test verification

## Modules Created
1. `internal/detection` - Script and language detection
2. `internal/unicode` - Unicode normalization and ASCII conversion  
3. `internal/transliteration` - Character mapping and conversion
4. `internal/nameparser` - Cultural name parsing
5. `internal/gender` - Gender inference from cultural markers

## Key Improvements
- Proper Unicode normalization using golang.org/x/text
- Language-specific ASCII mappings (German umlauts, Vietnamese diacritics, etc.)
- Cultural context-aware name parsing
- Gender inference with confidence scoring
- Comprehensive test coverage for key use cases

## Remaining Issues
- Minor compilation errors being resolved
- Need to verify test execution