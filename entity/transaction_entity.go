package entity

import "time"

type Transaction struct {
	ID         int     `json:"id"`
	UserID     int     `json:"user_id"`
	CategoryID int     `json:"category_id"`
	Date       time.Time  `json:"date"`       // Format: YYYY-MM-DD
	Amount     float32 `json:"amount"`     // Decimal value
	Content    string  `json:"content"`    // Optional description
}
