package orders

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/leenzstra/WBL0/internal/models"
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

func TestGetOk(t *testing.T) {
	db := prepareRepoMocks()
	repo := NewRepo(db)

	var o []models.OrderModel
	repo.DB.Scanner.(*mocks.IScanner).
		On("Select", context.Background(), repo.DB.Pool, &o, 
			"SELECT * FROM orders WHERE order_uid = $1", order.OrderUID).
		Times(1).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(2).(*[]models.OrderModel)
			(*arg) = append((*arg), order)
		})

	oneOrder, err := repo.Get(context.Background(), order.OrderUID)

	repo.DB.Scanner.(*mocks.IScanner).AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, oneOrder, &order)
}

func TestGetErrorOnSelect(t *testing.T) {
	db := prepareRepoMocks()
	repo := NewRepo(db)

	var o []models.OrderModel
	repo.DB.Scanner.(*mocks.IScanner).
		On("Select", context.Background(), repo.DB.Pool, &o, 
			"SELECT * FROM orders WHERE order_uid = $1", order.OrderUID).
		Times(1).
		Return(errors.New("some select error")).
		Run(func(args mock.Arguments) {
			arg := args.Get(2).(*[]models.OrderModel)
			(*arg) = append((*arg), order)
		})

	oneOrder, err := repo.Get(context.Background(), order.OrderUID)

	repo.DB.Scanner.(*mocks.IScanner).AssertExpectations(t)
	assert.Error(t, err)
	assert.Nil(t, oneOrder)
	assert.ErrorContains(t, err, "OrderRepo.Get.Select")
}

func TestGetErrorOnLen(t *testing.T) {
	db := prepareRepoMocks()
	repo := NewRepo(db)

	var o []models.OrderModel
	repo.DB.Scanner.(*mocks.IScanner).
		On("Select", context.Background(), repo.DB.Pool, &o, 
			"SELECT * FROM orders WHERE order_uid = $1", order.OrderUID).
		Times(1).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(2).(*[]models.OrderModel)
			(*arg) = append((*arg), order, order)
		})

	oneOrder, err := repo.Get(context.Background(), order.OrderUID)

	repo.DB.Scanner.(*mocks.IScanner).AssertExpectations(t)
	assert.Error(t, err)
	assert.Nil(t, oneOrder)
	assert.ErrorContains(t, err, "OrderRepo.Get.row.len")
}

func TestGetAllOk(t *testing.T) {
	db := prepareRepoMocks()
	repo := NewRepo(db)

	var o []models.OrderModel
	repo.DB.Scanner.(*mocks.IScanner).
		On("Select", context.Background(), repo.DB.Pool, &o, "SELECT * FROM orders").
		Times(1).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(2).(*[]models.OrderModel)
			(*arg) = append((*arg), order)
		})

	os, err := repo.GetAll(context.Background())

	repo.DB.Scanner.(*mocks.IScanner).AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(os))
	assert.Equal(t, os[0], order)
}

func TestGetAllErrorOnSelect(t *testing.T) {
	db := prepareRepoMocks()
	repo := NewRepo(db)

	var o []models.OrderModel
	repo.DB.Scanner.(*mocks.IScanner).
		On("Select", context.Background(), repo.DB.Pool, &o, "SELECT * FROM orders").
		Times(1).
		Return(errors.New("some select error")).
		Run(func(args mock.Arguments) {
			arg := args.Get(2).(*[]models.OrderModel)
			(*arg) = append((*arg), order)
		})

	os, err := repo.GetAll(context.Background())

	repo.DB.Scanner.(*mocks.IScanner).AssertExpectations(t)
	assert.Error(t, err)
	assert.Nil(t, os)
	assert.ErrorContains(t, err, "OrderRepo.GetAll.Select")
}
