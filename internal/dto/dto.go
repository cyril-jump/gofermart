package dto

type Task struct {
	UserID   string
	NumOrder string
	IsNew    bool
}

type AccrualResponse struct {
	NumOrder    string  `json:"order"`
	OrderStatus string  `json:"status"`
	Accrual     float64 `json:"accrual"`
}
