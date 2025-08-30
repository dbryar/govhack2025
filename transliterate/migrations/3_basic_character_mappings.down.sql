-- Remove basic character mappings
DELETE FROM character_mappings WHERE source_script IN ('cyrillic', 'latin', 'chinese');