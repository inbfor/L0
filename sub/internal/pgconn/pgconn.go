package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	db *pgxpool.Pool
}

func InitPG(ctx context.Context, connUrl string) (*DB, error) {

	pool, err := pgxpool.New(ctx, connUrl)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	db := &DB{
		db: pool,
	}

	return db, nil
}

func (pg *DB) Close() {
	pg.db.Close()
}
