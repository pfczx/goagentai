package memory

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pfczx/goagentai/memory/db/generated"
	"github.com/pressly/goose/v3"
)

type MemoryHandler struct {
	db      *sql.DB
	queries *generated.Queries
}

func (m *MemoryHandler) CloseDB() error {
	if m.db == nil {
		return fmt.Errorf("Unexpected nil value when closing DB")
	}
	return m.db.Close()
}

func InitDatabase() (*sql.DB, *generated.Queries, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, err
	}
	dbPath := filepath.Join(path, ".config", "goagent", "long_term.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot open database: %w", err)
	}
	//supress goose output
	silentLogger := log.New(io.Discard, "", 0)
	goose.SetLogger(silentLogger)
	db.SetMaxOpenConns(1)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, nil, err
	}
	migrationsPath := "memory/db/schema"
	if err := goose.Up(db, migrationsPath); err != nil {
		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, generated.New(db), nil

}

func NewMemoryHandler() (*MemoryHandler, error) {
	database, q, err := InitDatabase()
	if err != nil {
		return nil, err
	}
	return &MemoryHandler{
		queries: q,
		db:      database,
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

func (m *MemoryHandler) ClearShortTerm(profileID string) error {
	ctx := context.Background()
	return m.queries.ClearShortTermMemoryByProfile(ctx, profileID)

}

func (m *MemoryHandler) TrimLongTermToStorageSize(profileID string, size int) error {
	ctx := context.Background()
	count, err := m.queries.LongTermMemorySizeForProfile(ctx, profileID)
	if err != nil {
		return err
	}

	limit := count - int64(size)
	if limit <= 0 {

		return nil
	}

	params := generated.DeleteOldLongTermForProfileParams{
		ProfileID: profileID,
		Limit:     limit,
	}

	return m.queries.DeleteOldLongTermForProfile(ctx, params)

}

func (m *MemoryHandler) SaveLongTerm(profileID, content string, tf map[string]float32) error {
	ctx := context.Background()

	tfJSON, err := json.Marshal(tf)
	if err != nil {
		return fmt.Errorf("cannot marshal tfidf: %w", err)
	}

	arg := generated.InsertLongTermParams{
		ProfileID: profileID,
		Content:   content,
		Tf:        sql.NullString{String: string(tfJSON), Valid: true},
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
		out = append(out, MemoryChunk{
			Profile:   profileID,
			Summary:   memo.Content,
			TF:        tf,
			CreatedAt: memo.CreatedAt.Time,
		})
	}
	return out, nil

}

func (m *MemoryHandler) ClearLongTerm(profileID string) error {
	ctx := context.Background()
	return m.queries.ClearLongMemory(ctx, profileID)
}
