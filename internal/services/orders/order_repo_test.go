package orders

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/leenzstra/WBL0/mocks"
	"github.com/leenzstra/WBL0/pkg/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func prepareRepoMocks() *database.DB {
	pool := mocks.IPool{}
	logger := zap.NewNop()
	builder := database.PgxBuilder
	scanner := mocks.IScanner{}

	db, err := database.New(&pool, builder, &scanner, logger)
	if err != nil {
		panic(err)
	}

	return db
}

func TestNew(t *testing.T) {
	db := prepareRepoMocks()

	repo := NewRepo(db)

	assert.NotNil(t, repo)
}

func TestAddOk(t *testing.T) {
	db := prepareRepoMocks()
	repo := NewRepo(db)

	repo.DB.Pool.(*mocks.IPool).
		On("Exec", context.Background(), mock.AnythingOfType("string"), mock.AnythingOfType("string"),
			mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("DeliveryModel"),
			mock.AnythingOfType("PaymentModel"), mock.AnythingOfType("OrderItems"), mock.AnythingOfType("string"),
			mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"),
			mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("Time"), 
			mock.AnythingOfType("string")).
		Times(1).
		Return(pgconn.CommandTag{}, nil)

	err := repo.Add(context.Background(), order)

	repo.DB.Pool.(*mocks.IPool).AssertExpectations(t)
	assert.NoError(t, err)
}

func TestAddErrorOnExec(t *testing.T) {
	db := prepareRepoMocks()
	repo := NewRepo(db)

	repo.DB.Pool.(*mocks.IPool).
		On("Exec", context.Background(), mock.AnythingOfType("string"), mock.AnythingOfType("string"),
			mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("DeliveryModel"),
			mock.AnythingOfType("PaymentModel"), mock.AnythingOfType("OrderItems"), mock.AnythingOfType("string"),
			mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"),
			mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("Time"), 
			mock.AnythingOfType("string")).
		Times(1).
		Return(pgconn.CommandTag{}, errors.New("some error"))

	err := repo.Add(context.Background(), order)

	repo.DB.Pool.(*mocks.IPool).AssertExpectations(t)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "OrderRepo.Add.r.Pool.Exec")
}