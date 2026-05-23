-- Force all grammar seeds to re-run so somatome patterns get re-indexed
-- (previously 100-147, now 21-68 in the seed file)
DELETE FROM schema_seeds WHERE name LIKE '%grammar%';

-- Re-sequence any existing index_position gaps within each level
-- so GetByLevelAndIndex(level, n) always finds a pattern for n < COUNT(*)
WITH ranked AS (
  SELECT id,
         (ROW_NUMBER() OVER (PARTITION BY jlpt_level ORDER BY index_position ASC, created_at ASC) - 1) AS new_pos
  FROM grammar_patterns
)
UPDATE grammar_patterns
SET index_position = (SELECT new_pos FROM ranked WHERE ranked.id = grammar_patterns.id);

-- Reset all user grammar indices to 0 to avoid out-of-range positions
UPDATE user_progress SET current_grammar_index = 0;
