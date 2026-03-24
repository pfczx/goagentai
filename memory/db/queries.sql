-- name: InsertShortTerm :exec
INSERT INTO short_term_memory (profile_id, memory)
VALUES (?, ?);


-- name: GetShortTermByProfile :many
SELECT id, profile_id, memory, created_at
FROM short_term_memory
WHERE profile_id = ?;


-- name: InsertLongTerm :exec
INSERT INTO long_term_memory (
    profile_id,
    content,
    tf
) VALUES (?, ?, ?);


-- name: GetLongTermByProfile :many
SELECT id, profile_id, content, tf, created_at
FROM long_term_memory
WHERE profile_id = ?;


-- name: CountShortTermByProfile :one
SELECT COUNT(*) 
FROM short_term_memory
WHERE profile_id = ?;


-- name: ClearShortTermMemoryByProfile :exec
DELETE FROM short_term_memory
WHERE profile_id = ?;


-- name: LongTermMemorySizeForProfile :one
SELECT COUNT(*)
FROM long_term_memory
WHERE profile_id = ?;


-- name: DeleteOldLongTermForProfile :exec 
DELETE FROM long_term_memory AS ltm
WHERE ltm.id IN (
  SELECT sub.id 
  FROM long_term_memory AS sub 
  WHERE sub.profile_id = ? 
  ORDER BY sub.created_at ASC 
  LIMIT ?
);
