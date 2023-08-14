package psql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/evgen1067/hw12_13_14_15_calendar/internal/common"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/config"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/v4/stdlib" //nolint:blank-imports
)

type Storage struct {
	db *sql.DB
}

func NewStorage() storage.DBStorage {
	return new(Storage)
}

func (store *Storage) Connect(ctx context.Context, cfg *config.Config) (err error) {
	store.db, err = sql.Open("pgx", getDSN(cfg))
	if err != nil {
		return fmt.Errorf("failed to load driver: %w", err)
	}
	return store.db.PingContext(ctx)
}

func (store *Storage) Close() error {
	return store.db.Close()
}

func getDSN(cfg *config.Config) string {
	return fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Database,
		cfg.DB.SSLMode)
}

func (store *Storage) Create(ctx context.Context, event common.Event) (common.EventID, error) {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return event.ID, err
	}
	defer tx.Rollback()

	query := `INSERT INTO events (title, description, date_start, 
                    date_end, owner_id, notify_in) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err = tx.QueryRowContext(
		ctx,
		query,
		event.Title, event.Description, event.DateStart, event.DateEnd, event.OwnerID, event.NotifyIn).Scan(&event.ID)
	if err != nil {
		return event.ID, err
	}

	if err = tx.Commit(); err != nil {
		return event.ID, err
	}

	return event.ID, nil
}

func (store *Storage) Update(ctx context.Context, id common.EventID, event common.Event) (common.EventID, error) {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return id, err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	query := `UPDATE events SET title = $1, description = $2, date_start = $3, 
                  date_end = $4, owner_id = $5, notify_in = $6 WHERE id = $7`

	result, err := tx.ExecContext(
		ctx,
		query,
		event.Title, event.Description, event.DateStart, event.DateEnd, event.OwnerID, event.NotifyIn, id)
	if err != nil {
		return id, err
	}
	notFound, err := result.RowsAffected()
	if err != nil {
		return id, err
	}
	if notFound == 0 {
		return id, common.ErrNotFound
	}
	if err = tx.Commit(); err != nil {
		return event.ID, err
	}
	return id, nil
}

func (store *Storage) Delete(ctx context.Context, id common.EventID) (common.EventID, error) {
	query := `DELETE FROM events WHERE id = $1`
	result, err := store.db.ExecContext(ctx, query, id)
	if err != nil {
		return id, err
	}
	notFound, err := result.RowsAffected()
	if err != nil {
		return id, err
	}
	if notFound == 0 {
		return id, common.ErrNotFound
	}
	return id, nil
}

func (store *Storage) PeriodList(
	ctx context.Context,
	startPeriod time.Time,
	period storage.Period,
) ([]common.Event, error) {
	var endPeriod time.Time
	switch period {
	case "Day":
		endPeriod = startPeriod.AddDate(0, 0, 1)
	case "Week":
		endPeriod = startPeriod.AddDate(0, 0, 7)
	case "Month":
		endPeriod = startPeriod.AddDate(0, 1, 0)
	}
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	query := "SELECT * FROM events e WHERE e.date_end > $1 AND $2 > e.date_start"

	rows, err := tx.QueryContext(
		ctx,
		query,
		startPeriod.Format("2006-01-02 15:04"),
		endPeriod.Format("2006-01-02 15:04"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []common.Event
	for rows.Next() {
		var event common.Event
		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.DateStart,
			&event.DateEnd,
			&event.OwnerID,
			&event.NotifyIn); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

func (store *Storage) DayList(ctx context.Context, startDate time.Time) ([]common.Event, error) {
	period := storage.Period("Day")
	return store.PeriodList(ctx, startDate, period)
}

func (store *Storage) WeekList(ctx context.Context, startDate time.Time) ([]common.Event, error) {
	period := storage.Period("Week")
	return store.PeriodList(ctx, startDate, period)
}

func (store *Storage) MonthList(ctx context.Context, startDate time.Time) ([]common.Event, error) {
	period := storage.Period("Month")
	return store.PeriodList(ctx, startDate, period)
}

func (store *Storage) ClearOldEvents(ctx context.Context) error {
	lastYearDate := time.Now().AddDate(-1, 0, 0)
	query := `DELETE FROM events WHERE date_end < $1`

	result, err := store.db.ExecContext(ctx, query, lastYearDate.Format("2006-01-02 15:04"))
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (store *Storage) SchedulerList(ctx context.Context) ([]common.Notice, error) {
	var notices []common.Notice
	query := `SELECT id, title, date_start, owner_id
				FROM events
				WHERE date_start > (now() - notify_in * interval '1 hour')
				  AND date_start <= (now() + notify_in * interval '1 hour')`

	rows, err := store.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var notice common.Notice
		if err := rows.Scan(
			&notice.EventID,
			&notice.Title,
			&notice.Datetime,
			&notice.OwnerID,
		); err != nil {
			return nil, err
		}
		notices = append(notices, notice)
	}
	return notices, rows.Err()
}
