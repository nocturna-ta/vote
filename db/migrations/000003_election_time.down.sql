-- db/migrations/000003_election_time.down.sql
DROP INDEX IF EXISTS idx_election_times_unique_active;
DROP INDEX IF EXISTS idx_election_times_time_range;
DROP INDEX IF EXISTS idx_election_times_active;
DROP INDEX IF EXISTS idx_election_times_status;
DROP TABLE IF EXISTS election_times;