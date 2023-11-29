package database

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DB struct {
	Pool    *pgxpool.Pool
	Logger  *zap.Logger
	Builder *squirrel.StatementBuilderType
}

func New(connUrl string, logger *zap.Logger) (*DB, error) {
	config, err := pgxpool.ParseConfig(connUrl)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	logger.Info("Database started")

	return &DB{
		Pool:    pool,
		Logger:  logger,
		Builder: &builder,
	}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}
