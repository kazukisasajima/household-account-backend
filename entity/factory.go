package entity

func NewDomains() []interface{} {
	return []interface{}{
		User{},
		Category{},
		Transaction{},
		MonthlySummary{},
	}
}