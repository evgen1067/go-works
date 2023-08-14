package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/evgen1067/anti-bruteforce/internal/common"
	"github.com/evgen1067/anti-bruteforce/internal/config"
	_ "github.com/jackc/pgx/v4/stdlib" //nolint:blank-imports
)

type Repo struct {
	db  *sql.DB
	cfg *config.Config
}

func NewRepo(cfg *config.Config) *Repo {
	repo := new(Repo)
	repo.cfg = cfg
	return repo
}

func (r *Repo) Connect(ctx context.Context) (err error) {
	r.db, err = sql.Open("pgx", getDSN(r.cfg))
	if err != nil {
		return fmt.Errorf("failed to load driver: %w", err)
	}
	return r.db.PingContext(ctx)
}

func (r *Repo) Close() error {
	return r.db.Close()
}

func getDSN(cfg *config.Config) string {
	return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Database,
		cfg.DB.SSLMode)
}

func (r *Repo) AddToBlacklist(ctx context.Context, address string) error {
	return r.Add(ctx, common.Blacklist, address)
}

func (r *Repo) AddToWhitelist(ctx context.Context, address string) error {
	return r.Add(ctx, common.Whitelist, address)
}

func (r *Repo) ExistsInBlacklist(ctx context.Context, address string) (bool, error) {
	return r.Exists(ctx, common.Blacklist, address)
}

func (r *Repo) ExistsInWhitelist(ctx context.Context, address string) (bool, error) {
	return r.Exists(ctx, common.Whitelist, address)
}

func (r *Repo) DeleteFromBlacklist(ctx context.Context, address string) error {
	return r.Delete(ctx, common.Blacklist, address)
}

func (r *Repo) DeleteFromWhitelist(ctx context.Context, address string) error {
	return r.Delete(ctx, common.Whitelist, address)
}

func (r *Repo) Add(ctx context.Context, table common.TableName, address string) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() ////nolint:errcheck

	var (
		query  string
		status bool
	)
	switch table {
	case common.Blacklist:
		query = `INSERT INTO blacklist (address)
			  	 VALUES ($1::inet)
			 	 ON CONFLICT (address) DO NOTHING
			 	 RETURNING true status;`
	case common.Whitelist:
		query = `INSERT INTO whitelist (address)
			 	 VALUES ($1::inet)
			 	 ON CONFLICT (address) DO NOTHING
			 	 RETURNING true status;`
	}

	err = tx.QueryRowContext(
		ctx,
		query,
		address).Scan(&status)

	if errors.Is(sql.ErrNoRows, err) {
		return common.ErrIPExists
	} else if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *Repo) Exists(ctx context.Context, table common.TableName, address string) (bool, error) {
	var (
		query  string
		exists bool
	)
	switch table {
	case common.Blacklist:
		query = `SELECT true blacklist
			  	 FROM blacklist
			  	 WHERE address >>= $1::inet;`
	case common.Whitelist:
		query = `SELECT true whitelist
			 	 FROM whitelist
				 WHERE address >>= $1::inet;`
	}

	err := r.db.QueryRowContext(ctx, query, address).Scan(&exists)
	if errors.Is(sql.ErrNoRows, err) {
		return false, nil
	} else if err != nil {
		return exists, err
	}
	return exists, nil
}

func (r *Repo) Delete(ctx context.Context, table common.TableName, address string) (err error) {
	var (
		query  string
		status bool
	)
	switch table {
	case common.Blacklist:
		query = `DELETE
			 	 FROM blacklist
			 	 WHERE address = $1::inet
			 	 RETURNING true status;`
	case common.Whitelist:
		query = `DELETE
			 	 FROM whitelist
			 	 WHERE address = $1::inet
			  	 RETURNING true status;`
	}

	err = r.db.QueryRowContext(ctx, query, address).Scan(&status)

	if errors.Is(sql.ErrNoRows, err) {
		return common.ErrIPNotExists
	} else if err != nil {
		return err
	}

	return nil
}
