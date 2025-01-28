package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"household-account-backend/entity"
)

func TestMonthlySummary(t *testing.T) {
	monthlySummary := entity.MonthlySummary{
		ID:        1,
		UserID:    2,
		YearMonth: "2025-01",
		Income:    5000.00,
		Expense:   3000.00,
		Balance:   2000.00,
	}
	assert.Equal(t, 1, monthlySummary.ID)
	assert.Equal(t, 2, monthlySummary.UserID)
	assert.Equal(t, "2025-01", monthlySummary.YearMonth)
	assert.Equal(t, 5000.00, monthlySummary.Income)
	assert.Equal(t, 3000.00, monthlySummary.Expense)
	assert.Equal(t, 2000.00, monthlySummary.Balance)
}
