-- Clear JLPT seed tracking so seeds re-run with the corrected insertJLPTQuestion
-- which no longer double-marshals the options JSON array.
DELETE FROM schema_seeds WHERE name LIKE '%mock_test%' OR name LIKE '%jlpt%';
