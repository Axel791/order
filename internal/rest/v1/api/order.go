package api

type InputCreateOrder struct {
	UserID     int64  `json:"userID"`
	Code       string `json:"code"`
	TotalPrice int64  `json:"totalPrice"`
}

type OrderResponse struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"userID"`
	Code       string `json:"code"`
	TotalPrice int64  `json:"totalPrice"`
}
