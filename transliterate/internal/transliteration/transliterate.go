// Package transliteration provides comprehensive character-to-character mapping
// and conversion between different writing systems.
package transliteration

import (
	"context"
	"database/sql"
	"strings"
	"unicode/utf8"
	"errors"

	"encore.dev/storage/sqldb"
)

// Config holds transliteration configuration
type Config struct {
	UseDatabase    bool
	FallbackToASCII bool
	PreserveSpacing bool
	CaseSensitive  bool
}

// DefaultConfig returns sensible defaults
func DefaultConfig() Config {
	return Config{
		UseDatabase:    true,
		FallbackToASCII: true,
		PreserveSpacing: true,
		CaseSensitive:  false,
	}
}

// Result represents the result of a transliteration
type Result struct {
	Output     string
	Confidence float64
	Notes      []string
	Method     string // "database", "builtin", "fallback"
}

// Engine handles transliteration operations
type Engine struct {
	config Config
	db     *sqldb.Database
}

// NewEngine creates a new transliteration engine
func NewEngine(config Config, db *sqldb.Database) *Engine {
	return &Engine{
		config: config,
		db:     db,
	}
}

// Transliterate converts text from one script to another
func (e *Engine) Transliterate(ctx context.Context, text, fromScript, toScript, locale string) (*Result, error) {
	if !utf8.ValidString(text) {
		return nil, ErrInvalidUTF8
	}

	if text == "" {
		return &Result{Output: "", Confidence: 1.0, Method: "empty"}, nil
	}

	var result strings.Builder
	var notes []string
	var confidenceSum float64
	var charCount int

	// Process character by character
	for _, r := range text {
		charResult, err := e.transliterateRune(ctx, r, fromScript, toScript, locale)
		if err != nil {
			return nil, err
		}

		result.WriteString(charResult.Output)
		if charResult.Note != "" {
			notes = append(notes, charResult.Note)
		}
		confidenceSum += charResult.Confidence
		charCount++
	}

	// Calculate average confidence
	confidence := confidenceSum / float64(charCount)
	if charCount == 0 {
		confidence = 1.0
	}

	// Determine primary method used
	method := "mixed"
	if len(notes) == 0 {
		method = "builtin"
	}

	return &Result{
		Output:     result.String(),
		Confidence: confidence,
		Notes:      notes,
		Method:     method,
	}, nil
}

// RuneResult represents the result of transliterating a single rune
type RuneResult struct {
	Output     string
	Confidence float64
	Note       string
	Method     string
}

// transliterateRune converts a single rune
func (e *Engine) transliterateRune(ctx context.Context, r rune, fromScript, toScript, locale string) (*RuneResult, error) {
	sourceChar := string(r)

	// Try database lookup first
	if e.config.UseDatabase {
		if dbResult, err := e.lookupInDatabase(ctx, sourceChar, fromScript, toScript, locale); err == nil && dbResult != "" {
			return &RuneResult{
				Output:     dbResult,
				Confidence: 0.95,
				Method:     "database",
			}, nil
		}
	}

	// Try built-in rules
	if builtinResult := e.applyBuiltinRules(r, fromScript, toScript); builtinResult != "" {
		return &RuneResult{
			Output:     builtinResult,
			Confidence: 0.85,
			Method:     "builtin",
		}, nil
	}

	// Fallback to ASCII approximation
	if e.config.FallbackToASCII && toScript == "ascii" {
		asciiResult := e.approximateToASCII(r)
		confidence := 0.3
		note := ""
		if asciiResult == "?" {
			note = "Unknown character approximated"
			confidence = 0.1
		}
		return &RuneResult{
			Output:     asciiResult,
			Confidence: confidence,
			Note:       note,
			Method:     "fallback",
		}, nil
	}

	// Keep original character as last resort
	return &RuneResult{
		Output:     sourceChar,
		Confidence: 0.1,
		Note:       "Character unchanged",
		Method:     "unchanged",
	}, nil
}

// lookupInDatabase performs database lookup for character mapping
func (e *Engine) lookupInDatabase(ctx context.Context, sourceChar, fromScript, toScript, locale string) (string, error) {
	var targetChar string
	
	err := e.db.QueryRow(ctx, `
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
	`, sourceChar, fromScript, toScript, locale).Scan(&targetChar)

	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return targetChar, nil
}

// applyBuiltinRules applies hardcoded transliteration rules
func (e *Engine) applyBuiltinRules(r rune, fromScript, toScript string) string {
	switch fromScript {
	case "cyrillic":
		if toScript == "latin" || toScript == "ascii" {
			return e.transliterateCyrillic(r)
		}
	case "chinese":
		if toScript == "latin" || toScript == "ascii" {
			return e.transliterateChinese(r)
		}
	case "japanese":
		if toScript == "latin" || toScript == "ascii" {
			return e.transliterateJapanese(r)
		}
	case "arabic":
		if toScript == "latin" || toScript == "ascii" {
			return e.transliterateArabic(r)
		}
	case "greek":
		if toScript == "latin" || toScript == "ascii" {
			return e.transliterateGreek(r)
		}
	case "korean":
		if toScript == "latin" || toScript == "ascii" {
			return e.transliterateKorean(r)
		}
	case "hebrew":
		if toScript == "latin" || toScript == "ascii" {
			return e.transliterateHebrew(r)
		}
	case "thai":
		if toScript == "latin" || toScript == "ascii" {
			return e.transliterateThai(r)
		}
	}
	return ""
}

// transliterateCyrillic handles Cyrillic to Latin conversion
func (e *Engine) transliterateCyrillic(r rune) string {
	mapping := map[rune]string{
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
	
	return mapping[r]
}

// transliterateChinese handles Chinese to Latin conversion
func (e *Engine) transliterateChinese(r rune) string {
	// Comprehensive mapping for common Chinese characters
	mapping := map[rune]string{
		// Numbers
		'一': "Yi", '二': "Er", '三': "San", '四': "Si", '五': "Wu",
		'六': "Liu", '七': "Qi", '八': "Ba", '九': "Jiu", '十': "Shi",
		
		// Common surnames
		'李': "Li", '王': "Wang", '张': "Zhang", '刘': "Liu", '陈': "Chen",
		'杨': "Yang", '赵': "Zhao", '黄': "Huang", '周': "Zhou", '吴': "Wu",
		'徐': "Xu", '孙': "Sun", '胡': "Hu", '朱': "Zhu", '高': "Gao",
		'林': "Lin", '何': "He", '郭': "Guo", '马': "Ma", '罗': "Luo",
		'梁': "Liang", '宋': "Song", '郑': "Zheng", '谢': "Xie", '韩': "Han",
		'唐': "Tang", '冯': "Feng", '于': "Yu", '董': "Dong", '萧': "Xiao",
		'程': "Cheng", '曹': "Cao", '袁': "Yuan", '邓': "Deng", '许': "Xu",
		'傅': "Fu", '沈': "Shen", '曾': "Zeng", '彭': "Peng", '吕': "Lu",
		
		// Common given names
		'小': "Xiao", '大': "Da", '中': "Zhong", '文': "Wen", '明': "Ming",
		'华': "Hua", '建': "Jian", '国': "Guo", '民': "Min", '伟': "Wei",
		'龍': "Long", '龙': "Long", '凤': "Feng", '鳳': "Feng", '玉': "Yu",
		'金': "Jin", '春': "Chun", '红': "Hong", '军': "Jun", '强': "Qiang",
		'云': "Yun", '平': "Ping", '志': "Zhi", '刚': "Gang", '勇': "Yong",
		'磊': "Lei", '娜': "Na", '静': "Jing", '丽': "Li", '敏': "Min",
		'秀': "Xiu", '英': "Ying", '芳': "Fang", '燕': "Yan", '雪': "Xue",
		'琴': "Qin", '梅': "Mei", '莉': "Li", '兰': "Lan", '翠': "Cui",
		
		// Common words
		'你': "ni", '好': "hao", '是': "shi", '的': "de", '我': "wo",
		'他': "ta", '她': "ta", '们': "men", '有': "you", '在': "zai",
		'了': "le", '不': "bu", '就': "jiu", '人': "ren", '都': "dou",
		
		// Directions
		'东': "Dong", '南': "Nan", '西': "Xi", '北': "Bei",
		'上': "Shang", '下': "Xia", '左': "Zuo", '右': "You",
		'前': "Qian", '后': "Hou",
		
		// Time/descriptors
		'新': "Xin", '老': "Lao", '长': "Chang", '短': "Duan",
		'低': "Di", '快': "Kuai", '慢': "Man",
		'早': "Zao", '晚': "Wan",
	}
	
	return mapping[r]
}

// transliterateArabic handles Arabic to Latin conversion
func (e *Engine) transliterateArabic(r rune) string {
	mapping := map[rune]string{
		'ا': "a", 'ب': "b", 'ت': "t", 'ث': "th", 'ج': "j", 'ح': "h",
		'خ': "kh", 'د': "d", 'ذ': "dh", 'ر': "r", 'ز': "z", 'س': "s",
		'ش': "sh", 'ص': "s", 'ض': "d", 'ط': "t", 'ظ': "z", 'ع': "'",
		'غ': "gh", 'ف': "f", 'ق': "q", 'ك': "k", 'ل': "l", 'م': "m",
		'ن': "n", 'ه': "h", 'و': "w", 'ي': "y",
		
		// Additional Arabic letters
		'ء': "'", 'آ': "aa", 'أ': "a", 'إ': "i", 'ؤ': "u", 'ئ': "i",
		'ة': "h", 'ى': "a",
	}
	
	return mapping[r]
}

// transliterateGreek handles Greek to Latin conversion
func (e *Engine) transliterateGreek(r rune) string {
	mapping := map[rune]string{
		// Uppercase
		'Α': "A", 'Β': "B", 'Γ': "G", 'Δ': "D", 'Ε': "E", 'Ζ': "Z",
		'Η': "H", 'Θ': "Th", 'Ι': "I", 'Κ': "K", 'Λ': "L", 'Μ': "M",
		'Ν': "N", 'Ξ': "X", 'Ο': "O", 'Π': "P", 'Ρ': "R", 'Σ': "S",
		'Τ': "T", 'Υ': "Y", 'Φ': "Ph", 'Χ': "Ch", 'Ψ': "Ps", 'Ω': "O",
		
		// Lowercase
		'α': "a", 'β': "b", 'γ': "g", 'δ': "d", 'ε': "e", 'ζ': "z",
		'η': "h", 'θ': "th", 'ι': "i", 'κ': "k", 'λ': "l", 'μ': "m",
		'ν': "n", 'ξ': "x", 'ο': "o", 'π': "p", 'ρ': "r", 'σ': "s", 'ς': "s",
		'τ': "t", 'υ': "y", 'φ': "ph", 'χ': "ch", 'ψ': "ps", 'ω': "o",
	}
	
	return mapping[r]
}

// transliterateJapanese handles Japanese to Latin conversion (basic)
func (e *Engine) transliterateJapanese(r rune) string {
	// Basic Hiragana mappings
	hiragana := map[rune]string{
		'あ': "a", 'い': "i", 'う': "u", 'え': "e", 'お': "o",
		'か': "ka", 'き': "ki", 'く': "ku", 'け': "ke", 'こ': "ko",
		'が': "ga", 'ぎ': "gi", 'ぐ': "gu", 'げ': "ge", 'ご': "go",
		'さ': "sa", 'し': "shi", 'す': "su", 'せ': "se", 'そ': "so",
		'ざ': "za", 'じ': "ji", 'ず': "zu", 'ぜ': "ze", 'ぞ': "zo",
		'た': "ta", 'ち': "chi", 'つ': "tsu", 'て': "te", 'と': "to",
		'だ': "da", 'ぢ': "ji", 'づ': "zu", 'で': "de", 'ど': "do",
		'な': "na", 'に': "ni", 'ぬ': "nu", 'ね': "ne", 'の': "no",
		'は': "ha", 'ひ': "hi", 'ふ': "fu", 'へ': "he", 'ほ': "ho",
		'ば': "ba", 'び': "bi", 'ぶ': "bu", 'べ': "be", 'ぼ': "bo",
		'ぱ': "pa", 'ぴ': "pi", 'ぷ': "pu", 'ぺ': "pe", 'ぽ': "po",
		'ま': "ma", 'み': "mi", 'む': "mu", 'め': "me", 'も': "mo",
		'や': "ya", 'ゆ': "yu", 'よ': "yo",
		'ら': "ra", 'り': "ri", 'る': "ru", 'れ': "re", 'ろ': "ro",
		'わ': "wa", 'を': "wo", 'ん': "n",
	}
	
	return hiragana[r]
}

// transliterateKorean handles Korean to Latin conversion (basic)
func (e *Engine) transliterateKorean(r rune) string {
	// This is a simplified approach - full Korean transliteration is complex
	// Would need proper Hangul decomposition
	return ""
}

// transliterateHebrew handles Hebrew to Latin conversion
func (e *Engine) transliterateHebrew(r rune) string {
	mapping := map[rune]string{
		'א': "'", 'ב': "b", 'ג': "g", 'ד': "d", 'ה': "h", 'ו': "v",
		'ז': "z", 'ח': "ch", 'ט': "t", 'י': "y", 'כ': "kh", 'ל': "l",
		'מ': "m", 'נ': "n", 'ס': "s", 'ע': "'", 'פ': "p", 'צ': "ts",
		'ק': "q", 'ר': "r", 'ש': "sh", 'ת': "t",
	}
	
	return mapping[r]
}

// transliterateThai handles Thai to Latin conversion
func (e *Engine) transliterateThai(r rune) string {
	// Basic Thai consonants
	mapping := map[rune]string{
		'ก': "k", 'ข': "kh", 'ค': "kh", 'ง': "ng", 'จ': "j", 'ฉ': "ch",
		'ช': "ch", 'ซ': "s", 'ญ': "y", 'ด': "d", 'ต': "t", 'ถ': "th",
		'ท': "th", 'น': "n", 'บ': "b", 'ป': "p", 'ผ': "ph", 'ฝ': "f",
		'พ': "ph", 'ฟ': "f", 'ภ': "ph", 'ม': "m", 'ย': "y", 'ร': "r",
		'ล': "l", 'ว': "w", 'ศ': "s", 'ษ': "s", 'ส': "s", 'ห': "h",
		'อ': "'", 'ฮ': "h",
		
		// Vowels
		'า': "a", 'ิ': "i", 'ี': "i", 'ึ': "ue", 'ื': "ue", 'ุ': "u", 'ู': "u",
		'เ': "e", 'แ': "ae", 'โ': "o", 'ใ': "ai", 'ไ': "ai",
	}
	
	return mapping[r]
}

// approximateToASCII provides fallback ASCII approximation
func (e *Engine) approximateToASCII(r rune) string {
	// Handle already-ASCII characters
	if r < 128 {
		return string(r)
	}

	// Use our Unicode normalization for ASCII conversion
	// This is a simplified version - would integrate with unicode package
	
	// Common approximations
	approximations := map[rune]string{
		// Accented vowels -> base vowels
		'á': "a", 'à': "a", 'â': "a", 'ã': "a", 'ä': "a", 'å': "a",
		'é': "e", 'è': "e", 'ê': "e", 'ë': "e",
		'í': "i", 'ì': "i", 'î': "i", 'ï': "i",
		'ó': "o", 'ò': "o", 'ô': "o", 'õ': "o", 'ö': "o", 'ø': "o",
		'ú': "u", 'ù': "u", 'û': "u", 'ü': "u",
		'ý': "y", 'ỳ': "y", 'ŷ': "y", 'ÿ': "y",
		
		// Other common characters
		'ç': "c", 'ñ': "n", 'ß': "ss",
	}
	
	if approx, exists := approximations[r]; exists {
		return approx
	}
	
	// Default fallback
	return "?"
}

// Custom errors
var (
	ErrInvalidUTF8 = errors.New("invalid UTF-8 input")
)