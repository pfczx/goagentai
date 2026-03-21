-- +goose Up
CREATE TABLE short_term_memory (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_id TEXT NOT NULL,
    memory TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE long_term_memory (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_id TEXT NOT NULL,
    content TEXT NOT NULL,
    tf TEXT,
    keywords TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_short_term_profile ON short_term_memory(profile_id);
CREATE INDEX idx_long_term_profile ON long_term_memory(profile_id);

-- +goose Down
DROP TABLE short_term_memory;
DROP TABLE long_term_memory;
