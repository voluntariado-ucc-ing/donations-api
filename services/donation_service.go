package services

import (
	"net/http"

	"github.com/voluntariado-ucc-ing/donations-api/clients"
	"github.com/voluntariado-ucc-ing/donations-api/domain"
)

var (
	DonationService donationServiceInterface = &donationService{}
)

type donationServiceInterface interface {
	CreateDonation(request domain.DonationRequest) (*domain.DonationRequest, domain.ApiError)
	GetDonatorByMail(mail string) (*domain.Donor, domain.ApiError)
	GetDonation(id int64) (*domain.Donation, domain.ApiError)
	GetDonatorById(id int64) (*domain.Donor, domain.ApiError)
	GetAllDonations() ([]domain.Donation, domain.ApiError)
}

type donationService struct{}

func (d donationService) GetDonation(id int64) (*domain.Donation, domain.ApiError) {
	return clients.GetDonation(id)
}

func (d donationService) CreateDonation(request domain.DonationRequest) (*domain.DonationRequest, domain.ApiError) {
	var donorId int64 = 0
	donor, err := clients.GetDonatorByMail(request.Donor.Mail)
	if err != nil {
		if err.Status() != http.StatusNotFound {
			return nil, err
		}
		donorId, err = clients.InsertDonor(request.Donor)
		if err != nil {
			return nil, err
		}
	} else {
		donorId = donor.DonorId
	}

	for index := range request.Donations {
		directionId, err := clients.InsertDirection(*request.Donations[index].Direction)
		if err != nil {
			return nil, err
		}

		request.Donations[index].DirectionId = directionId
		request.Donations[index].DonorId = donorId

		donationId, err := clients.InsertDonation(request.Donations[index])
		if err != nil {
			return nil, err
		}

		request.Donations[index].DonationId = donationId
		request.Donations[index].Direction = nil
	}

	return &request, nil
}

func (d donationService) GetDonatorByMail(mail string) (*domain.Donor, domain.ApiError) {
	donor, err := clients.GetDonatorByMail(mail)
	if err != nil {
		return nil, err
	}
	return donor, nil
}

func (d donationService) GetDonatorById(id int64) (*domain.Donor, domain.ApiError) {
	donor, err := clients.GetDonatorById(id)
	if err != nil {
		return nil, err
	}
	return donor, nil
}

func (d donationService) GetAllDonations() ([]domain.Donation, domain.ApiError) {
	res := make([]domain.Donation, 0)
	ids, err := clients.GetAllDonationsIds()
	if err != nil {
		return nil, err
	}

	input := make(chan domain.DonationConcurrent)
	defer close(input)
	for _, id := range ids {
		go d.getConcurrentDonation(id, input)
	}

	for i := 0; i < len(ids); i++ {
		result := <-input
		if result.Error != nil {
			return nil, result.Error
		}
		res = append(res, *result.Donation)
	}

	return res, nil
}

func (d donationService) getConcurrentDonation(id int64, output chan domain.DonationConcurrent) {
	vol, err := d.GetDonation(id)
	output <- domain.DonationConcurrent{Donation: vol, Error: err}
	return

}
