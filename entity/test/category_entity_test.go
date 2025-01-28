package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"household-account-backend/entity"
)

func TestCategory(t *testing.T) {
	category := entity.Category{
		ID:     1,
		UserID: 2,
		Name:   "Food",
		Type:   "expense",
	}
	assert.Equal(t, 1, category.ID)
	assert.Equal(t, 2, category.UserID)
	assert.Equal(t, "Food", category.Name)
	assert.Equal(t, "expense", category.Type)
}
