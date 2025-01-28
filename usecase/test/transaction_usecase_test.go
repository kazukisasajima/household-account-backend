package usecase_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"household-account-backend/entity"
	"household-account-backend/usecase"
)

type mockTransactionRepository struct {
	mock.Mock
}

func NewMockTransactionRepository() *mockTransactionRepository {
	return new(mockTransactionRepository)
}

func (m *mockTransactionRepository) CreateTransaction(transaction *entity.Transaction) (*entity.Transaction, error) {
	args := m.Called(transaction)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *mockTransactionRepository) GetTransactionByID(transactionID int) (*entity.Transaction, error) {
	args := m.Called(transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *mockTransactionRepository) GetTransactionsByUserID(userID int) ([]entity.Transaction, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Transaction), args.Error(1)
}

func (m *mockTransactionRepository) UpdateTransaction(transaction *entity.Transaction) (*entity.Transaction, error) {
	args := m.Called(transaction)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *mockTransactionRepository) DeleteTransaction(transactionID int) error {
	args := m.Called(transactionID)
	return args.Error(0)
}

type TransactionUseCaseSuite struct {
	suite.Suite
	transactionUseCase usecase.TransactionUseCase
}

func TestTransactionUseCaseSuite(t *testing.T) {
	suite.Run(t, new(TransactionUseCaseSuite))
}

func (suite *TransactionUseCaseSuite) SetupTest() {
	mockRepo := NewMockTransactionRepository()
	suite.transactionUseCase = usecase.NewTransactionUseCase(mockRepo)
}

func (suite *TransactionUseCaseSuite) TestCreateTransaction() {
	transaction := &entity.Transaction{
		UserID:     1,
		CategoryID: 1,
		Date:       time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
		Amount:     100.00,
		Content:    "Groceries",
	}

	mockRepo := NewMockTransactionRepository()
	suite.transactionUseCase = usecase.NewTransactionUseCase(mockRepo)
	mockRepo.On("CreateTransaction", transaction).Return(transaction, nil)

	createdTransaction, err := suite.transactionUseCase.CreateTransaction(transaction)
	suite.Assert().Nil(err)
	suite.Assert().Equal(100.00, createdTransaction.Amount)
	suite.Assert().Equal("Groceries", createdTransaction.Content)
}

func (suite *TransactionUseCaseSuite) TestGetTransactionByID() {
	transaction := &entity.Transaction{
		ID:         1,
		UserID:     1,
		CategoryID: 1,
		Date:       time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
		Amount:     100.00,
		Content:    "Groceries",
	}

	mockRepo := NewMockTransactionRepository()
	suite.transactionUseCase = usecase.NewTransactionUseCase(mockRepo)
	mockRepo.On("GetTransactionByID", transaction.ID).Return(transaction, nil)

	retrievedTransaction, err := suite.transactionUseCase.GetTransactionByID(transaction.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal(100.00, retrievedTransaction.Amount)
	suite.Assert().Equal("Groceries", retrievedTransaction.Content)
}

func (suite *TransactionUseCaseSuite) TestGetTransactionsByUserID() {
	transactions := []entity.Transaction{
		{
			ID:         1,
			UserID:     1,
			CategoryID: 1,
			Date:       time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			Amount:     100.00,
			Content:    "Groceries",
		},
		{
			ID:         2,
			UserID:     1,
			CategoryID: 2,
			Date:       time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC),
			Amount:     200.00,
			Content:    "Rent",
		},
	}

	mockRepo := NewMockTransactionRepository()
	suite.transactionUseCase = usecase.NewTransactionUseCase(mockRepo)
	mockRepo.On("GetTransactionsByUserID", 1).Return(transactions, nil)

	retrievedTransactions, err := suite.transactionUseCase.GetTransactionsByUserID(1)
	suite.Assert().Nil(err)
	suite.Assert().Equal(2, len(retrievedTransactions))
	suite.Assert().Equal("Groceries", retrievedTransactions[0].Content)
	suite.Assert().Equal("Rent", retrievedTransactions[1].Content)
}

func (suite *TransactionUseCaseSuite) TestUpdateTransaction() {
	transaction := &entity.Transaction{
		ID:         1,
		UserID:     1,
		CategoryID: 1,
		Date:       time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
		Amount:     150.00,
		Content:    "Updated Groceries",
	}

	mockRepo := NewMockTransactionRepository()
	suite.transactionUseCase = usecase.NewTransactionUseCase(mockRepo)
	mockRepo.On("UpdateTransaction", transaction).Return(transaction, nil)

	updatedTransaction, err := suite.transactionUseCase.UpdateTransaction(transaction)
	suite.Assert().Nil(err)
	suite.Assert().Equal(150.00, updatedTransaction.Amount)
	suite.Assert().Equal("Updated Groceries", updatedTransaction.Content)
}

func (suite *TransactionUseCaseSuite) TestDeleteTransaction() {
	mockRepo := NewMockTransactionRepository()
	suite.transactionUseCase = usecase.NewTransactionUseCase(mockRepo)
	mockRepo.On("DeleteTransaction", 1).Return(nil)

	err := suite.transactionUseCase.DeleteTransaction(1)
	suite.Assert().Nil(err)
}
