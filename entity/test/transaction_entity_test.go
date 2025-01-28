package entity_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"household-account-backend/entity"
)

func TestTransaction(t *testing.T) {
	transaction := entity.Transaction{
		ID:         1,
		UserID:     2,
		Date:       time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
		Amount:     1000.50,
		Content:    "Grocery shopping",
	}
	assert.Equal(t, 1, transaction.ID)
	assert.Equal(t, 2, transaction.UserID)
	assert.Equal(t, 3, transaction.CategoryID)
	assert.Equal(t, "2025-01-01", transaction.Date)
	assert.Equal(t, 1000.50, transaction.Amount)
	assert.Equal(t, "Grocery shopping", transaction.Content)
}
