package sqlite

import (
	"BotAuc/lib/e"
	"BotAuc/storage"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
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

func (s *Storage) SaveData(ctx context.Context, auc *storage.Auction) error {
	q := `DELETE FROM aucs WHERE URL = ?;
			INSERT INTO aucs (Name, URL, StartDate, EndDate, Status) VALUES (?, ?, ?, ?, ?)`

	if _, err := s.db.ExecContext(ctx, q, auc.URL, auc.Name, auc.URL, auc.StartDate, auc.EndDate, auc.Status); err != nil {
		return e.Wrap("can't save auc", err)
	}
	return nil
}

func (s *Storage) IsExists(ctx context.Context, auc *storage.Auction) (bool, error) {
	q := `SELECT COUNT(*) FROM aucs WHERE URL = ?`

	var count int
	if err := s.db.QueryRowContext(ctx, q, auc.URL).Scan(&count); err != nil {
		return false, e.Wrap("can't check auc", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS aucs (Name TEXT, URL TEXT, StartDate DATETIME, EndDate DATETIME, Status TEXT)`
	if _, err := s.db.ExecContext(ctx, q); err != nil {
		err = e.Wrap("can't create table auc", err)
		return err
	}
	return nil
}

func (s *Storage) ActualizeAucs(ctx context.Context, urls *storage.UrlsAlias) error {
	q := `DELETE FROM aucs WHERE NOT URL IN (%s)`
	q = fmt.Sprintf(q, listOfURLParams(urls))

	params := make([]interface{}, 0)
	for _, url := range *urls {
		params = append(params, url)
	}

	if _, err := s.db.ExecContext(ctx, q, params...); err != nil {
		return e.Wrap("can't save auc", err)
	}
	return nil
}

func (s *Storage) GetFutureAucs(ctx context.Context) (string, error) {
	return "", nil
}

func listOfURLParams(urls *storage.UrlsAlias) string {
	res := strings.Repeat("?, ", len(*urls)-1) + "?"
	return res
}
