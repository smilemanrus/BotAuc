package sqlite

const makeAlertQuery = `BEGIN TRANSACTION;
				CREATE TABLE IF NOT EXISTS aucsWithoutEvent
				(
				Name      TEXT,
				URL       TEXT,
				StartDate DATETIME,
				EndDate   DATETIME
				);
				
				INSERT INTO aucsWithoutEvent (Name, URL, StartDate, EndDate)
				SELECT aucs.Name, aucs.URL, aucs.StartDate, aucs.EndDate
				FROM aucs
				LEFT JOIN alerts on aucs.URL = alerts.URL
				AND alerts.EVENTDATE IS NULL
				WHERE Status = '%aucStatus'
				AND ((UNIXEPOCH(StartDate) - UNIXEPOCH()) / 60 / 60) = 0
				AND alerts.URL IS NULL
				;
				INSERT INTO alerts (USERID, NOTIFICATIONTYPE, URL, EVENTDATE)
				SELECT subscribes.ID,
				'%notificationType%',
				aucsWithoutEvent.URL,
				NULL
				
				FROM subscribes
				CROSS JOIN aucsWithoutEvent;
				DROP TABLE aucsWithoutEvent;
				COMMIT TRANSACTION`

const makingTempTableForFix = `BEGIN TRANSACTION;
								CREATE TABLE IF NOT EXISTS DataForFix
								(
									USERID           INTEGER,
									NOTIFICATIONTYPE TEXT,
									URL              TEXT
								);
								%logic%
								;
								
								DROP TABLE DataForFix;
								COMMIT TRANSACTION`
