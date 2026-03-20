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
    tf_idf,
    keywords
) VALUES (?, ?, ?, ?);

-- name: GetLongTermByProfile :many
SELECT id, profile_id, content, tf_idf, keywords, created_at
FROM long_term_memory
WHERE profile_id = ?;

-- name: CountShortTermByProfile :one
SELECT COUNT(*) 
FROM short_term_memory
WHERE profile_id = ?;
