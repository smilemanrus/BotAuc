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
				WHERE Status = '%aucStatus%'
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

const trunsctnQuery = `BEGIN TRANSACTION;
								%logic%
								;
								COMMIT TRANSACTION`

const mkngDataForFix = `INSERT INTO DataForFix (USERID, NOTIFICATIONTYPE, URL, EVENTDATE)
						%insertTables%
						%logic%
						DROP TABLE DataForFix`
