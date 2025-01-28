package gateway

import (
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"household-account-backend/entity"
)

type TransactionRepository interface {
	CreateTransaction(transaction *entity.Transaction) (*entity.Transaction, error)
	GetTransactionByID(userID int, transactionID int) (*entity.Transaction, error)
	GetTransactionsByUserID(userID int) ([]entity.Transaction, error)
	UpdateTransaction(transaction *entity.Transaction) (*entity.Transaction, error)
	DeleteTransaction(userID int, transactionID int) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (tr *transactionRepository) CreateTransaction(transaction *entity.Transaction) (*entity.Transaction, error) {
	if err := tr.db.Create(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (tr *transactionRepository) GetTransactionByID(userID int, transactionID int) (*entity.Transaction, error) {
	transaction := &entity.Transaction{}
	if err := tr.db.Where("id = ? AND user_id = ?", transactionID, userID).First(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (tr *transactionRepository) GetTransactionsByUserID(userID int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := tr.db.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (tr *transactionRepository) UpdateTransaction(transaction *entity.Transaction) (*entity.Transaction, error) {
	// 既存データの取得
	selectedTransaction, err := tr.GetTransactionByID(transaction.UserID, transaction.ID)
	if err != nil {
		return nil, err
	}

	// フィールドをコピー（空の値を無視）
	if err := copier.CopyWithOption(selectedTransaction, transaction, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	// 更新
	if err := tr.db.Save(selectedTransaction).Error; err != nil {
		return nil, err
	}

	return selectedTransaction, nil
}


func (tr *transactionRepository) DeleteTransaction(userID int, transactionID int) error {
	if err := tr.db.Where("id = ? AND user_id = ?", transactionID, userID).Delete(&entity.Transaction{}).Error; err != nil {
		return err
	}
	return nil
}
