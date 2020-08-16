package services

import (
	"github.com/voluntariado-ucc-ing/donations-api/clients"
	"github.com/voluntariado-ucc-ing/donations-api/domain"
	"net/http"
)

var (
	DonationService donationServiceInterface = &donationService{}
)

type donationServiceInterface interface {
	CreateDonation(request domain.DonationRequest) (*domain.DonationRequest, domain.ApiError)
	GetDonatorByMail(mail string) (*domain.Donor, domain.ApiError)
}

type donationService struct {}

func (d donationService) CreateDonation(request domain.DonationRequest) (*domain.DonationRequest, domain.ApiError) {
	var donorId int64 = 0
	donor, err := clients.GetDonator(request.Donor.Mail)
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
	donor, err := clients.GetDonator(mail)
	if err != nil {
		return nil, err
	}
	return donor, nil
}
