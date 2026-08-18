-- db/queries.sql

-- name: RegisterHive :one
INSERT INTO hives (name, hostname, init_system)
VALUES ($1, $2, $3)
RETURNING id, name, hostname, init_system, created_at;

-- name: GetHiveByName :one
SELECT id, name, hostname, init_system, created_at
FROM hives
WHERE name = $1 LIMIT 1;

-- name: InsertHarvestedLog :exec
INSERT INTO harvested_logs (hive_id, timestamp, severity, systemd_unit, component, message, raw_payload)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: ListLogsPaged :many
SELECT l.id, l.timestamp, l.severity, l.systemd_unit, l.component, l.message, h.name as hive_name
FROM harvested_logs l
JOIN hives h ON l.hive_id = h.id
ORDER BY l.timestamp DESC
LIMIT $1 OFFSET $2;

-- name: SearchLogsByKeyword :many
SELECT l.id, l.timestamp, l.severity, l.systemd_unit, l.component, l.message, h.name as hive_name
FROM harvested_logs l
JOIN hives h ON l.hive_id = h.id
WHERE to_tsvector('english', l.message) @@ to_tsquery('english', $1)
ORDER BY l.timestamp DESC
LIMIT $2;
