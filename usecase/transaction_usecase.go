package usecase

import (
	"household-account-backend/adapter/gateway"
	"household-account-backend/entity"
)

type TransactionUseCase interface {
	CreateTransaction(transaction *entity.Transaction) (*entity.Transaction, error)
	GetTransactionByID(userID int, transactionID int) (*entity.Transaction, error)
	GetTransactionsByUserID(userID int) ([]entity.Transaction, error)
	UpdateTransaction(transaction *entity.Transaction) (*entity.Transaction, error)
	DeleteTransaction(userID int, transactionID int) error
}

type transactionUseCase struct {
	transactionRepository gateway.TransactionRepository
}

func NewTransactionUseCase(transactionRepository gateway.TransactionRepository) TransactionUseCase {
	return &transactionUseCase{
		transactionRepository: transactionRepository,
	}
}

func (tu *transactionUseCase) CreateTransaction(transaction *entity.Transaction) (*entity.Transaction, error) {
	return tu.transactionRepository.CreateTransaction(transaction)
}

func (tu *transactionUseCase) GetTransactionByID(userID int, transactionID int) (*entity.Transaction, error) {
	return tu.transactionRepository.GetTransactionByID(userID, transactionID)
}

func (tu *transactionUseCase) GetTransactionsByUserID(userID int) ([]entity.Transaction, error) {
	return tu.transactionRepository.GetTransactionsByUserID(userID)
}

func (tu *transactionUseCase) UpdateTransaction(transaction *entity.Transaction) (*entity.Transaction, error) {
	return tu.transactionRepository.UpdateTransaction(transaction)
}

func (tu *transactionUseCase) DeleteTransaction(userID int, transactionID int) error {
	return tu.transactionRepository.DeleteTransaction(userID, transactionID)
}
