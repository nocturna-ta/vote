-- db/migrations/000003_election_time.up.sql
CREATE TABLE IF NOT EXISTS election_times (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'scheduled',
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,

    -- Constraints
    CONSTRAINT check_election_time_end_after_start CHECK (end_time > start_time),
    CONSTRAINT check_election_status CHECK (status IN ('scheduled', 'active', 'ended'))
    );

-- Index for efficient queries
CREATE INDEX IF NOT EXISTS idx_election_times_status ON election_times (status);
CREATE INDEX IF NOT EXISTS idx_election_times_active ON election_times (is_active);
CREATE INDEX IF NOT EXISTS idx_election_times_time_range ON election_times (start_time, end_time);

-- Ensure only one active election at a time
CREATE UNIQUE INDEX IF NOT EXISTS idx_election_times_unique_active
    ON election_times (is_active)
    WHERE is_active = TRUE AND is_deleted = FALSE;