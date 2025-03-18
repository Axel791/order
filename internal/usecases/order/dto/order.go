package dto

type CreateOrder struct {
	UserID     int64
	Code       string
	TotalPrice int64
}

type Order struct {
	ID         int64
	UserID     int64
	Code       string
	TotalPrice int64
}
