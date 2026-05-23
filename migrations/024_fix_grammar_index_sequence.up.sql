-- Force all grammar seeds to re-run so somatome patterns get re-indexed
DELETE FROM schema_seeds WHERE name LIKE '%grammar%';

-- Re-sequence index_position gaps within each level using window function
WITH ranked AS (
  SELECT id,
         (ROW_NUMBER() OVER (PARTITION BY jlpt_level ORDER BY index_position ASC, created_at ASC) - 1) AS new_pos
  FROM grammar_patterns
)
UPDATE grammar_patterns
SET index_position = ranked.new_pos
FROM ranked
WHERE grammar_patterns.id = ranked.id;

-- Reset all user grammar indices to avoid out-of-range positions
UPDATE user_progress SET current_grammar_index = 0;
