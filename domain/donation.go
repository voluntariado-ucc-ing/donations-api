package domain

const (
	StatusToBeConfirmed = "to_be_confirmed"
	StatusToBeRetrieved = "to_be_retrieved"
	StatusRetrieved     = "retrieved"
	StatusRejected      = "rejected"
)

type Donation struct {
	DonationId  int64      `json:"donation_id"`
	Quantity    int64      `json:"quantity"`
	Unit        string     `json:"unit"`
	Description string     `json:"description"`
	Element     string     `json:"element"`
	TypeId      int64      `json:"type_id"`
	DonorId     int64      `json:"donator_id,omitempty"`
	Date        string     `json:"donation_date"`
	DirectionId int64      `json:"direction_id"`
	Direction   *Direction `json:"direction,omitempty"`
	Status      string     `json:"status"`
}

type DonationRequest struct {
	Donations []Donation `json:"donations"`
	Donor     Donor      `json:"donator"`
}

type DonatorRequest struct {¿
	Donor     Donor      `json:"donator"`
}

type DonationConcurrent struct {
	Donation *Donation
	Error    ApiError
}

type Direction struct {
	DirectionId int64  `json:"direction_id"`
	Street      string `json:"street"`
	Number      int64  `json:"number"`
	Details     string `json:"details"`
	City        string `json:"city"`
	PostalCode  int64  `json:"postal_code"`
}

type StatusRequest struct {
	Status string `json:"status"`
}

func (s *StatusRequest) IsValidStatus() bool {
	return s.Status == StatusRejected ||
		s.Status == StatusRetrieved ||
		s.Status == StatusToBeConfirmed ||
		s.Status == StatusToBeRetrieved
}
