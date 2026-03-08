package storage

import (
	"sync"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sqlx.DB
	mu sync.RWMutex
}

func NewDB(file string) (*Storage, error) {
	db, err := sqlx.Open("sqlite", file)
	if err != nil {
		return nil, err
	}

	return &Storage{db, sync.RWMutex{}}, nil
}

func (st *Storage) InitSchemas(schemas ...string) error {
	for _, schema := range schemas {
		if _, err := st.db.Exec(schema); err != nil {
			return err
		}
	}
	return nil
}

func (st *Storage) Close() error {
	return st.db.Close()
}
