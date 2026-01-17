CREATE TABLE placement_test_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    test_score INTEGER NOT NULL,
    assigned_level VARCHAR(10) NOT NULL,
    completed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_placement_user ON placement_test_results(user_id);
