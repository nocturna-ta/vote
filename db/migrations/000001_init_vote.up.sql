CREATE TABLE IF NOT EXISTS votes (
    id UUID PRIMARY KEY,
    voter_id UUID NOT NULL,
    election_pair_id UUID NOT NULL,
    voted_at TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL,
    transaction_hash VARCHAR(255),
    region VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOL NOT NULL DEFAULT FALSE
);