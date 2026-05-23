-- Remove seed tracking for JLPT mock tests so they re-run with the ON CONFLICT upsert,
-- which will populate the options column for questions that were seeded without it.
DELETE FROM schema_seeds WHERE name LIKE '%mock_test%' OR name LIKE '%jlpt%';
