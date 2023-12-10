package repository

import (
	"context"
	"time"

	"github.com/iqbaludinm/hr-microservice/profile-service/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	WithTransaction(ctx context.Context, fn func(pgx.Tx) error) error
	WithoutTransaction(ctx context.Context, fn func(*pgxpool.Pool) error) error
}

type StoreImpl struct {
	Db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &StoreImpl{
		Db: db,
	}
}

// This function is used for starting transaction. like POST, PUT, DELETE
func (r *StoreImpl) WithTransaction(ctx context.Context, fn func(pgx.Tx) error) error {
	// context.WithTimeout is used for setting a timeout for this operation.
	c, cancel := context.WithTimeout(context.Background(), time.Duration(config.TimeOutDuration)*time.Second)
	defer cancel()

	// begin transaction for this operation.
	tx, err := r.Db.Begin(c)
	if err != nil {
		return err
	}

	// run function with transaction db.
	// if funtion return error then rollback and return error
	if err := fn(tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	// commit transaction
	return tx.Commit(ctx)
}

// This function is used for query without transactions. like GET
func (r *StoreImpl) WithoutTransaction(ctx context.Context, fn func(*pgxpool.Pool) error) error {
	// run function with context db.
	if err := fn(r.Db); err != nil {
		return err
	}

	return nil
}
