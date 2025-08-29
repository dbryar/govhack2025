-- Transliteration database schema
-- Tracks input/output pairs with confidence levels and locale information

-- Primary transliteration results table
CREATE TABLE transliterations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    input_text TEXT NOT NULL,
    output_text TEXT NOT NULL,
    input_script VARCHAR(50) NOT NULL,  -- e.g., 'cyrillic', 'chinese', 'arabic'
    output_script VARCHAR(50) NOT NULL, -- e.g., 'latin', 'ascii'
    input_locale VARCHAR(10),           -- ISO language code e.g., 'zh-CN', 'ru-RU'
    confidence_score DECIMAL(3,2),      -- 0.00-1.00 confidence level
    usage_count INTEGER DEFAULT 1,     -- Track frequency of this transliteration
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Character mapping reference table for building transliteration rules
CREATE TABLE character_mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_char VARCHAR(10) NOT NULL,   -- Original character (supports multi-byte)
    target_char VARCHAR(50) NOT NULL,   -- Target transliteration
    source_script VARCHAR(50) NOT NULL,
    target_script VARCHAR(50) NOT NULL,
    locale VARCHAR(10),                 -- Optional locale specificity
    frequency_weight DECIMAL(3,2) DEFAULT 0.50, -- How common this mapping is
    context_rules TEXT,                 -- JSON rules for when to apply this mapping
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- User feedback and corrections table
CREATE TABLE transliteration_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transliteration_id UUID REFERENCES transliterations(id),
    suggested_output TEXT NOT NULL,
    feedback_type VARCHAR(20) NOT NULL, -- 'correction', 'alternative', 'preferred'
    user_context TEXT,                  -- Optional context about the correction
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_transliterations_input ON transliterations(input_text);
CREATE INDEX idx_transliterations_scripts ON transliterations(input_script, output_script);
CREATE INDEX idx_transliterations_locale ON transliterations(input_locale);
CREATE INDEX idx_character_mappings_source ON character_mappings(source_char, source_script);
CREATE INDEX idx_character_mappings_locale ON character_mappings(locale, source_script);