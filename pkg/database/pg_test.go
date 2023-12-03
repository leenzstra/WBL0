package database

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/leenzstra/WBL0/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func prepareMocks() (*mocks.IPool, squirrel.StatementBuilderType, *mocks.IScanner, *zap.Logger) {
	pool := mocks.IPool{}
	logger := zap.NewNop()
	builder := PgxBuilder
	scanner := mocks.IScanner{}

	return &pool, builder, &scanner, logger
}

func TestNew(t *testing.T) {
	p, b, s, l := prepareMocks()
	db, err := New(p, b, s, l)

	assert.NoError(t, err)
	assert.NotNil(t, db)
}