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

func (m *MemoryHandler) SaveLongTerm(profileID, content string, tf map[string]float32, keywords []string) error {
	ctx := context.Background()

	tfJSON, err := json.Marshal(tf)
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
		Tf:        sql.NullString{String: string(tfJSON), Valid: true},
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

func (m *MemoryHandler) GetLongTerm(profileID string) ([]MemoryChunk, error) {
	ctx := context.Background()
	memory, err := m.queries.GetLongTermByProfile(ctx, profileID)
	if err != nil {
		return nil, err
	}
	var out []MemoryChunk
	for _, memo := range memory {
		var tf map[string]float32
		err := json.Unmarshal([]byte(memo.Tf.String), &tf)
		if err != nil {
			return nil, err
		}
		var keywords []string
		err = json.Unmarshal([]byte(memo.Keywords.String), &keywords)
		if err != nil {
			return nil, err
		}
		out = append(out, MemoryChunk{
			Profile:   profileID,
			Summary:   memo.Content,
			TF:        tf,
			Keywords:  keywords,
			CreatedAt: memo.CreatedAt.Time,
		})
	}
	return out, nil

}
