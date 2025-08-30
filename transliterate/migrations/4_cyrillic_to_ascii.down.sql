-- Remove Cyrillic to ASCII mappings
DELETE FROM character_mappings WHERE source_script = 'cyrillic' AND target_script = 'ascii';