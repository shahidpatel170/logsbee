-- db/schema.sql
CREATE TABLE IF NOT EXISTS hives (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    hostname VARCHAR(255) NOT NULL,
    init_system VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS harvested_logs (
    id BIGSERIAL PRIMARY KEY,
    hive_id INT REFERENCES hives(id) ON DELETE CASCADE NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    severity VARCHAR(50) NOT NULL,
    systemd_unit VARCHAR(255),
    component VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    raw_payload JSONB
);

-- index
CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON harvested_logs(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_logs_severity ON harvested_logs(severity);
CREATE INDEX IF NOT EXISTS idx_logs_search_vectors ON harvested_logs USING gin(to_tsvector('english', message));
