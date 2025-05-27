ALTER TABLE votes
    ADD COLUMN IF NOT EXISTS error_message TEXT,
    ADD COLUMN IF NOT EXISTS retry_count INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS processed_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_votes_status ON votes (status);

CREATE INDEX IF NOT EXISTS idx_votes_processed_at ON votes (processed_at);