package domain

type Donor struct {
	DonorId     int64  `json:"donator_id,omitempty"`
	Mail        string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
}
