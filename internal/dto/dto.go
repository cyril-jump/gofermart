package dto

type Task struct {
	NumOrder string
	IsNew    bool
}

type AccrualResponse struct {
	UserID      string  `json:"id, omitempty"`
	NumOrder    string  `json:"order"`
	OrderStatus string  `json:"status"`
	Accrual     float64 `json:"accrual"`
}

type User struct {
	UserID   string `json:"id, omitempty"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
