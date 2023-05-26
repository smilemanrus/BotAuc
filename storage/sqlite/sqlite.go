package sqlite

import (
	"BotAuc/lib/e"
	"BotAuc/storage"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"time"
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
	q := `CREATE TABLE IF NOT EXISTS aucs (Name TEXT, URL TEXT, StartDate DATETIME, EndDate DATETIME, Status TEXT);
			CREATE TABLE IF NOT EXISTS subscribes (ID INTEGER, USERNAME TEXT);
			CREATE TABLE IF NOT EXISTS alerts (USERID INTEGER, NOTIFICATIONTYPE TEXT, URL TEXT, EVENTDATE DATETIME)`
	if _, err := s.db.ExecContext(ctx, q); err != nil {
		err = e.Wrap("can't create table auc", err)
		return err
	}
	return nil
}

func (s *Storage) ActualizeAucs(ctx context.Context, urls *storage.UrlsAlias) error {
	q := `DELETE FROM aucs WHERE NOT URL IN (%s)`
	q = fmt.Sprintf(q, listOfParams(len(*urls), "?", ", "))

	params := make([]interface{}, 0)
	for _, url := range *urls {
		params = append(params, url)
	}

	if _, err := s.db.ExecContext(ctx, q, params...); err != nil {
		return e.Wrap("can't save auc", err)
	}
	return nil
}

func (s *Storage) GetFutureAucs(ctx context.Context) (msg string, err error) {
	q := `select Name,
				   URL,
				   CAST(StartDate AS VARCHAR),
				   CAST(EndDate AS VARCHAR)
			from aucs
			WHERE Status = 'ready'`

	row, err := s.db.QueryContext(ctx, q)
	if err != nil {
		err = e.Wrap("can't make request", err)
	}
	defer func() { _ = row.Close() }()

	aucs := make([]string, 0)
	for row.Next() { // Iterate and fetch the records from result cursor
		var name string
		var url string
		var startDate string
		var endDate string
		if err = row.Scan(&name, &url, &startDate, &endDate); err != nil {
			err = e.Wrap("can't scan row", err)
			return "", err
		}

		aucs = append(aucs, fmt.Sprintf(`[%s](%s) с %s по %s`, name, url, startDate, endDate))
	}
	msg = strings.Join(aucs, `\n`)
	return msg, err
}

func (s *Storage) SubscrToAucs(ctx context.Context, chatID int, username string) error {
	q := `INSERT into subscribes(ID, USERNAME)
			SELECT ?, ?
			WHERE NOT EXISTS(select 1 FROM subscribes WHERE ID = ?)`
	if _, err := s.db.ExecContext(ctx, q, chatID, username, chatID); err != nil {
		return e.Wrap("can't exec query to auc subscribing", err)
	}
	return nil
}

func (s *Storage) UnSubscrFormAucs(ctx context.Context, chatID int) error {
	q := `DELETE FROM subscribes WHERE ID = ?`
	if _, err := s.db.ExecContext(ctx, q, chatID); err != nil {
		return e.Wrap("can't exec query to delete subscribing", err)
	}
	return nil
}

func (s *Storage) GetAucsBfrHour(ctx context.Context, eventType string) (storage.EventsData, error) {
	q := `SELECT alerts.USERID,
				   aucs.Name,
				   aucs.URL,
				   aucs.StartDate,
				   aucs.EndDate
			
			FROM alerts
					 INNER JOIN aucs
								ON alerts.URL = aucs.URL
									AND aucs.Status = ?
									AND alerts.NOTIFICATIONTYPE = ?
			ORDER BY alerts.USERID, aucs.StartDate`
	qMakeAlert := strings.ReplaceAll(makeAlertQuery, "'%aucStatus%'", "?")
	qMakeAlert = strings.ReplaceAll(qMakeAlert, "'%notificationType%'", "?")
	aucStatus := "ready"
	if _, err := s.db.ExecContext(ctx, qMakeAlert, aucStatus, eventType); err != nil {
		return nil, e.Wrap("can't exec insert to query for fixing alert", err)
	}

	row, err := s.db.QueryContext(ctx, q, aucStatus, eventType)
	if err != nil {
		return nil, e.Wrap("can't make request", err)
	}
	defer func() { _ = row.Close() }()

	eventsData := make(storage.EventsData, 0)
	for row.Next() { // Iterate and fetch the records from result cursor
		var userID int
		var name string
		var url string
		var startDate string
		var endDate string
		if err = row.Scan(&userID, &name, &url, &startDate, &endDate); err != nil {
			return nil, e.Wrap("can't scan row", err)
		}
		msg := fmt.Sprintf(`[%s](%s) с %s по %s`, name, url, startDate, endDate)
		eventsData[userID] = append(eventsData[userID], storage.NewEventData(msg, url))
	}
	return eventsData, nil
}

func (s *Storage) FixSendingAlert(ctx context.Context, eventsData storage.EventsData, notyType string) error {
	queryBase := trunsctnQuery
	insertQuery := `INSERT INTO alerts (USERID, NOTIFICATIONTYPE, URL, EVENTDATE)
					SELECT DataForFix.*
					;
					`
	deletingQuery := `DELETE FROM alerts
						WHERE USERID IN (%listOfSubs%) AND NOTIFICATIONTYPE = ? AND EVENTDATE IS NULL
						;
						%logic%`
	littleInsTable := mkngDataForFix
	selectsForIns := `SELECT
							?,--USERID
							?,--NOTIFICATIONTYPE
							?,--URL
							? --EVENTDATE`
	unionStr := `
				UNION ALL
	`
	params := make([]interface{}, 0)
	currDate := time.Now()
	//Надо воткнуть список параметров в params согласно запросу deletingQuery в параметр listOfSubs
	for userID, eventData := range eventsData {
		for _, sendedMsg := range eventData {
			params = append(params, userID)
			params = append(params, notyType)
			params = append(params, sendedMsg.URL)
			params = append(params, currDate)
		}
	}
	selectsForIns = listOfParams(len(eventsData), selectsForIns, unionStr)

	subsParams := listOfParams(len(eventsData), "?", ", ")
	deletingQuery = strings.Replace(deletingQuery, "%listOfSubs%", subsParams, 1)

	littleInsTable = strings.Replace(littleInsTable, "%insertTables%", selectsForIns, 1)
	mainLogic := strings.Replace(littleInsTable, "%logic%", deletingQuery, 1)
	mainLogic = strings.Replace(mainLogic, "%logic%", insertQuery, 1)
	queryBase = strings.Replace(queryBase, "%logic%", mainLogic, 1)
	if _, err := s.db.ExecContext(ctx, queryBase, params); err != nil {
		return e.Wrap("can't exec query by fixing alert", err)
	}
	return nil
}

func listOfParams(count int, insString, separator string) string {
	pattern := fmt.Sprintf("%s%s", insString, separator)
	res := strings.Repeat(pattern, count-1) + insString
	return res
}
