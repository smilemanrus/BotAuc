package sqlite

import (
	"BotAuc/lib/e"
	"database/sql"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, e.Wrap("Can't open db", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap("Can't connect to db", err)
	}
	return &Storage{db: db}, nil
}
