package domain

/* STATUS
toBeConfirmed
toBeRetrieved
retrieved
rejected
*/

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

type Direction struct {
	DirectionId int64  `json:"direction_id"`
	Street      string `json:"street"`
	Number      int64  `json:"number"`
	Details     string `json:"details"`
	City        string `json:"city"`
	PostalCode  int64  `json:"postal_code"`
}
