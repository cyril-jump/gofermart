package dto

type Task struct {
	UserID   string
	NumOrder string
	IsNew    bool
}

type AccrualResponse struct {
	NumOrder    string  `json:"order"`
	OrderStatus string  `json:"status"`
	Accrual     float32 `json:"accrual"`
}

type Order struct {
	NumOrder    string  `json:"number"`
	OrderStatus string  `json:"status"`
	Accrual     float32 `json:"accrual"`
	UploadedAt  string  `json:"uploaded_at"`
}

type User struct {
	UserID   string `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserBalance struct {
	Current   string `json:"current"`
	Withdrawn string `json:"withdrawn"`
}
