package dto

type Task struct {
	NumOrder string
	IsNew    bool
}

type AccrualResponse struct {
	UserID      string  `json:"user_id, omitempty"`
	NumOrder    string  `json:"order"`
	OrderStatus string  `json:"status"`
	Accrual     float64 `json:"accrual"`
	UploadedAt  string  `json:"uploaded_at,omitempty"`
}

type User struct {
	UserID   string `json:"user_id, omitempty"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
