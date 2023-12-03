package database

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	PgxBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type DB struct {
	Pool    IPool
	Logger  *zap.Logger
	Builder *squirrel.StatementBuilderType
	Scanner IScanner
}

func NewPgxScanner() (*pgxscan.API, error) {
	dbscanApi, err := pgxscan.NewDBScanAPI()
	if err != nil {
		return nil, err
	}

	return pgxscan.NewAPI(dbscanApi)
}

func NewPgxPool(connUrl string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connUrl)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return pool, err
}

func New(pool IPool, builder squirrel.StatementBuilderType, scanner IScanner, logger *zap.Logger) (*DB, error) {
	logger.Info("Database started")

	return &DB{
		Pool:    pool,
		Logger:  logger,
		Builder: &builder,
		Scanner: scanner,
	}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}
