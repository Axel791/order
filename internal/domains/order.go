package domains

import "github.com/Axel791/appkit"

type Order struct {
	ID         int64
	UserID     int64
	Code       string
	TotalPrice int64
}

func (v *Order) ValidateUserID() error {
	if v.UserID <= 0 {
		return appkit.ValidationError("error invalid user id")
	}
	return nil
}

func (v *Order) ValidateTotalPrice() error {
	if v.TotalPrice <= 0 {
		return appkit.ValidationError("error invalid total price")
	}
	return nil
}
