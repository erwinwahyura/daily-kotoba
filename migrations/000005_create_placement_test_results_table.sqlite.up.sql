CREATE TABLE placement_test_results (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    test_score INTEGER NOT NULL,
    assigned_level TEXT NOT NULL,
    completed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_placement_user ON placement_test_results(user_id);