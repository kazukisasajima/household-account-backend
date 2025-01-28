package entity

type MonthlySummary struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	YearMonth string  `json:"year_month"` // Format: YYYY-MM
	Income    float64 `json:"income"`
	Expense   float64 `json:"expense"`
	Balance   float64 `json:"balance"`
}
