package memory

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/pfczx/goagentai/memory/db/generated"
	"github.com/pressly/goose/v3"
	"os"
	"path/filepath"
)

type MemoryHandler struct {
	queries *generated.Queries
}

func InitDatabase() (*generated.Queries, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dbPath := filepath.Join(path, ".config", "goagent", "long_term.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open database: %w", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	migrationsPath := "memory/db/schema"
	if err := goose.Up(db, migrationsPath); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return generated.New(db), nil

}

func NewMemoryHandler() (*MemoryHandler, error) {
	q, err := InitDatabase()
	if err != nil {
		return nil, err
	}
	return &MemoryHandler{
		queries: q,
	}, nil
}

func (m *MemoryHandler) SaveShortTerm(profileID, memory string) error {
	ctx := context.Background()

	arg := generated.InsertShortTermParams{
		ProfileID: profileID,
		Memory:    memory,
	}
	return m.queries.InsertShortTerm(ctx, arg)
}

func (m *MemoryHandler) GetShortTerm(profileID string) ([]generated.ShortTermMemory, error) {
	ctx := context.Background()
	return m.queries.GetShortTermByProfile(ctx, profileID)
}

func (m *MemoryHandler) SaveLongTerm(profileID, content string, tfidf map[string]float64, keywords []string) error {
	ctx := context.Background()

	tfidfJSON, err := json.Marshal(tfidf)
	if err != nil {
		return fmt.Errorf("cannot marshal tfidf: %w", err)
	}
	keywordsJSON, err := json.Marshal(keywords)
	if err != nil {
		return fmt.Errorf("cannot marshal keywords: %w", err)
	}

	arg := generated.InsertLongTermParams{
		ProfileID: profileID,
		Content:   content,
		TfIdf:     sql.NullString{String: string(tfidfJSON), Valid: true},
		Keywords:  sql.NullString{String: string(keywordsJSON), Valid: true},
	}

	return m.queries.InsertLongTerm(ctx, arg)
}

func (m *MemoryHandler) CountShortTerm(profileID string) (int, error) {
	ctx := context.Background()
	count, err := m.queries.CountShortTermByProfile(ctx, profileID)
	if err != nil {
		return 1, err
	}
	return int(count), nil
}
