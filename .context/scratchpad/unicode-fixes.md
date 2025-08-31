# Unicode Fixes Task

## Current Issues
- Unicode errors producing `?` in name output
- Need better language/script detection
- Cultural names and gender markers not preserved properly
- Frontend examples not working

## Test Cases to Validate (lines 733-737)
1. `Doctor Nguyễn Văn Minh` → title `Dr.`, family `NGUYEN`, first `Minh`, middle `["Van"]`, gender `M?` (~0.6)
2. `Prof. Jürgen Groß` → title `Prof.`, family `GROSS`, first `Jurgen`, middle `[]`
3. `李小龍` (zh) → likely `Li Xiaolong` (unidecode), family `LI`, given `Xiaolong`
4. `Tanaka-san Yoko` → do **not** map `-san` to title; output note
5. `Maria del Carmen Núñez` → family `NUNEZ` (note: lost tilde), keep `del` particle with family

## Frontend Changes Required
- Change text field to single input line
- Enter key should submit the name
- Update examples to working samples

## Technical Focus
- Fix unicode handling with proper Go packages
- Improve transliteration accuracy
- Better cultural name preservation